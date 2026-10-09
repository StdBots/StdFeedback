package clone

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/StdBots/StdFeedback/internal/config"
	"github.com/StdBots/StdFeedback/internal/credit"
	"github.com/StdBots/StdFeedback/internal/database"
)

// CloneBot represents a cloned bot instance.
type CloneBot struct {
	API         *tgbotapi.BotAPI
	OwnerID     int64
	DB          *database.MongoDB
	BotUsername string
	stopChan    chan struct{}
}

// NewCloneBot creates and starts a clone bot.
func NewCloneBot(botToken string, ownerID int64, cfg *config.Config, db *database.MongoDB) (*CloneBot, error) {
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	return &CloneBot{
		API:         api,
		OwnerID:     ownerID,
		DB:          db,
		BotUsername: api.Self.UserName,
		stopChan:    make(chan struct{}),
	}, nil
}

// Start begins the polling loop for the clone bot.
func (c *CloneBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.API.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				c.processUpdate(&update)
			}
		case <-c.stopChan:
			log.Printf("Clone bot %s stopped.", c.BotUsername)
			return
		}
	}
}

// processUpdate handles messages for the clone bot.
func (c *CloneBot) processUpdate(update *tgbotapi.Update) {
	msg := update.Message
	if msg == nil {
		return
	}

	if msg.IsCommand() && msg.Command() == "start" {
		text := fmt.Sprintf("Hello! I am a feedback clone bot.\n\n%s", credit.GetCreditFooter())
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		c.API.Send(reply)
		return
	}

	// If owner replies to a forwarded message
	if msg.From.ID == c.OwnerID && msg.ReplyToMessage != nil {
		// Forward reply back to user
		reply := tgbotapi.NewMessage(msg.ReplyToMessage.Chat.ID, msg.Text+"\n\n"+credit.GetCreditFooter())
		c.API.Send(reply)
		return
	}

	// Otherwise, forward user message to owner
	forward := tgbotapi.NewForward(c.OwnerID, msg.Chat.ID, msg.MessageID)
	c.API.Send(forward)

	ack := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Message sent to the bot owner.\n\n%s", credit.GetCreditFooter()))
	c.API.Send(ack)
}

// Stop stops the clone bot.
func (c *CloneBot) Stop() {
	close(c.stopChan)
}

// GetInfo returns information about the clone bot.
func (c *CloneBot) GetInfo() string {
	return fmt.Sprintf("Clone Bot @%s | Owner: %d\nPowered by STD DEEPANSHU", c.BotUsername, c.OwnerID)
}
