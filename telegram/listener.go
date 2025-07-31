package telegram

import (
	"digits_say/digits"
	"digits_say/storage"
	"digits_say/telegram/menus"
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
				State:        map[string]interface{}{},
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

			msg := menus.MakeStartMenu(*user, update)
			l.bot.Send(msg)
			return
		}
		if exist {
			msg := menus.MakeStartMenu(*user, update)
			l.bot.Send(msg)
			return
		}

	} else if update.Message.Text == "/menu" {
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
		msg := menus.MakeMainMenu(*user, update.Message.MessageID, update.Message.Chat.ID, "Главное меню.")
		_, err = l.bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending main menu: ", err)
		}
	} else if update.Message.Text == "/Начать" {

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
		user.State["Register"] = "Birthdate"
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
		user.State["Register"] = "Email"
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
		user.State["Register"] = "FullName"
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
	} else if callback.Data == "State=RegisterFinished" {
		user, _, err := l.storage.GetUserByID(fmt.Sprintf("tg%d", callback.From.ID))
		if err != nil {
			log.Println("Error getting user by Telegram ID: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при получении пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.State["Register"] = "Finished"
		_, err = l.storage.UpdateUser(*user)
		if err != nil {
			log.Println("Error updating user state: ", err)
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.")
			msg.ReplyToMessageID = callback.Message.MessageID
			l.bot.Send(msg)
			return
		}
		msg := menus.MakeMainMenu(*user, callback.Message.MessageID, callback.Message.Chat.ID, "Регистрация завершена. Вот главное меню.")
		_, err = l.bot.Send(msg)
		fmt.Println(err)
		return
	} else {
		log.Println("Unknown callback data: ", callback.Data)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Извини, я не понимаю эту команду.")
		msg.ReplyToMessageID = callback.Message.MessageID
		l.bot.Send(msg)
		return
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

	if user.State["Register"] == "Birthdate" {
		birthdate, err := time.Parse("02.01.2006", update.Message.Text)
		if err != nil {
			log.Println("Error parsing birthdate: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты. Пожалуйста, введи дату в формате ДД.ММ.ГГГГ.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.Birthdate = birthdate.Format("2006-01-02") //models.CustomDateTime{Time: birthdate}
		user.State["Register"] = ""
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
		msg = menus.MakeStartMenu(*user, update)
		l.bot.Send(msg)

		/*if user.Birthdate != "" && user.Email != "" && user.FullName != "" {
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
		}*/
		return
	} else if user.State["Register"] == "Email" {
		_, err := mail.ParseAddress(update.Message.Text)
		if err != nil {
			log.Println("Error parsing email: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат Email. Пожалуйста, введи корректный Email.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
		user.Email = update.Message.Text
		user.State["Register"] = ""
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
		msg = menus.MakeStartMenu(*user, update)
		l.bot.Send(msg)

		/*if user.Birthdate != "" && user.Email != "" && user.FullName != "" {
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
		}*/
		return
	} else if user.State["Register"] == "FullName" {
		user.FullName = update.Message.Text
		user.State["Register"] = ""
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
		msg = menus.MakeStartMenu(*user, update)
		l.bot.Send(msg)

		/*if user.Birthdate != "" && user.Email != "" && user.FullName != "" {
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

		}*/
		return
	} else if update.Message.Text == "Число сознания" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			consciousnessNumber, err := digits.GetConsciousnessNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating consciousness number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего сознания. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твоё число сознания: %d", consciousnessNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}

	} else if update.Message.Text == "Число действия" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			actionNumber, err := digits.GetActionNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating action number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего действия. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твоё число дейсвия: %d", actionNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
	} else if update.Message.Text == "Число кармы" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			karmaNumber, err := digits.GetKarmaNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating karma number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего числа кармы. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твоё число кармы: %d", karmaNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
	} else if update.Message.Text == "Число месяца" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			monthNumber, err := digits.GetMonthNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating month number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего месяца. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твоё число месяца: %d", monthNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
	} else if update.Message.Text == "Число года" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			monthNumber, err := digits.GetYearNumber(user.Birthdate)
			if err != nil {
				log.Println("Error calculating month number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего года. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Твоё число года: %d", monthNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
	}else if update.Message.Text == "Личный день" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			monthNumber, err := digits.GetPrivateDay()
			if err != nil {
				log.Println("Error calculating private day number: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте твоего личного дня. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Число личного дня: %d", monthNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
	}else if update.Message.Text == "Общий день" {
		if user.Birthdate != "" && user.State["Register"] == "Finished" {
			monthNumber, err := digits.GetPublicDay(user.Birthdate)
			if err != nil {
				log.Println("Error calculating shared day: ", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при расчёте общего дня. Попробуй позже.")
				msg.ReplyToMessageID = update.Message.MessageID
				l.bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Число общего дня: %d", monthNumber))
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже, что ты ещё не завершил регистрацию. Напиши /start, чтобы начать.")
			msg.ReplyToMessageID = update.Message.MessageID
			l.bot.Send(msg)
			return
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, я ещё не знаю, что тебе отвечать.")
		msg.ReplyToMessageID = update.Message.MessageID
		l.bot.Send(msg)
		return
	}
}
