package telegram

import (
	"digits_say/digits"
	"digits_say/storage"
	"fmt"
	"log"
	"net/mail"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type Config struct {
	ApiToken string
	Debug    bool
	Timeout  int
	DBConfig storage.DBConfig
}

type TelegramListener struct {
	*Config
	bot          *tgbotapi.BotAPI
	updateConfig *tgbotapi.UpdateConfig

	storage *storage.Storage
}

func NewListener(config Config) (TelegramBot *TelegramListener, err error) {
	storage := storage.Storage{
		DBConfig: config.DBConfig,
	}
	err = storage.ConnectToSurreal()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SurrealDB: %w", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.ApiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}
	if config.Debug == true {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = config.Timeout

	return &TelegramListener{
		Config:       &config,
		bot:          bot,
		updateConfig: &updateConfig,
		storage:      &storage,
	}, nil
}

func (l *TelegramListener) Start() {
	log.Println("Online as a Telegram bot @" + l.bot.Self.UserName)

	for update := range l.bot.GetUpdatesChan(*l.updateConfig) {
		if update.Message != nil {
			if update.Message.IsCommand() {
				l.HandleCommad(update)
			} else if update.Message.Text != "" {
				l.HandleText(update)
			}
		}
		if update.CallbackQuery != nil {
			l.HandleCallbacks(update.CallbackQuery)
		}
	}
}

func (l *TelegramListener) HandleCommad(update tgbotapi.Update) {
	if update.Message.Text == "/start" {
		user, exist, err := l.storage.GetUserByID(fmt.Sprintf("tg%d", update.Message.From.ID))
		if err != nil {
			log.Println("Error getting user by Telegram ID: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при получении пользователя. Попробуй позже.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}

		if !exist {
			userData := storage.User{
				ID:           models.RecordID{Table: "Users", ID: fmt.Sprintf("tg%d", update.Message.From.ID)},
				UserName:     update.Message.From.UserName,
				Name:         update.Message.From.FirstName,
				Surname:      update.Message.From.LastName,
				LanguageCode: update.Message.From.LanguageCode,
			}
			fmt.Printf("Registering new user with data: %#v\n", userData)
			registeredUser, err := l.storage.RegisterNewUser(userData)
			if err != nil {
				log.Println("Error registering new user: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при регистрации. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			user = registeredUser
			exist = true

			msg := MakeStartMenu(*user, update)
			l.bot.Send(msg)
			return
		}
		if exist {
			msg := MakeStartMenu(*user, update)
			l.bot.Send(msg)
			return
		}

	}
}

func (l *TelegramListener) HandleCallbacks(callback *tgbotapi.CallbackQuery) {
	if callback.Data == "State=RegisterBirthdate" {
		user, _, err := l.storage.GetUserByID(fmt.Sprintf("tg%d", callback.From.ID))
		if err != nil {
			log.Println("Error getting user by Telegram ID: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при получении пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.State = "RegisterBirthdate"
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user state: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введи свою дату рождения в формате ДД.ММ.ГГГГ:")
		msg.ReplyToMessageID = callback.Message.MessageID
		l.bot.Send(msg)
	} else if callback.Data == "State=RegisterEmail" {
		user, _, err := l.storage.GetUserByID(fmt.Sprintf("tg%d", callback.From.ID))
		if err != nil {
			log.Println("Error getting user by Telegram ID: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при получении пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.State = "RegisterEmail"
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user state: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введи свой Email:")
		msg.ReplyToMessageID = callback.Message.MessageID
		l.bot.Send(msg)
	} else if callback.Data == "State=RegisterFullName" {
		user, _, err := l.storage.GetUserByID(fmt.Sprintf("tg%d", callback.From.ID))
		if err != nil {
			log.Println("Error getting user by Telegram ID: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при получении пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.State = "RegisterFullName"
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user state: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введи свое полное имя латинскими буквами как в загране или банковской карте(Nadezda, Vitaliy):")
		msg.ReplyToMessageID = callback.Message.MessageID
		l.bot.Send(msg)
	}
}

func (l *TelegramListener) HandleText(update tgbotapi.Update) {
	user, exist, err := l.storage.GetUserByID(fmt.Sprintf("tg%d", update.Message.From.ID))
	if err != nil {
		log.Println("Error getting user by Telegram ID: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при получении пользователя. Попробуй позже.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		return
	}

	if !exist {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не зарегистрирован. Напиши /start, чтобы начать.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		return
	}

	if user.State == "RegisterBirthdate" {
		birthdate, err := time.Parse("02.01.2006", update.Message.Text)
		if err != nil {
			log.Println("Error parsing birthdate: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты. Пожалуйста, введи дату в формате ДД.ММ.ГГГГ.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.Birthdate = birthdate.Format("2006-01-02") //models.CustomDateTime{Time: birthdate}
		user.State = ""
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user birthdate: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при обновлении даты рождения. Попробуй позже.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Твоя дата рождения успешно сохранена.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		msg = MakeStartMenu(*user, update)
		l.bot.Send(msg)

		if user.Birthdate != "" && user.Email != "" && user.FullName != "" {
			consciousnessNumber, err := digits.GetConsciousnessNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating consciousness number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего сознания. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твой номер сознания: %d", consciousnessNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
		}
	} else if user.State == "RegisterEmail" {
		_, err := mail.ParseAddress(update.Message.Text)
		if err != nil {
			log.Println("Error parsing email: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат Email. Пожалуйста, введи корректный Email.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.Email = update.Message.Text
		user.State = ""
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user email: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при обновлении Email. Попробуй позже.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Твой Email успешно сохранён.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		msg = MakeStartMenu(*user, update)
		l.bot.Send(msg)

		if user.Birthdate != "" && user.Email != "" && user.FullName != "" {
			consciousnessNumber, err := digits.GetConsciousnessNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating consciousness number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего сознания. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твой номер сознания: %d", consciousnessNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
		}
	} else if user.State == "RegisterFullName" {
		user.FullName = update.Message.Text
		user.State = ""
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user full name: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при обновлении полного имени. Попробуй позже.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Твоё полное имя успешно сохранено.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		msg = MakeStartMenu(*user, update)
		l.bot.Send(msg)

		if user.Birthdate != "" && user.Email != "" && user.FullName != "" {
			consciousnessNumber, err := digits.GetConsciousnessNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating consciousness number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего сознания. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твой номер сознания: %d", consciousnessNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
		}
	} else if user.State == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, я ещё не знаю, что тебе отвечать.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		return
	}
}

func MakeStartMenu(user storage.User, update tgbotapi.Update) tgbotapi.MessageConfig {
	buttons := map[string]tgbotapi.InlineKeyboardButton{
		"Birthdate": tgbotapi.NewInlineKeyboardButtonData("Ввести дату рождения", "State=RegisterBirthdate"),
		"Email":     tgbotapi.NewInlineKeyboardButtonData("Ввести Email", "State=RegisterEmail"),
		"FullName":  tgbotapi.NewInlineKeyboardButtonData("Ввести полное имя", "State=RegisterFullName"),
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

	Keyboard := tgbotapi.NewInlineKeyboardRow()

	if user.FullName != "" && user.Birthdate != "" && user.FullName != "" {
		Keyboard = append(Keyboard, buttons["Birthdate"])
		Keyboard = append(Keyboard, buttons["Email"])
		Keyboard = append(Keyboard, buttons["FullName"])
		RegisterMarkup = tgbotapi.NewInlineKeyboardMarkup(Keyboard)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет "+user.Name+", вот что мне сейчас изветно о тебе.\n"+CurrentData+"\nМожешь изменить свою анкету.\n")
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = RegisterMarkup
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	} else {
		if user.Birthdate == "" {
			Keyboard = append(Keyboard, buttons["Birthdate"])
		}
		if user.Email == "" {
			Keyboard = append(Keyboard, buttons["Email"])
		}
		if user.FullName == "" {
			Keyboard = append(Keyboard, buttons["FullName"])
		}
		RegisterMarkup = tgbotapi.NewInlineKeyboardMarkup(Keyboard)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет "+user.Name+", вот что мне сейчас изветно о тебе.\n"+CurrentData+"\nЗаполни свою анкету.\n")
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = RegisterMarkup
		msg.ReplyToMessageID = update.Message.MessageID
		return msg
	}

}
