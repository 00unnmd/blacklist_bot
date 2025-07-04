package handlers

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

func (h *BotHandler) addUserPhoneNumber(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))
	h.bannedUser.PhoneNumber = ""

	err := c.Send("➕ Добавление пользователя \n Шаг 1. Введите номер телефона", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		// validate phone number
		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: Номер телефона слишком короткий.")
		}
		h.bannedUser.PhoneNumber = ctx.Text()
		return h.addUserFullName(c)
	})

	return nil
}

func (h *BotHandler) addUserFullName(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("Назад", "add_user_phone_number")
	markup.Inline(markup.Row(btnCancel, btnPrev))
	h.bannedUser.FullName = ""

	err := c.Send("➕ Добавление пользователя \n Шаг 2. Введите ФИО", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		// validate fullname
		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: ФИО слишком короткое.")
		}
		h.bannedUser.FullName = ctx.Text()
		return h.addUserDescription(c)
	})

	return nil
}

func (h *BotHandler) addUserDescription(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("Назад", "add_user_full_name")
	markup.Inline(markup.Row(btnCancel, btnPrev))
	h.bannedUser.Description = ""

	err := c.Send("➕ Добавление пользователя \n Шаг 3. Введите описание", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		// validate description
		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: описание слишком короткое.")
		}
		h.bannedUser.Description = ctx.Text()
		return h.addUserBirthday(c)
	})

	return nil
}

func (h *BotHandler) addUserBirthday(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("Назад", "add_user_description")
	btnSkip := markup.Data("Пропустить", "add_user_city")
	markup.Inline(markup.Row(btnCancel, btnPrev, btnSkip))
	h.bannedUser.BirthDay = ""

	err := c.Send("➕ Добавление пользователя \n Шаг 4. Введите дату рождения", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		// validate birthday
		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: дата рождения слишком короткая.")
		}
		h.bannedUser.BirthDay = ctx.Text()
		return h.addUserCity(c)
	})

	return nil
}

func (h *BotHandler) addUserCity(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("Назад", "add_user_birthday")
	btnSkip := markup.Data("Пропустить", "add_user_school_format")
	markup.Inline(markup.Row(btnCancel, btnPrev, btnSkip))
	h.bannedUser.City = ""

	err := c.Send("➕ Добавление пользователя \n Шаг 5. Введите город", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		// validate city (no need)
		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: название города слишком короткое.")
		}
		h.bannedUser.City = ctx.Text()
		return h.addUserSchoolFormat(c)
	})

	return nil
}

func (h *BotHandler) addUserSchoolFormat(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnPrev := markup.Data("Назад", "add_user_city")
	btnSkip := markup.Data("Пропустить", "add_user_confirmation")
	markup.Inline(markup.Row(btnCancel, btnPrev, btnSkip))
	h.bannedUser.SchoolFormat = ""

	err := c.Send("➕ Добавление пользователя \n Шаг 6. Введите формат школы (оффлайн/онлайн)", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		// validate city (no need)
		if len(ctx.Text()) < 1 {
			ctx.Send("❌ Ошибка: формат школы слишком короткий.")
		}
		h.bannedUser.SchoolFormat = ctx.Text()
		return h.addUserConfirmation(c)
	})

	return nil
}

func (h *BotHandler) addUserConfirmation(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnConfirm := markup.Data("✅ Сохранить", "save_user")
	markup.Inline(markup.Row(btnCancel, btnConfirm))

	strF := fmt.Sprintf("➕ Добавление пользователя \n"+
		"Новый пользователь: \n"+
		"Номер телефона: %s \n"+
		"ФИО: %s \n"+
		"Описание: %s \n"+
		"Дата рождения: %s \n"+
		"Город: %s \n"+
		"Формат школы: %s \n",
		h.bannedUser.PhoneNumber, h.bannedUser.FullName, h.bannedUser.Description, h.bannedUser.BirthDay, h.bannedUser.City, h.bannedUser.SchoolFormat)
	err := c.Send(strF, markup)
	if err != nil {
		return err
	}

	return nil
}

func (h *BotHandler) saveBannedUser(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	btnRepeat := markup.Data("➕ Добавить пользователя", "add_user_phone_number")
	markup.Inline(markup.Row(btnCancel, btnRepeat))

	if err := h.db.AddBannedUser(h.bannedUser); err != nil {
		return c.Send("❌ Ошибка при добавлении пользователя: "+err.Error(), markup)
	}

	return c.Send("✅ Пользователь успешно добавлен!", markup)
}

func (h *BotHandler) findUserHandler(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	btnCancel := markup.Data("Ⓜ️ В главное меню", "main_menu")
	markup.Inline(markup.Row(btnCancel))

	err := c.Send("🔍 Поиск пользователя. Введите номер телефона для поиска:", markup)
	if err != nil {
		return err
	}

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		btnRepeat := markup.Data("🔍 Найти пользователя", "find_user")
		markup.Inline(markup.Row(btnCancel, btnRepeat))
		phoneNumber := ctx.Text()

		user, err := h.db.FindBannedUser(phoneNumber)
		if err != nil {
			return ctx.Send("❌ Ошибка при поиске пользователя: "+err.Error(), markup)
		}

		if user == nil {
			return ctx.Send("❌ Пользователь с таким номером телефона не найден", markup)
		}

		return ctx.Send(fmt.Sprintf(
			"🔍 Пользователь найден! \nНомер: %s\nФИО: %s\nОписание: %s",
			user.PhoneNumber, user.FullName, user.Description,
		), markup)
	})

	return nil
}
