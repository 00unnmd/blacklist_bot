package handlers

import (
	"blacklist_bot/internal/models"
	"blacklist_bot/utils/validation"
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
)

func (h *BotHandler) addUserPhoneNumber(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))
	h.bannedUser.PhoneNumber = ""

	err := h.EditBotMessage(c.Chat().ID, "➕ Добавление преподавателя \n Шаг 1. Введите номер телефона", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		normalizedPhone, err := validation.ValidateAndNormalizePhone(ctx.Text())
		if err != nil {
			errMsg := fmt.Sprintf("❌ Ошибка: %s\nВведите номер еще раз.", err)
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, errMsg)
		}

		h.bannedUser.PhoneNumber = normalizedPhone
		return h.addUserFullName(c)
	})

	return nil
}

func (h *BotHandler) addUserFullName(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("⬅️ Назад", "add_user_phone_number")
	markup.Inline(markup.Row(btnCancel, btnPrev))
	h.bannedUser.FullName = ""

	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, "➕ Добавление преподавателя \n Шаг 2. Введите ФИО", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		err := validation.ValidateDescriptionStr(ctx.Text())
		if err != nil {
			errMsg := fmt.Sprintf("❌ Ошибка: %s\nВведите ФИО еще раз.", err)
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, errMsg)
		}

		h.bannedUser.FullName = ctx.Text()
		return h.addUserDescription(c)
	})

	return nil
}

func (h *BotHandler) addUserDescription(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("⬅️ Назад", "add_user_full_name")
	markup.Inline(markup.Row(btnCancel, btnPrev))
	h.bannedUser.Description = ""

	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, "➕ Добавление преподавателя \n Шаг 3. Введите описание", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		err := validation.ValidateDescriptionStr(ctx.Text())
		if err != nil {
			errMsg := fmt.Sprintf("❌ Ошибка: %s\nВведите описание еще раз.", err)
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, errMsg)
		}

		h.bannedUser.Description = ctx.Text()
		return h.addUserBirthday(c)
	})

	return nil
}

func (h *BotHandler) addUserBirthday(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("⬅️ Назад", "add_user_description")
	btnSkip := markup.Data("➡️ Пропустить", "add_user_city")
	markup.Inline(markup.Row(btnCancel, btnPrev, btnSkip))
	h.bannedUser.BirthDay = ""

	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, "➕ Добавление преподавателя \n Шаг 4. Введите дату рождения в формате 01.01.2000", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		err := validation.ValidateBirthdayStr(ctx.Text())
		if err != nil {
			errMsg := fmt.Sprintf("❌ Ошибка: %s\nВведите дату рождения еще раз.", err)
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, errMsg)
		}

		h.bannedUser.BirthDay = ctx.Text()
		return h.addUserCity(c)
	})

	return nil
}

func (h *BotHandler) addUserCity(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("⬅️ Назад", "add_user_birthday")
	btnSkip := markup.Data("➡️ Пропустить", "add_user_school_format")
	markup.Inline(markup.Row(btnCancel, btnPrev, btnSkip))
	h.bannedUser.City = ""

	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, "➕ Добавление преподавателя \n Шаг 5. Введите город", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		err := validation.ValidateCityStr(ctx.Text())
		if err != nil {
			errMsg := fmt.Sprintf("❌ Ошибка: %s\nВведите город еще раз.", err)
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, errMsg)
		}

		h.bannedUser.City = ctx.Text()
		return h.addUserSchoolFormat(c)
	})

	return nil
}

