package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"time"
)

type Config struct {
	TelegramBotToken string
	DBName           string
	DBUser           string
	DBPass           string
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	db, err := sql.Open("postgres", "user="+configuration.DBUser+" dbname="+configuration.DBName+" password="+configuration.DBPass+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for {
		select {
		case update := <-updates:
			UserName := update.Message.From.UserName
			ChatID := update.Message.Chat.ID
			Text := update.Message.Text
			log.Printf("[%s] %d %s", UserName, ChatID, Text)
			_, err := db.Exec("INSERT INTO messages VALUES (($1), ($2), ($3));", UserName, time.Now(), Text)

			if err != nil {
				log.Printf("Error during inserting ", err)
			}
		}
	}
}
