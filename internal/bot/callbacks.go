package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/StdBots/StdFeedback/internal/credit"
)

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	// Acknowledge the callback query to remove the loading state
	callback := tgbotapi.NewCallback(cb.ID, "")
	if _, err := b.API.Request(callback); err != nil {
		log.Printf("Failed to answer callback query: %v", err)
	}

	switch cb.Data {
	case "close":
		deleteMsg := tgbotapi.NewDeleteMessage(cb.Message.Chat.ID, cb.Message.MessageID)
		b.API.Request(deleteMsg)

	case "notif_toggle":
		// Toggle notification logic
		text := "Notifications toggled!"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Notifications: 🔕", "notif_toggle"),
			),
		)
		edit := tgbotapi.NewEditMessageTextAndMarkup(cb.Message.Chat.ID, cb.Message.MessageID, text, keyboard)
		b.API.Send(edit)

	case "help":
		text := credit.GetHelpMessage() + "\n\n" + credit.GetFooter()
		edit := tgbotapi.NewEditMessageTextAndMarkup(cb.Message.Chat.ID, cb.Message.MessageID, text, credit.GetInlineButtons())
		b.API.Send(edit)

	case "about":
		text := credit.GetAboutMessage() + "\n\n" + credit.GetFooter()
		edit := tgbotapi.NewEditMessageTextAndMarkup(cb.Message.Chat.ID, cb.Message.MessageID, text, credit.GetInlineButtons())
		b.API.Send(edit)

	case "start":
		text := credit.GetStartMessage() + "\n\n" + credit.GetFooter()
		edit := tgbotapi.NewEditMessageTextAndMarkup(cb.Message.Chat.ID, cb.Message.MessageID, text, credit.GetInlineButtons())
		b.API.Send(edit)

	case "stats":
		if b.IsSudo(cb.From.ID) {
			text := "Stats:\nUsers: 0\nMessages: 0"
			edit := tgbotapi.NewEditMessageText(cb.Message.Chat.ID, cb.Message.MessageID, text)
			b.API.Send(edit)
		}
	}
}
