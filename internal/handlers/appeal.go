package handlers

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

func (h *BotHandler) addAppealHandler(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))

	link := `<a href="https://t.me/unnmd00">администратору</a>`
	msgStr := fmt.Sprintf("📝 Напишите текст обращения %s", link)
	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, msgStr, &telebot.SendOptions{ParseMode: telebot.ModeHTML, ReplyMarkup: markup})
	if err != nil {
		return err
	}

	return nil
}
