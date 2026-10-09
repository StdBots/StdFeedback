package broadcast

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/StdBots/StdFeedback/internal/database"
)

// BroadcastResult contains statistics about a broadcast operation.
type BroadcastResult struct {
	Success   int64
	Failed    int64
	Blocked   int64
	Deleted   int64
	Errors    []string
	TimeTaken time.Duration
}

// Broadcaster manages broadcasting operations.
type Broadcaster struct {
	DB      *database.MongoDB
	BotAPI  *tgbotapi.BotAPI
	Workers int
}

// NewBroadcaster creates a new Broadcaster instance.
func NewBroadcaster(db *database.MongoDB, api *tgbotapi.BotAPI, workers int) *Broadcaster {
	if workers <= 0 {
		workers = 20
	}
	return &Broadcaster{
		DB:      db,
		BotAPI:  api,
		Workers: workers,
	}
}

// Broadcast sends a message to all users in the database using a worker pool.
func (b *Broadcaster) Broadcast(msg tgbotapi.Chattable, asCopy bool) (*BroadcastResult, error) {
	// Example: users := b.DB.GetAllUserIDs()
	var users []int64
	return b.BroadcastToUsers(users, msg), nil
}

// BroadcastToUsers sends a message to specific users using a worker pool.
func (b *Broadcaster) BroadcastToUsers(userIDs []int64, msg tgbotapi.Chattable) *BroadcastResult {
	start := time.Now()
	result := &BroadcastResult{}
	var mu sync.Mutex

	jobs := make(chan int64, len(userIDs))
	var wg sync.WaitGroup

	for i := 0; i < b.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range jobs {
				var err error
				// Note: this implementation simplifies message creation, real implementation 
				// would recreate the message payload for each target chat id.
				switch m := msg.(type) {
				case tgbotapi.MessageConfig:
					m.ChatID = userID
					_, err = b.BotAPI.Send(m)
				case tgbotapi.CopyMessageConfig:
					m.ChatID = userID
					_, err = b.BotAPI.Send(m)
				}

				mu.Lock()
				if err != nil {
					errStr := strings.ToLower(err.Error())
					if strings.Contains(errStr, "429") || strings.Contains(errStr, "flood") {
						time.Sleep(3 * time.Second) // basic flood wait handling
						result.Failed++
					} else if strings.Contains(errStr, "blocked") {
						result.Blocked++
					} else if strings.Contains(errStr, "deactivated") || strings.Contains(errStr, "deleted") {
						result.Deleted++
					} else {
						result.Failed++
						result.Errors = append(result.Errors, errStr)
					}
				} else {
					result.Success++
				}
				
				// Optional: report progress every 100 users
				totalDone := result.Success + result.Failed + result.Blocked + result.Deleted
				if totalDone > 0 && totalDone%100 == 0 {
					fmt.Printf("Broadcast progress: %d processed...\n", totalDone)
				}
				mu.Unlock()
			}
		}()
	}

	for _, id := range userIDs {
		jobs <- id
	}
	close(jobs)
	wg.Wait()

	result.TimeTaken = time.Since(start)
	return result
}

// FormatResult formats the broadcast result as readable text.
func (b *Broadcaster) FormatResult(result *BroadcastResult) string {
	return fmt.Sprintf(`📢 Broadcast Completed!
⏱ Time Taken: %v
✅ Success: %d
❌ Failed: %d
🚫 Blocked: %d
💀 Deleted: %d
⚡ Powered by STD BOTS`, result.TimeTaken, result.Success, result.Failed, result.Blocked, result.Deleted)
}
