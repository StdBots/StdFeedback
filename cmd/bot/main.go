package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/StdBots/StdFeedback/internal/analytics"
	"github.com/StdBots/StdFeedback/internal/bot"
	"github.com/StdBots/StdFeedback/internal/config"
	"github.com/StdBots/StdFeedback/internal/credit"
	"github.com/StdBots/StdFeedback/internal/database"
	"github.com/StdBots/StdFeedback/internal/ticket"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// 2. Load config
	cfg := config.MustLoad()

	// 3. Print ASCII banner
	credit.PrintBanner()

	// 4. Verify credit integrity
	if intact, _ := credit.VerifyIntegrity(); !intact {
		log.Println("WARNING: Credit tampered! Please respect the STD DEEPANSHU and STD BOTS copyright.")
	}

	// 5. Connect to MongoDB
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 6. Initialize ticket manager & analytics
	ticketManager := ticket.NewTicketManager(db)
	_ = ticketManager
	analyticsManager := analytics.NewAnalytics(db)

	// 7. Start analytics daily reset goroutine
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			if err := analyticsManager.ResetDaily(context.Background()); err != nil {
				log.Printf("Failed to reset daily analytics: %v", err)
			}
		}
	}()

	// 8. Create bot
	b, err := bot.NewBot(cfg, db)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// 9. Start credit periodic checker in background
	go credit.CheckCreditPeriodically(6*time.Hour, func() bool {
		intact, _ := credit.VerifyIntegrity()
		return intact
	})

	// 10. Report fork status
	go credit.ReportForkStatus(b.API.Self.UserName, true)

	// 11. Log bot started
	log.Printf("Bot started successfully! @%s\n", b.API.Self.UserName)

	// 12. Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 13. Start bot polling
	go func() {
		b.Start()
	}()

	// Wait for termination signal
	<-quit
	log.Println("Received shutdown signal. Stopping gracefully...")

	// 14. On shutdown signal
	b.Stop()

	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Bot stopped gracefully")
}
