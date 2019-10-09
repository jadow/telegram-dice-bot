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
	configPath     = "config/config.json"
	dice4          = "4 sided dice"
	dice6          = "6 sided dice"
	dice8          = "8 sided dice"
	dice10         = "10 sided dice"
	dice12         = "12 sided dice"
	dice20         = "20 sided dice"
	help           = "telegram bot for rolling dice\n/open to open keyboard\n/close to close keyboard"
	openKeyboard   = "roll dice?"
	removeKeyboard = "removing keyboard"
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

var diceStringInt = map[string]int{
	dice4:  4,
	dice6:  6,
	dice8:  8,
	dice10: 10,
	dice12: 12,
	dice20: 20,
}

var mapIntString = map[int]string{
	1: "1",
	2: "2",
	3: "3",
	4: "4",
	5: "5",
	6: "6",
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

func getName(user tgbotapi.User) string {
	if user.UserName != "" {
		return user.UserName
	}
	return user.FirstName
}

func getDice(randomLimit int) int {
	return rand.Intn(randomLimit) + 1
}

func getDiceMessage(name string, random int) string {
	return fmt.Sprintf("%s rolled a %d", name, random)
}

func getDiceRollMessage(name string, dice string, random int) string {
	return fmt.Sprintf("%s rolled a %s and got a %d ", name, dice, random)
}

//identify using title
func getQuery(update tgbotapi.Update) tgbotapi.InlineConfig {
	results := []interface{}{}
	for i, v := range diceStringInt {
		result := tgbotapi.NewInlineQueryResultArticle(i, i,
			getDiceRollMessage(getName(*update.InlineQuery.From), i, getDice(v)))
		result.Description = i
		results = append(results, result)
	}

	return tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       results,
	}
}

func getCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start", "help":
			msg.Text = help
		case "open":
			msg.Text = openKeyboard
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.Text = removeKeyboard
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
	} else if randomLimit, ok := diceStringInt[update.Message.Text]; ok {
		msg.Text = getDiceMessage(getName(*update.Message.From), getDice(randomLimit))
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

		if update.InlineQuery != nil {
			res := getQuery(update)
			bot.AnswerInlineQuery(res)
		} else if update.Message != nil {
			msg := getCommand(update)
			bot.Send(msg)
		}
		//else continue
	}
}
