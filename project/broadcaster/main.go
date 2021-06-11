package main

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Todo struct {
	Id        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func mustLookupEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Errorf("environment value not set: %v", key))
	}
	return value
}

func main() {
	telegramAPIToken := mustLookupEnv("TELEGRAM_API_TOKEN")
	telegramChannelID := mustLookupEnv("TELEGRAM_CHANNEL_ID")
	natsURL := mustLookupEnv("NATS_URL")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	rx := make(chan *Todo)
	ec.BindRecvQueueChan("todos", "broadcaster", rx)

	bot, err := tgbotapi.NewBotAPI(telegramAPIToken)
	if err != nil {
		panic(err)
	}

	for todo := range rx {
		fmt.Printf("received: %+xv", &todo)
		msg := tgbotapi.NewMessageToChannel(telegramChannelID, fmt.Sprintf("New todo with text: %v", todo.Text))
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
