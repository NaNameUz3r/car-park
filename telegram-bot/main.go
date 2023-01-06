package main

import (
	"car-park/telegram-bot/clients/telegram"
	"os"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())
	// fetcher := fetcher.New()
	// processor := processor.New()
	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := os.Getenv("TELEGRAM_APIKEY")
	if token == "" {
		panic("No token in ENV. Set it with TELEGRAM_APIKEY key")
	} else {
		return token
	}
}
