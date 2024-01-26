package main

import (
	"flag"
	"log"

	tgClient "github.com/vladislavsherwood/TelegramBot/clients/telegram"
	"github.com/vladislavsherwood/TelegramBot/consumer/event_consumer"
	"github.com/vladislavsherwood/TelegramBot/events/telegram"
	"github.com/vladislavsherwood/TelegramBot/lib/storage/files"
)

const (
	tgBotHost = "api.telegram.org"

	storagePath = "files_storage"

	batchSize = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}

func mustToken() string {
	token := flag.String(
		"t",
		"",
		"token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}

func mustBotHost() string {
	host := flag.String(
		"host",
		"",
		"host for access to telegram bot")
	flag.Parse()

	if *host == "" {
		log.Fatal("host is not specified")
	}
	return *host
}
