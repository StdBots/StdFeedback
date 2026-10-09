package analytics

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/StdBots/StdFeedback/internal/database"
)

// Analytics tracks bot usage stats.
type Analytics struct {
	DB           *database.MongoDB
	messageCount int64
	dailyUsers   map[string]map[int64]bool
	mu           sync.RWMutex
}

// NewAnalytics creates a new Analytics instance.
func NewAnalytics(db *database.MongoDB) *Analytics {
	return &Analytics{
		DB:         db,
		dailyUsers: make(map[string]map[int64]bool),
	}
}

// TrackMessage increments message count and tracks user.
func (a *Analytics) TrackMessage(userID int64) {
	atomic.AddInt64(&a.messageCount, 1)
	a.TrackNewUser(userID)
}

// TrackNewUser tracks daily active user.
func (a *Analytics) TrackNewUser(userID int64) {
	today := time.Now().Format("2006-01-02")
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if a.dailyUsers[today] == nil {
		a.dailyUsers[today] = make(map[int64]bool)
	}
	a.dailyUsers[today][userID] = true
}

// GetDailyActiveUsers returns the count of active users today.
func (a *Analytics) GetDailyActiveUsers() int64 {
	today := time.Now().Format("2006-01-02")
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if users, ok := a.dailyUsers[today]; ok {
		return int64(len(users))
	}
	return 0
}

// GetTotalMessages returns the total message count.
func (a *Analytics) GetTotalMessages() int64 {
	return atomic.LoadInt64(&a.messageCount)
}

// GetStats returns formatted stats text.
func (a *Analytics) GetStats() string {
	// These values would normally be fetched from the DB
	totalUsers := int64(0)
	openTickets := int64(0)
	activeClones := int64(0)
	
	return fmt.Sprintf(`📊 **Bot Statistics**
├ 👥 Total Users: %d
├ 📨 Total Messages: %d  
├ 📅 Active Today: %d
├ 🎫 Open Tickets: %d
├ 🤖 Active Clones: %d
└ ⚡ Powered by STD BOTS`, 
		totalUsers, 
		a.GetTotalMessages(), 
		a.GetDailyActiveUsers(), 
		openTickets, 
		activeClones)
}

// ResetDaily resets daily counters.
func (a *Analytics) ResetDaily() {
	today := time.Now().Format("2006-01-02")
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.dailyUsers = make(map[string]map[int64]bool)
	a.dailyUsers[today] = make(map[int64]bool)
}

// StartDailyReset starts a goroutine that resets counters at midnight.
func (a *Analytics) StartDailyReset() {
	go func() {
		for {
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			a.ResetDaily()
		}
	}()
}
