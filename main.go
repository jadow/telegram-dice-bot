package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/jadow/Telegram-Dice-Bot/config"
	"github.com/tkanos/gonfig"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	configPath = "config/config.json"
	dice4      = "4 side dice"
	dice6      = "6 side dice"
	dice8      = "8 side dice"
	dice10     = "10 side dice"
	dice12     = "12 side dice"
	dice20     = "20 side dice"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(dice4),
		tgbotapi.NewKeyboardButton(dice6),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(dice8),
		tgbotapi.NewKeyboardButton(dice10),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(dice12),
		tgbotapi.NewKeyboardButton(dice20),
	),
)

func configure() config.Configuration {

	// server details
	c := &config.Configuration{}
	err := gonfig.GetConf(configPath, c)
	if err != nil {
		log.Println("configuration error. unable to configure server")
	}

	return *c
}

func main() {
	config := configure()

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Print(".")
	for update := range updates {
		if update.Message == nil {
			continue
		}
		reply := ""
		randomLimit := 0
		switch update.Message.Text {
		case dice4:
			randomLimit = 4
		case dice6:
			randomLimit = 6
		case dice8:
			randomLimit = 8
		case dice10:
			randomLimit = 10
		case dice12:
			randomLimit = 12
		case dice20:
			randomLimit = 20
		}

		if randomLimit > 0 {
			reply = fmt.Sprintf("%s rolled a %d", update.Message.From.UserName, rand.Intn(randomLimit)+1)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		bot.Send(msg)
	}

}
