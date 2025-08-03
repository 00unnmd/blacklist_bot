package utils

import (
	"gopkg.in/telebot.v3"
	"sync"
	"time"
)

type MessageTracker struct {
	mu       sync.Mutex
	messages map[int64][]int
}

func NewMessageTracker() *MessageTracker {
	return &MessageTracker{
		messages: make(map[int64][]int),
	}
}

func (mt *MessageTracker) TrackMessage(chatID int64, messageID int) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	mt.messages[chatID] = append(mt.messages[chatID], messageID)
}

func (mt *MessageTracker) TrackMessages() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Message() != nil {
				mt.mu.Lock()
				mt.messages[c.Chat().ID] = append(mt.messages[c.Chat().ID], c.Message().ID)
				mt.mu.Unlock()
			}
			return next(c)
		}
	}
}

func (mt *MessageTracker) GetLastBotMessageId(chatId int64) (int, error) {
	if messages, exist := mt.messages[chatId]; exist && len(messages) > 0 {
		return messages[len(messages)-1], nil
	}
	return 0, nil
}

func (mt *MessageTracker) ClearChatHistory(bot *telebot.Bot, chatID int64) {
	mt.mu.Lock()
	messageIDs, exists := mt.messages[chatID]
	if !exists || len(messageIDs) <= 1 {
		mt.mu.Unlock()
		return
	}

	messagesToDelete := messageIDs[:len(messageIDs)-1]
	mt.messages[chatID] = messageIDs[len(messageIDs)-1:]
	mt.mu.Unlock()

	go func(ids []int) {
		for i, msgID := range ids {
			_ = bot.Delete(&telebot.Message{
				Chat: &telebot.Chat{ID: chatID},
				ID:   msgID,
			})

			if i%20 == 0 && i > 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}(messagesToDelete)
}
