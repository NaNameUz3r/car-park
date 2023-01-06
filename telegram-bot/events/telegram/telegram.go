package telegram

import (
	"car-park/telegram-bot/clients/telegram"
	"car-park/telegram-bot/events"
	"car-park/telegram-bot/lib/er"
	"car-park/telegram-bot/storage"
	"errors"
)

var ErrUnknownEvent = errors.New("unknown event type")
vat ErrUnknownMetaType = errors.New("unknown meta type")

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, er.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		p.processMessage(event)
	default:
		return er.Wrap("can't process message", ErrUnknownEvent)
	}
}

func event(u telegram.Update) events.Event {
	updateType := getUpdateType(u)
	res := events.Event{
		Type: updateType,
		Text: getUpdateText(u),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res

}

func getUpdateType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func getUpdateText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}

func (p *Processor) processMessage(event events.Event) {
	meta, err := getMeta(event)
	if err != nil {
		return er.Wrap("can't process message", err)
	}
}

func getMeta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{},er.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}