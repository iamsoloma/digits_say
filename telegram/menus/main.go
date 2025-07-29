package menus

import (
	"digits_say/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MakeMainMenu(user storage.User, MessageID int, chatId int64, message string) tgbotapi.MessageConfig {
	buttons := map[string]tgbotapi.KeyboardButton{
		"Start":         tgbotapi.NewKeyboardButton("/start"),
		"Consciousness": tgbotapi.NewKeyboardButton("Число сознания"),
		"Action":        tgbotapi.NewKeyboardButton("Число действия"),
		"Karma":         tgbotapi.NewKeyboardButton("Число кармы"),
		"Year":          tgbotapi.NewKeyboardButton("Число года"),
		"Month":         tgbotapi.NewKeyboardButton("Число месяца"),
		"PrivateDay":    tgbotapi.NewKeyboardButton("Личный день"),
		"SharedDay":     tgbotapi.NewKeyboardButton("Общий день"),
	}

	RegisterMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons["Start"],
			buttons["Consciousness"],
			buttons["Action"],
		),
		tgbotapi.NewKeyboardButtonRow(
			buttons["Karma"],
			buttons["Year"],
			buttons["Month"],
		),
		tgbotapi.NewKeyboardButtonRow(
			buttons["PrivateDay"],
			buttons["SharedDay"],
		),
	)

	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = RegisterMarkup
	msg.ReplyToMessageID = MessageID
	return msg

}
