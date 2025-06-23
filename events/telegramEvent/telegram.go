package telegramEvent

import (
	"PracticeBot/clients/telegramClients"
	"PracticeBot/events"
	"PracticeBot/storage"
	"fmt"
)

type EventManager struct {
	tgClient *telegramClients.Client
	storage  storage.Storage
	offset   int
}

type Meta struct {
	ChatID   int
	UserName string
}

func NewEventManager(tgClient *telegramClients.Client, st storage.Storage) *EventManager {
	return &EventManager{
		tgClient: tgClient,
		storage:  st,
	}
}

func (m *EventManager) Fetch(limit int) ([]events.Event, error) {
	updates, err := m.tgClient.GetUpdates(m.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		result = append(result, toEvent(update))
	}

	m.offset = updates[len(updates)-1].UpdateID + 1

	return result, nil
}

func (m *EventManager) Process(event events.Event) error {
	switch event.Type {
	case events.MessageEvent:
		err := m.processMessage(event)
		if err != nil {
			return fmt.Errorf("can't process message: %w", err)
		}
	default:
		return fmt.Errorf("unknown event type: %v", event.Type)
	}
	return nil
}

func (m *EventManager) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't get metadata: %w", err)
	}

	text := event.Text
	chatID := meta.ChatID
	username := meta.UserName

	if err := m.DoCmd(text, chatID, username); err != nil {
		return fmt.Errorf("can't do process page: %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta for event: %v", event.Type)
	}
	return result, nil
}

func toEvent(update telegramClients.Update) events.Event {
	eventType := fetchType(update)
	result := events.Event{
		Type: eventType,
		Text: fetchText(update),
	}

	if eventType == events.MessageEvent {
		result.Meta = Meta{
			ChatID:   update.Message.Chat.ChatID,
			UserName: update.Message.From.Username,
		}
	}

	return result
}

func fetchType(update telegramClients.Update) events.EventType {
	if update.Message == nil {
		return events.UnknownEvent
	}
	return events.MessageEvent
}

func fetchText(update telegramClients.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
