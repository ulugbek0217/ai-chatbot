package filters

import (
	"strings"

	"github.com/go-telegram/bot/models"
)

func IsGroup(upd *models.Update) bool {
	if upd.Message.Chat.Type == models.ChatTypeGroup || upd.Message.Chat.Type == models.ChatTypeSupergroup {
		return true
	}
	return false
}

func AboutAI(upd *models.Update) bool {
	if strings.Contains(strings.ToLower(upd.Message.Text), "ии") ||
		strings.Contains(strings.ToLower(upd.Message.Text), "искусственный интеллект") {
		return true
	}
	return false
}
