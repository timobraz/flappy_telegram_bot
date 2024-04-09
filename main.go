package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7181024480:AAF1_1hvuUr3LfPUis4omoVZZXO80GCjTdk")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message{
			if strings.Contains("@catty_flappy_bot", update.Message.Text) {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://timobraz.github.io/cat-flappy/")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}
