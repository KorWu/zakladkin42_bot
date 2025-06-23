package event_consumer

import (
	"PracticeBot/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("ERR consumer fetching: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("ERR consumer handling: %s", err.Error())

			continue
		}
	}
}

/*
TODO:
1. Потеря событий: ретраи, возвращение в хранилище, фолбэ, подтверждение
2. обработка всей пачки: останавливаться после первой ошибки, вести счетчик
3. параллельная обработка (sync.WaitGroup)
*/

func (c Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("new event: %v", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("ERR processor: %s", err.Error())

			continue
		}
	}

	return nil
}
