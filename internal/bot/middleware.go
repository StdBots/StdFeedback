package bot

import (
	"strings"
	"sync"
	"time"
)

// RateLimiter provides sliding window rate limiting.
type RateLimiter struct {
	maxRequests int
	window      time.Duration
	requests    map[int64][]time.Time
	mu          sync.RWMutex
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make(map[int64][]time.Time),
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	for range ticker.C {
		now := time.Now()
		rl.mu.Lock()
		for userID, times := range rl.requests {
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) < rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, userID)
			} else {
				rl.requests[userID] = valid
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if the user is allowed to make a request.
func (rl *RateLimiter) Allow(userID int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.requests[userID]

	var valid []time.Time
	for _, t := range times {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.maxRequests {
		rl.requests[userID] = valid // keep the valid ones
		return false
	}

	valid = append(valid, now)
	rl.requests[userID] = valid
	return true
}

// AntiSpam provides spam detection functionality.
type AntiSpam struct {
	userSpamScore   map[int64]int
	lastMessage     map[int64]string
	lastMessageTime map[int64]time.Time
	spamPatterns    []string
	mu              sync.RWMutex
}

// NewAntiSpam creates a new AntiSpam instance.
func NewAntiSpam() *AntiSpam {
	return &AntiSpam{
		userSpamScore:   make(map[int64]int),
		lastMessage:     make(map[int64]string),
		lastMessageTime: make(map[int64]time.Time),
		spamPatterns: []string{
			"http://", "https://", "t.me/", "crypto", "bitcoin", "investment", "join",
		},
	}
}

// IsSpam detects repeated messages, flooding, and common spam patterns.
func (as *AntiSpam) IsSpam(userID int64, text string) bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	now := time.Now()
	isSpam := false

	lowerText := strings.ToLower(text)
	for _, pattern := range as.spamPatterns {
		if strings.Contains(lowerText, pattern) {
			isSpam = true
			as.userSpamScore[userID] += 2
			break
		}
	}

	if lastMsg, ok := as.lastMessage[userID]; ok && lastMsg == text {
		isSpam = true
		as.userSpamScore[userID] += 1
	}

	if lastTime, ok := as.lastMessageTime[userID]; ok {
		if now.Sub(lastTime) < 500*time.Millisecond {
			isSpam = true
			as.userSpamScore[userID] += 1
		}
	}

	as.lastMessage[userID] = text
	as.lastMessageTime[userID] = now

	return isSpam
}

// AutoBanCheck checks if the user has triggered spam too many times.
func (as *AntiSpam) AutoBanCheck(userID int64) bool {
	as.mu.RLock()
	defer as.mu.RUnlock()

	// if score > 10, trigger ban logic
	return as.userSpamScore[userID] > 10
}