func (h *BotHandler) addUserSchoolFormat(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnOfflineFormat := markup.Data("🏫 Оффлайн", "school_offline")
	btnOnlineFormat := markup.Data("🌐 Онлайн", "school_online")
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("⬅️ Назад", "add_user_city")
	btnSkip := markup.Data("➡️ Пропустить", "add_user_confirmation")
	markup.Inline(
		markup.Row(btnOfflineFormat, btnOnlineFormat),
		markup.Row(btnCancel, btnPrev, btnSkip),
	)
	h.bannedUser.SchoolFormat = ""

	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, "➕ Добавление преподавателя \n Шаг 6. Выберите формат школы (Оффлайн/Онлайн)", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(&btnOfflineFormat, func(ctx telebot.Context) error {
		h.bannedUser.SchoolFormat = "Оффлайн"
		return h.addUserConfirmation(ctx)
	})

	h.bot.Handle(&btnOnlineFormat, func(ctx telebot.Context) error {
		h.bannedUser.SchoolFormat = "Онлайн"
		return h.addUserConfirmation(ctx)
	})

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		err := h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, "❌ Ошибка: формат школы необходимо выбрать.")
		if err != nil {
			return err
		}

		return h.addUserSchoolFormat(c)
	})

	return nil
}

func (h *BotHandler) addUserConfirmation(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnConfirm := markup.Data("✅ Сохранить", "save_user")
	markup.Inline(markup.Row(btnCancel, btnConfirm))

	strF := fmt.Sprintf("➕ Добавление преподавателя \n"+
		"Проверьте информацию и подтвердите добавление. \n"+
		"Номер телефона: +%s \n"+
		"ФИО: %s \n"+
		"Описание: %s \n"+
		"Дата рождения: %s \n"+
		"Город: %s \n"+
		"Формат школы: %s \n",
		h.bannedUser.PhoneNumber, h.bannedUser.FullName, h.bannedUser.Description, h.bannedUser.BirthDay, h.bannedUser.City, h.bannedUser.SchoolFormat)
	err := h.SendAndTrack(c.Recipient(), c.Chat().ID, strF, markup)
	if err != nil {
		return err
	}

	return nil
}

func (h *BotHandler) saveBannedUser(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnRepeat := markup.Data("➕ Добавить преподавателя", "add_user_phone_number")
	markup.Inline(markup.Row(btnCancel, btnRepeat))

	if err := h.db.AddBannedUser(h.bannedUser); err != nil {
		return h.SendAndTrack(c.Recipient(), c.Chat().ID, "❌ Ошибка при добавлении преподавателя: "+err.Error(), markup)
	}

	return h.SendAndTrack(c.Recipient(), c.Chat().ID, "✅ Пользователь успешно добавлен!", markup)
}

func (h *BotHandler) findUserHandler(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))

	err := h.EditBotMessage(c.Chat().ID, "🔍 Поиск преподавателя.\nВведите номер телефона или ФИО для поиска", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		btnRepeat := markup.Data("🔍 Найти преподавателя", "find_user")
		markup.Inline(markup.Row(btnCancel, btnRepeat))
		input := ctx.Text()

		isPhoneNumber := validation.IsPhoneNumber(input)
		var users []models.BannedUser
		var err error

		if isPhoneNumber {
			normalizedPhone, errV := validation.ValidateAndNormalizePhone(input)
			if errV != nil {
				errMsg := fmt.Sprintf("❌ Ошибка: %s\nВведите номер еще раз.", errV)
				return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, errMsg)
			}
			users, err = h.db.FindBannedUserByPhone(normalizedPhone)
		} else {
			users, err = h.db.FindBannedUserByName(input)
		}

		if err != nil {
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, "❌ Ошибка при поиске преподавателя: "+err.Error(), markup)
		}

		if len(users) < 1 {
			var searchInput string
			if isPhoneNumber {
				searchInput = "номером телефона"
			} else {
				searchInput = "ФИО"
			}

			msg := fmt.Sprintf("❌ Пользователь с таким %s не найден", searchInput)
			return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, msg, markup)
		}

		var usersBuilder strings.Builder
		for _, item := range users {
			usersBuilder.WriteString(fmt.Sprintf(
				"\n\nФИО: %s\nОписание: %s\nГород: %s\nФормат школы: %s",
				item.FullName, item.Description, item.City, item.SchoolFormat,
			))
		}
		usersStr := usersBuilder.String()

		return h.SendAndTrack(ctx.Recipient(), ctx.Chat().ID, fmt.Sprintf("🔍 Пользователь найден!%s", usersStr), markup)
	})

	return nil
}
