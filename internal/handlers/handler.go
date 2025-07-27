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
	tracker    *MessageTracker
}

func New(bot *telebot.Bot, db *database.Database) *BotHandler {
	return &BotHandler{
		bot:        bot,
		db:         db,
		bannedUser: models.BannedUser{},
		tracker:    NewMessageTracker(),
	}
}

func (h *BotHandler) SetupHandlers() {
	h.bot.Use(h.tracker.TrackMessages())
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

// SendAndTrack added for tracking bot messages (user messages tracked with TrackMessages)
func (h *BotHandler) SendAndTrack(to telebot.Recipient, chatId int64, text string, options ...interface{}) error {
	msg, err := h.bot.Send(to, text, options...)
	if err == nil {
		h.tracker.TrackMessage(chatId, msg.ID)
	}
	return err
}

func (h *BotHandler) showMainMenu(c telebot.Context) error {
	h.tracker.ClearChatHistory(h.bot, c.Chat().ID)

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
		return h.SendAndTrack(c.Recipient(), c.Chat().ID, "Выберите пункт из меню")
	})

	return h.SendAndTrack(c.Recipient(), c.Chat().ID, "Ⓜ️ Главное меню", markup)
}
