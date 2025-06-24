package event_consumer

import (
	"PracticeBot/events"
	"log"
	"time"
	"sync"
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
Потеря событий: ретраи, возвращение в хранилище, фолбэ, подтверждение
*/

func (c Consumer) handleEvents(eventsList []events.Event) error {

	wg := sync.WaitGroup{}

	var cntErrors countErrors

	for _, event := range eventsList {
		log.Printf("new event: %v", event.Text)

		cntErrors.mu.Lock()
		if cntErrors.cnt >= 5 {
			log.Printf("%d events were not processed in a row. Programm stopped", cntErrors.cnt)
			break
		}
		cntErrors.mu.Unlock()

		wg.Add(1)
		go func(ev events.Event) {
			defer wg.Done()
			err := c.processor.Process(ev)
			if err != nil {
				log.Printf("ERR processing: %s", err.Error())
				cntErrors.mu.Lock()
				cntErrors.cnt++
				cntErrors.mu.Unlock()
			}
		}(event)
	}
	wg.Wait()
	return nil
}

type countErrors struct {
	cnt int
	mu  sync.Mutex
}