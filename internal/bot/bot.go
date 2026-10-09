package bot

import (
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/StdBots/StdFeedback/internal/config"
	"github.com/StdBots/StdFeedback/internal/credit"
	"github.com/StdBots/StdFeedback/internal/database"
)

// Bot represents the core bot engine.
type Bot struct {
	API             *tgbotapi.BotAPI
	Config          *config.Config
	DB              *database.MongoDB
	Cache           *database.Cache
	RunningClones   map[string]*tgbotapi.BotAPI
	clonesMutex     sync.RWMutex
	UserMessageMap  map[int64]int64 // maps message_id to user_id for reply tracking
	messageMapMutex sync.RWMutex
	rateLimiter     *RateLimiter
	antiSpam        *AntiSpam
	stopChan        chan struct{}
}

// NewBot creates a new Bot instance.
func NewBot(cfg *config.Config, db *database.MongoDB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)
	log.Println(credit.GetCreditBanner())

	return &Bot{
		API:            api,
		Config:         cfg,
		DB:             db,
		RunningClones:  make(map[string]*tgbotapi.BotAPI),
		UserMessageMap: make(map[int64]int64),
		rateLimiter:    NewRateLimiter(5, 5*time.Second), // allow 5 requests per 5s
		antiSpam:       NewAntiSpam(),
		stopChan:       make(chan struct{}),
	}, nil
}

// Start starts the bot polling loop.
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			go b.processUpdate(update)
		case <-b.stopChan:
			log.Println("Stopping bot polling...")
			b.API.StopReceivingUpdates()
			return
		}
	}
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	close(b.stopChan)
}

// processUpdate routes the update to appropriate handlers.
func (b *Bot) processUpdate(update tgbotapi.Update) {
	var userID int64
	var text string

	if update.Message != nil {
		userID = update.Message.From.ID
		text = update.Message.Text
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.From.ID
	}

	// 1. Run rate limit middleware
	if userID != 0 && !b.IsSudo(userID) {
		if !b.rateLimiter.Allow(userID) {
			log.Printf("Rate limit exceeded for user %d", userID)
			return
		}

		if text != "" {
			if b.antiSpam.IsSpam(userID, text) {
				log.Printf("Spam detected from user %d", userID)
				return
			}
		}
	}

	if update.Message != nil {
		// 2. Check if message is a command
		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			return
		}

		// 3. Check if reply from owner
		if update.Message.ReplyToMessage != nil && b.IsSudo(update.Message.From.ID) {
			b.handleAdminReply(update.Message)
			return
		}

		// 4. Check if regular message in private chat (feedback)
		if update.Message.Chat.IsPrivate() && !b.IsSudo(update.Message.From.ID) {
			b.handleFeedback(update.Message)
			return
		}
	} else if update.CallbackQuery != nil {
		// 5. Check if callback query
		b.handleCallbackQuery(update.CallbackQuery)
		return
	}
}

// handleCommand routes commands to specific logic
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	log.Printf("Received command %s from %d", message.Command(), message.From.ID)
	switch message.Command() {
	case "start":
		b.handleStart(message)
	case "help":
		b.handleHelp(message)
	case "about":
		b.handleAbout(message)
	case "settings":
		b.handleSettings(message)
	case "donate":
		b.handleDonate(message)
	case "ban":
		b.handleBan(message)
	case "unban":
		b.handleUnban(message)
	case "stats":
		b.handleStats(message)
	case "broadcast":
		b.handleBroadcast(message)
	case "users":
		b.handleUsers(message)
	case "clone":
		b.handleClone(message)
	case "deleteclone":
		b.handleDeleteClone(message)
	case "clones":
		b.handleClones(message)
	default:
		// ignore unknown
	}
}

// handleAdminReply routes replies from sudo users to the original user
func (b *Bot) handleAdminReply(message *tgbotapi.Message) {
	b.messageMapMutex.RLock()
	userID, ok := b.UserMessageMap[int64(message.ReplyToMessage.MessageID)]
	b.messageMapMutex.RUnlock()

	if ok {
		msg := tgbotapi.NewMessage(userID, credit.WatermarkMessage(message.Text))
		_, err := b.API.Send(msg)
		if err != nil {
			log.Printf("Failed to send reply to user %d: %v", userID, err)
		}
	} else {
		log.Println("Could not find user ID for reply")
	}
}

// handleFeedback forwards regular messages to admin
func (b *Bot) handleFeedback(message *tgbotapi.Message) {
	msg, err := b.ForwardMessage(b.Config.OwnerID, message.Chat.ID, message.MessageID)
	if err == nil {
		b.messageMapMutex.Lock()
		b.UserMessageMap[int64(msg.MessageID)] = message.From.ID
		b.messageMapMutex.Unlock()

		confirm := tgbotapi.NewMessage(message.Chat.ID, "✅ Message sent to admin! You will get a reply soon.")
		b.API.Send(confirm)
	} else {
		log.Printf("Failed to forward feedback: %v", err)
	}
}

// handleCallbackQuery processes callback queries
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	log.Printf("Received callback query %s from %d", query.Data, query.From.ID)
	b.handleCallback(query)
}

// SendMessage sends a text message with the bot watermark applied.
func (b *Bot) SendMessage(chatID int64, text string, parseMode string) (tgbotapi.Message, error) {
	watermarkedText := credit.WatermarkMessage(text)
	msg := tgbotapi.NewMessage(chatID, watermarkedText)
	if parseMode != "" {
		msg.ParseMode = parseMode
	}
	return b.API.Send(msg)
}

// SendMessageWithKeyboard sends a message with an inline keyboard.
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	watermarkedText := credit.WatermarkMessage(text)
	msg := tgbotapi.NewMessage(chatID, watermarkedText)
	msg.ReplyMarkup = keyboard
	return b.API.Send(msg)
}

// ForwardMessage forwards a message from one chat to another.
func (b *Bot) ForwardMessage(chatID int64, fromChatID int64, messageID int) (tgbotapi.Message, error) {
	msg := tgbotapi.NewForward(chatID, fromChatID, messageID)
	return b.API.Send(msg)
}

// DeleteMessage deletes a message in a chat.
func (b *Bot) DeleteMessage(chatID int64, messageID int) error {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := b.API.Request(msg)
	return err
}

// IsSudo checks if a user is the owner or an authorized user.
func (b *Bot) IsSudo(userID int64) bool {
	if b.IsOwner(userID) {
		return true
	}
	for _, id := range b.Config.AuthUsers {
		if id == userID {
			return true
		}
	}
	return false
}

// IsOwner checks if a user is the owner.
func (b *Bot) IsOwner(userID int64) bool {
	return userID == b.Config.OwnerID
}
