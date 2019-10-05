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

var menu = map[string]int{
	dice4:  4,
	dice6:  6,
	dice8:  8,
	dice10: 10,
	dice12: 12,
	dice20: 20,
}

func configure() config.Configuration {

	// server details
	c := &config.Configuration{}
	err := gonfig.GetConf(configPath, c)
	if err != nil {
		log.Println("configuration error. unable to configure server")
	}

	return *c
}

func getMessage(update tgbotapi.Update) string {
	if randomLimit, ok := menu[update.Message.Text]; ok {
		return fmt.Sprintf("%s rolled a %d", update.Message.From.UserName, rand.Intn(randomLimit)+1)
	}
	return ""
}

func controlMenu(message string) interface{} {
	switch message {
	case "open":
		return numericKeyboard
	case "close":
		return tgbotapi.NewRemoveKeyboard(true)
	}
	return nil
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

		message := update.Message.Text
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, getMessage(update))
		msg.ReplyMarkup = controlMenu(message)

		bot.Send(msg)
	}

}
