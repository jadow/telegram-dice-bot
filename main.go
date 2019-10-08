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

var diceStringInt = map[string]int{
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

func getName(user tgbotapi.User) string {
	if user.UserName != "" {
		return user.UserName
	}
	return user.FirstName
}

func getDice(randomLimit int) int {
	return rand.Intn(randomLimit) + 1
}

func getMessage(name string, random int) string {
	return fmt.Sprintf("%s rolled a %d", name, random)
}

func getQuery(update tgbotapi.Update) tgbotapi.InlineConfig {

	result4 := tgbotapi.NewInlineQueryResultArticle("1", dice4,
		getMessage(getName(*update.InlineQuery.From), getDice(diceStringInt[dice4])))
	result4.Description = dice4

	result6 := tgbotapi.NewInlineQueryResultArticle("2", dice6,
		getMessage(getName(*update.InlineQuery.From), getDice(diceStringInt[dice6])))
	result6.Description = dice6

	result8 := tgbotapi.NewInlineQueryResultArticle("3", dice8,
		getMessage(getName(*update.InlineQuery.From), getDice(diceStringInt[dice8])))
	result8.Description = dice8

	result10 := tgbotapi.NewInlineQueryResultArticle("4", dice10,
		getMessage(getName(*update.InlineQuery.From), getDice(diceStringInt[dice10])))
	result10.Description = dice10

	result12 := tgbotapi.NewInlineQueryResultArticle("5", dice12,
		getMessage(getName(*update.InlineQuery.From), getDice(diceStringInt[dice12])))
	result12.Description = dice12

	result20 := tgbotapi.NewInlineQueryResultArticle("6", dice20,
		getMessage(getName(*update.InlineQuery.From), getDice(diceStringInt[dice20])))
	result20.Description = dice20

	return tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       []interface{}{result4, result6, result8, result10, result12, result20},
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
			msg.Text = closeKeyboard
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
	} else if randomLimit, ok := diceStringInt[update.Message.Text]; ok {
		msg.Text = getMessage(getName(*update.Message.From), getDice(randomLimit))
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
