package telegram

import (
	"errors"
	"github.com/vladislavsherwood/TelegramBot/lib/e"
	"github.com/vladislavsherwood/TelegramBot/lib/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// doCmd представляет из себя подобие api-роутера
// т.е. по формату и содержанию текста и сообщения будет пониматься роль команды
func (p *Processor) doCmd(text string, chatID int, username string) error {
	// Удаляем пробелы
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)
	// Проверка если текст является ссылкой
	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:

		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:

		return p.sendHello(chatID)
	default:

		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: Save Page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	IsExist, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	// Если страница уже есть, пользователю возвращается сообщение msgAlreadyExists
	if IsExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	// Если страницу сохранить не удалось, возвращается ошибка
	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: Can't Send Random", err) }()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

// Функция отправляет msgHelp
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

// Функция отправляет msgHello
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// Функция возвращает, добавлять ли значение в датабазу
func isAddCmd(text string) bool {
	return isURL(text)
}

// Функция возвращает, является ли text ссылкой, ссылки без протокола не будут восприниматься как ссылки
func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
