package telegram

import (
	"log"
	"net/http"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type bot struct {
	*tgbotapi.BotAPI
}

var Bot *bot

func init() {
	proxyURI := viper.GetString("telegram.proxy")
	transport := &http.Transport{}
	if proxyURI != "" {
		proxy, err := url.Parse(proxyURI)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 10,
	}
	botAPI, err := tgbotapi.NewBotAPIWithClient(viper.GetString("telegram.token"),
		client)
	if err != nil {
		log.Panicf("Error loading telegram: %v", err)
	}
	Bot = &bot{botAPI}
}

func (b *bot) Start() {
	go func() {
		updatesConfig := tgbotapi.NewUpdate(0)
		updatesConfig.Timeout = 60

		updatesChan := b.GetUpdatesChan(updatesConfig)
		for update := range updatesChan {
			switch {
			case update.Message != nil:
				go b.Message()
			case update.CallbackQuery != nil:
				b.CallbackQuery()
			}
		}
	}()
}
