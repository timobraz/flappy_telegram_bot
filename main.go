package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Play struct {
	Id    int
	Time  string
	Score int
	Hash  string
}

func main() {

	db, _ := sql.Open("sqlite3", "database.db")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS plays (
	id INTEGER NOT NULL PRIMARY KEY,
	time DATETIME NOT NULL,
	score INTEGER NOT NULL,
	hash TEXT NOT NULL
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/upload", func(c *gin.Context) {
		hash := c.DefaultQuery("hash", "guest")
		score := c.DefaultQuery("score", "0")
		log.Printf("Firstname: %s; Score: %s", hash, score)
		_, err := db.Exec(`INSERT INTO plays (time, score, hash) VALUES (datetime('now'), ?, ?);`, score, hash)
		if err != nil {
			log.Fatal(err)
		}

		c.String(http.StatusOK, "INSERTED %s %s", hash, score)
	})
	r.GET("/", func(c *gin.Context) {
		res, err := db.Query(`SELECT * FROM plays;`)

		if err != nil {
			log.Fatal(err)
		}
		defer res.Close()
		var plays []Play
		for res.Next() {
			var id int
			var time string
			var score int
			var hash string
			err = res.Scan(&id, &time, &score, &hash)
			if err != nil {
				log.Fatal(err)
			}
			plays = append(plays, Play{id, time, score, hash})
		}
		_, err = json.Marshal(plays)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, plays)
	})
	r.Run()

	bot, err := tgbotapi.NewBotAPI("7181024480:AAF1_1hvuUr3LfPUis4omoVZZXO80GCjTdk")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if strings.Contains("catty_flappy_bot", update.Message.Text) {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Play the game here: https://timobraz.github.io/cat-flappy/")
				bot.Send(msg)
			}
		}
	}
}
