package telegram

import "github.com/vladislavsherwood/TelegramBot/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
}

func New(client *telegram.Client)
