package main

import (
	"flag"
	"github.com/vladislavsherwood/TelegramBot/clients/telegram"
	"log"
)

func main() {

	tgClient := telegram.New(mustBotHost(), mustToken())

	//token = flags.Get(token)

	//fetcher = fetcher.New(tgClient)

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
