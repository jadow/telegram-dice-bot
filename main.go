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
	configPath    = "config/config.json"
	dice4         = "4 side dice"
	dice6         = "6 side dice"
	dice8         = "8 side dice"
	dice10        = "10 side dice"
	dice12        = "12 side dice"
	dice20        = "20 side dice"
	help          = "telegram bot for rolling dice\n/open to open keyboard\n/close to close keyboard"
	openKeyboard  = "roll dice?"
	closeKeyboard = "closing keyboard"
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

func getName(update tgbotapi.Update) string {
	if update.Message.From.UserName != "" {
		return update.Message.From.UserName
	}
	return update.Message.From.FirstName
}

func getMessage(update tgbotapi.Update) string {
	if randomLimit, ok := menu[update.Message.Text]; ok {
		return fmt.Sprintf("%s rolled a %d", getName(update), rand.Intn(randomLimit)+1)
	}
	return update.Message.Text
}

func checkCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start", "help":
			msg.Text = help
		case "open":
			msg.Text = openKeyboard
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.Text = closeKeyboard
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
	} else {
		msg.Text = getMessage(update)
	}
	return msg
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

		msg := checkCommand(update)
		bot.Send(msg)
	}
}
