package handlers

import (
	"blacklist_bot/internal/database"
	"blacklist_bot/internal/models"
	"gopkg.in/telebot.v3"
	"log"
	"strings"
)

type BotHandler struct {
	bot        *telebot.Bot
	db         *database.Database
	bannedUser models.BannedUser
}

func New(bot *telebot.Bot, db *database.Database) *BotHandler {
	return &BotHandler{
		bot:        bot,
		db:         db,
		bannedUser: models.BannedUser{},
	}
}

func (h *BotHandler) SetupHandlers() {
	h.bot.Handle("/start", h.showMainMenu)

	h.bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		if c.Message() != nil {
			if err := c.Delete(); err != nil {
				log.Println("Не удалось удалить сообщение меню:", err)
			}
		}

		// for some reason c.Callback().Data str starts with \n separator
		data := strings.TrimSpace(c.Callback().Data)
		switch data {
		case "add_user_phone_number":
			return h.addUserPhoneNumber(c)
		case "add_user_full_name":
			return h.addUserFullName(c)
		case "add_user_description":
			return h.addUserDescription(c)
		case "add_user_birthday":
			return h.addUserBirthday(c)
		case "add_user_city":
			return h.addUserCity(c)
		case "add_user_school_format":
			return h.addUserSchoolFormat(c)
		case "add_user_confirmation":
			return h.addUserConfirmation(c)
		case "save_user":
			return h.saveBannedUser(c)
		case "find_user":
			return h.findUserHandler(c)
		case "add_appeal":
			return h.addAppealHandler(c)
		case "main_menu":
			return h.showMainMenu(c)
		}

		return nil
	})
}

func (h *BotHandler) showMainMenu(c telebot.Context) error {
	h.bannedUser = models.BannedUser{}
	markup := &telebot.ReplyMarkup{}
	btnAddUser := markup.Data("➕ Добавить пользователя", "add_user_phone_number")
	btnFindUser := markup.Data("🔍 Найти пользователя", "find_user")
	btnAddAppeal := markup.Data("📝 Оставить обращение", "add_appeal")

	markup.Inline(
		markup.Row(btnAddUser, btnFindUser),
		markup.Row(btnAddAppeal),
	)

	h.bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		return c.Send("Выберите пункт из меню")
	})

	return c.Send("Ⓜ️ Главное меню", markup)
}
