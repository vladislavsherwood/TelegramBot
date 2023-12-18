package telegram

import (
	"errors"
	"github.com/vladislavsherwood/TelegramBot/clients/telegram"
	"github.com/vladislavsherwood/TelegramBot/events"
	"github.com/vladislavsherwood/TelegramBot/lib/e"
	"github.com/vladislavsherwood/TelegramBot/lib/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	chatID   int
	Username string
}

var (
	ErrUnknownErrType  = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	// Получаем апдейты (понятие в телеграмме)
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	// Если список апдейтов оказался пустым, то говорим, что ничего не нашли
	if len(updates) == 0 {
		return nil, nil
	}
	// Подгтовливается переменная для результата, заранее аллоцируя память для нее
	res := make([]events.Event, 0, len(updates))

	//Перебираем все полученные апдейты и преобразуем их в тип event - более общую сущность, которую используем далее
	for _, u := range updates {
		res = append(res, event(u))
	}
	//Обновляем параметр offset, чтобы в следующий раз получить следующую пачку изменений
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownErrType)

	}
}

func (p *Processor) processMessage(event events.Event) error {
	// Получаем мету при помощи функции meta
	meta, err := meta(event)
	// Обрабатываем ошибку
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			chatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
