package filters

import (
	"strings"

	"github.com/go-telegram/bot/models"
)

func IsGreeting(upd *models.Update) bool {
	if upd.Message == nil {
		return false
	}
	return strings.Contains(strings.ToLower(upd.Message.Text), "привет")
}

func IsGroup(upd *models.Update) bool {
	if upd.Message.Chat.Type == models.ChatTypeGroup || upd.Message.Chat.Type == models.ChatTypeSupergroup {
		return true
	}
	return false
}

func IsAboutAI(upd *models.Update) bool {
	if strings.Contains(strings.ToLower(upd.Message.Text), "ии") ||
		strings.Contains(strings.ToLower(upd.Message.Text), "искусственный интеллект") ||
		strings.Contains(strings.ToLower(upd.Message.Text), "бизнес") {
		return true
	}
	return false
}
