package main

import (
	"PracticeBot/clients/telegramClients"
	"PracticeBot/consumer/event_consumer"
	"PracticeBot/events/telegramEvent"
	"PracticeBot/storage/sqlite"
	"context"
	"flag"
	"log"
)

const (
	tgBotHost   = "	api.telegram.org"
	storagePath = "data/sqlite/storage.db"
	batchSize   = 100
)

func main() {
	tgClient := telegramClients.NewClient(tgBotHost, getToken())
	st, err := sqlite.NewStorage(storagePath)
	if err != nil {
		log.Fatal("can't make storage", err)
	}

	ctx := context.Context(context.Background())

	if err := st.Init(ctx); err != nil {
		log.Fatal("failed to init storage", err)
	}

	eventManager := telegramEvent.NewEventManager(tgClient, st)

	consumer := event_consumer.NewConsumer(eventManager, eventManager, batchSize)

	log.Println("starting telegram bot")

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

func getToken() string {
	token := flag.String("token-bot", "", "bot token")
	flag.Parse()
	if *token == "" {
		log.Fatal("no token provided")
	}
	return *token
}
