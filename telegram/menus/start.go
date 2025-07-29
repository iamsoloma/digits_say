package menus

import (
	"fmt"

	"digits_say/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MakeStartMenu(user storage.User, update tgbotapi.Update) tgbotapi.MessageConfig {
	buttons := map[string]tgbotapi.InlineKeyboardButton{
		"Birthdate": tgbotapi.NewInlineKeyboardButtonData("Ввести дату рождения", "State=RegisterBirthdate"),
		"Email":     tgbotapi.NewInlineKeyboardButtonData("Ввести Email", "State=RegisterEmail"),
		"FullName":  tgbotapi.NewInlineKeyboardButtonData("Ввести полное имя", "State=RegisterFullName"),
		"Finished":  tgbotapi.NewInlineKeyboardButtonData("Завершить регистрацию", "State=RegisterFinished"),
	}

	RegisterMarkup := tgbotapi.NewInlineKeyboardMarkup()

	birth := ""
	if user.Birthdate == "" {
		birth = ""
	} else {
		birth = user.Birthdate[8:10] + "." + user.Birthdate[5:7] + "." + user.Birthdate[:4]
	}

	CurrentData := fmt.Sprintf(
		"Полное имя: %s\nEmail: %s\nДата рождения: %s",
		user.FullName, user.Email, birth,
	)

	Row1 := tgbotapi.NewInlineKeyboardRow()
	Row2 := tgbotapi.NewInlineKeyboardRow()
	Row3 := tgbotapi.NewInlineKeyboardRow()

	if user.FullName != "" && user.Birthdate != "" && user.FullName != "" && user.State["Register"] != "Finished" {
		Row1 = append(Row1, buttons["Birthdate"])
		Row1 = append(Row1, buttons["FullName"])
		Row2 = append(Row2, buttons["Email"])
		Row3 = append(Row3, buttons["Finished"])
		RegisterMarkup = tgbotapi.NewInlineKeyboardMarkup(Row1, Row2, Row3)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет "+user.Name+", вот что мне сейчас изветно о тебе.\n"+CurrentData+"\nМожешь изменить свою анкету.\n")
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = RegisterMarkup
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	} else if user.State["Register"] == "Finished" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет "+user.Name+", вот что мне сейчас изветно о тебе.\n"+CurrentData)
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	} else {
		if user.Birthdate == "" {
			Row1 = append(Row1, buttons["Birthdate"])
		}
		if user.Email == "" {
			Row1 = append(Row1, buttons["Email"])
		}
		if user.FullName == "" {
			Row2 = append(Row2, buttons["FullName"])
		}
		RegisterMarkup = tgbotapi.NewInlineKeyboardMarkup(Row1, Row2, Row3)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет "+user.Name+", вот что мне сейчас изветно о тебе.\n"+CurrentData+"\nЗаполни свою анкету.\n")
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = RegisterMarkup
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	}

}
