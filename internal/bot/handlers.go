package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/StdBots/StdFeedback/internal/credit"
)

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	// Add user to DB logic here
	// Notify LOG_CHANNEL if new user

	text := credit.GetStartMessage() + "\n\n" + credit.GetFooter()
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = credit.GetInlineButtons()

	if _, err := b.API.Send(reply); err != nil {
		log.Printf("Error sending start message: %v", err)
	}
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	text := credit.GetHelpMessage() + "\n\n" + credit.GetFooter()
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = credit.GetInlineButtons()

	b.API.Send(reply)
}

func (b *Bot) handleAbout(msg *tgbotapi.Message) {
	text := credit.GetAboutMessage() + "\n\n" + credit.GetFooter()
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = credit.GetInlineButtons()

	b.API.Send(reply)
}

func (b *Bot) handleSettings(msg *tgbotapi.Message) {
	text := "Settings:\nToggle notifications below:"
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Notifications: 🔔", "notif_toggle"),
		),
	)
	reply.ReplyMarkup = keyboard

	b.API.Send(reply)
}

func (b *Bot) handleDonate(msg *tgbotapi.Message) {
	text := "Donate to support the developer!"
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

// Admin Commands
func (b *Bot) handleBan(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	// Parsing logic for /ban user_id duration reason
	text := "User banned."
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

func (b *Bot) handleUnban(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "User unbanned."
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

func (b *Bot) handleStats(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "Bot Stats:\nUsers: 0\nMessages: 0\nClones: 0"
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

func (b *Bot) handleBroadcast(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "Broadcast started."
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

func (b *Bot) handleUsers(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "Total Users: 0"
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

// Clone Commands
func (b *Bot) handleClone(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "Starting clone bot..."
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

func (b *Bot) handleDeleteClone(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "Clone bot deleted."
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}

func (b *Bot) handleClones(msg *tgbotapi.Message) {
	if !b.IsSudo(msg.From.ID) {
		return
	}
	text := "Active Clones:\nNone"
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	b.API.Send(reply)
}
