package telegramEvent

import (
	"PracticeBot/clients/telegramClients"
	"PracticeBot/storage"
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (m *EventManager) DoCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	if isAddCmd(text) {
		if err := m.savePage(chatID, text, username); err != nil {
			return fmt.Errorf("can't save page: %w", err)
		}
		return nil
	}

	switch text {
	case RndCmd:
		if err := m.sendRandomPage(chatID, username); err != nil {
			return fmt.Errorf("can't send random page: %w", err)
		}
		return nil
	case HelpCmd:
		if err := m.sendHelp(chatID); err != nil {
			return fmt.Errorf("can't send help: %w", err)
		}
		return nil
	case StartCmd:
		if err := m.sendHello(chatID); err != nil {
			return fmt.Errorf("can't send hello: %w", err)
		}
		return nil
	default:
		if err := m.tgClient.SendMessage(chatID, msgUnknownCmd); err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
		return nil
	}
}

func (m *EventManager) savePage(chatID int, pageURL, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	sendMsg := newMessageSender(chatID, m.tgClient)

	ctx := context.Context(context.Background())

	isExists, err := m.storage.IsExists(ctx, page)
	if err != nil {
		return fmt.Errorf("can't understand exist page or not: %w", err)
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := m.storage.Save(ctx, page); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	if err := sendMsg(msgSaved); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (m *EventManager) sendRandomPage(chatID int, username string) error {
	ctx := context.Context(context.Background())
	page, err := m.storage.PickRandom(ctx, username)
	if err != nil {
		return fmt.Errorf("can't pick random page: %w", err)
	}
	sendMsg := newMessageSender(chatID, m.tgClient)
	if page == nil {
		if err = sendMsg(msgNoSavedPages); err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
		return nil
	}

	if err := sendMsg(page.URL); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	if err := m.storage.Remove(ctx, page); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (m *EventManager) sendHelp(chatID int) error {
	sendMsg := newMessageSender(chatID, m.tgClient)
	if err := sendMsg(msgHelp); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (m *EventManager) sendHello(chatID int) error {
	sendMsg := newMessageSender(chatID, m.tgClient)
	if err := sendMsg(msgHello); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func newMessageSender(chatID int, tgClient *telegramClients.Client) func(msg string) error {
	return func(msg string) error {
		return tgClient.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Scheme != "" && u.Host != ""
}
