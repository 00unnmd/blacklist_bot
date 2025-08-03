package handlers

import (
	"blacklist_bot/internal/database"
	"blacklist_bot/internal/models"
	"blacklist_bot/utils"
	"gopkg.in/telebot.v3"
	"strings"
)

type BotHandler struct {
	bot        *telebot.Bot
	db         *database.Database
	bannedUser models.BannedUser
	tracker    *utils.MessageTracker
}

func New(bot *telebot.Bot, db *database.Database) *BotHandler {
	return &BotHandler{
		bot:        bot,
		db:         db,
		bannedUser: models.BannedUser{},
		tracker:    utils.NewMessageTracker(),
	}
}

func (h *BotHandler) SetupHandlers() {
	h.bot.Use(h.tracker.TrackMessages())
	h.bot.Handle("/start", h.showStart)
	h.bot.Handle("/about", h.showAbout)
	h.bot.Handle("/stop", h.showStop)

	h.bot.Handle(telebot.OnText, func(c telebot.Context) error {
		startMsg := "🚀 Выберите /start чтобы открыть главное меню.\n\n🤖 Выберите /about чтобы посмотреть информацию о боте."
		return h.SendAndTrack(c.Recipient(), c.Chat().ID, startMsg, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
	})

	h.bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
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
			return h.showStart(c)
		default:
			return h.showStart(c)
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

func (h *BotHandler) EditBotMessage(chatId int64, what interface{}, opts ...interface{}) error {
	msgID, err := h.tracker.GetLastBotMessageId(chatId)
	if err != nil {
		return err
	}

	_, err = h.bot.Edit(&telebot.Message{
		Chat: &telebot.Chat{ID: chatId},
		ID:   msgID,
	}, what, opts...)

	return err
}

func (h *BotHandler) showStart(c telebot.Context) error {
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
		h.tracker.ClearChatHistory(h.bot, c.Chat().ID)
		return h.SendAndTrack(c.Recipient(), c.Chat().ID, "Ⓜ️ Выберите пункт из меню", markup)
	})

	return h.SendAndTrack(c.Recipient(), c.Chat().ID, "Ⓜ️ Главное меню", markup)
}

func (h *BotHandler) showAbout(c telebot.Context) error {
	aboutMsg := "👋 <b>Добро пожаловать!</b>\n\n" +
		"Этот бот позволяет работать с черным списком преподавателей.\n\n" +
		"Основные возможности:\n" +
		"• Добавление преподавателей в ЧС\n" +
		"• Поиск преподавателей по базе\n" +
		"• Возможность оставить обращение администратору\n\n" +
		"🚀 Выберите /start чтобы открыть главное меню."

	return h.SendAndTrack(c.Recipient(), c.Chat().ID, aboutMsg, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
}

func (h *BotHandler) showStop(c telebot.Context) error {
	h.tracker.ClearChatHistory(h.bot, c.Chat().ID)
	h.bot.Handle(telebot.OnText, func(c telebot.Context) error {
		startMsg := "🚀 Выберите /start чтобы открыть главное меню.\n\n🤖 Выберите /about чтобы посмотреть информацию о боте."
		return h.SendAndTrack(c.Recipient(), c.Chat().ID, startMsg, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
	})

	stopMsg := "🛑 Работа с ботом завершена.\n\n🚀 Чтобы начать снова, выберите /start."
	return h.SendAndTrack(c.Recipient(), c.Chat().ID, stopMsg)
}
