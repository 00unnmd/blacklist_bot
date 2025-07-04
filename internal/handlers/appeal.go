package handlers

import (
	"blacklist_bot/internal/models"
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
)

func (h *BotHandler) addAppealHandler(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))

	err := c.Send("📝 Напишите текст обращения", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		btnRepeat := markup.Data("📝 Оставить обращение", "add_appeal")
		markup.Inline(markup.Row(btnCancel, btnRepeat))

		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: пустой текст обращения.", markup)
		}

		appeal := models.Appeal{
			Question:   ctx.Text(),
			Initiator:  ctx.Sender().Username,
			IsAnswered: false,
		}

		if err := h.db.AddAppeal(appeal); err != nil {
			return ctx.Send("❌ Ошибка при добавлении обращения: "+err.Error(), markup)
		}

		return ctx.Send("✅ Обращение успешно добавлено!", markup)
	})

	return nil
}

func (h *BotHandler) listAppealsHandler(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))

	appeals, err := h.db.GetUnansweredAppeals()
	if err != nil {
		return c.Send("Ошибка при получении списка обращений: "+err.Error(), markup)
	}

	if len(appeals) == 0 {
		return c.Send("Нет неотвеченных обращений", markup)
	}

	var builder strings.Builder
	builder.WriteString("📋 Список обращений \n\n")

	for _, appeal := range appeals {
		builder.WriteString(fmt.Sprintf(
			"Обращение: №%d \n"+
				"Инициатор: @%s \n"+
				"Вопрос: %s \n\n",
			appeal.ID,
			appeal.Initiator,
			appeal.Question,
		))
	}

	err = c.Send(builder.String(), markup, telebot.ModeHTML)
	if err != nil {
		return err
	}

	return nil
}
