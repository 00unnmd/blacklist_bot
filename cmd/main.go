package main

import (
	"blacklist_bot/internal/database"
	"blacklist_bot/internal/handlers"
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	pref := telebot.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	handler := handlers.New(bot, db)
	handler.SetupHandlers()

	go bot.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down bot...")
}
