package main

import (
	"digits_say/storage"
	"digits_say/telegram"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	config := Init()

	bot, err := telegram.NewListener(*config)

	if err != nil {
		panic(err)
	}

	bot.Start()
}

func Init() *telegram.Config {
	config := telegram.Config{}
	dbconfig := storage.DBConfig{}
	config.DBConfig = dbconfig

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file error: " + err.Error())
	}

	config.ApiToken = os.Getenv("TGAPIToken")

	config.DBConfig.ConnectionURL = os.Getenv("SurrealConnectionURL")
	config.DBConfig.Username = os.Getenv("SurrealUser")
	config.DBConfig.Password = os.Getenv("SurrealPassword")
	config.DBConfig.Namespace = os.Getenv("SurrealNamespace")
	config.DBConfig.Database = os.Getenv("SurrealDatabase")

	if env := os.Getenv("Debug"); "true" == strings.ToLower(env) {
		config.Debug = true
	} else {
		config.Debug = false
	}

	timeout, err := strconv.Atoi(os.Getenv("Timeout"))
	if err != nil {
		log.Fatalln("Can`t read enviroment variable 'Timeout': ", err)
	}
	config.Timeout = timeout

	return &config
}
