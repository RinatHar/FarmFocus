package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Config struct {
	Token string
}

func loadConfig() *Config {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	return &Config{
		Token: token,
	}
}

func main() {
	// Настройка логгера для Docker
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Загрузка конфигурации
	config := loadConfig()
	log.Println("Config loaded successfully")

	// Инициализация API клиента с обработкой ошибки
	api, err := maxbot.New(config.Token)
	if err != nil {
		log.Fatalf("Failed to create bot API client: %v", err)
	}
	log.Println("Bot API client created successfully")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Получение информации о боте
	info, err := api.Bots.GetBot(ctx)
	if err != nil {
		log.Printf("Error getting bot info: %v", err)
	} else {
		log.Printf("Bot info: Name=%s, ID=%s", info.Name, info.UserId)
	}

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		sig := <-exit
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	log.Println("Starting to listen for updates...")

	updateCount := 0
	for upd := range api.GetUpdates(ctx) {
		updateCount++
		log.Printf("Received update #%d: %T", updateCount, upd)

		switch upd := upd.(type) {
		case *schemes.MessageCreatedUpdate:
			log.Printf("Message from user %s: %s", upd.Message.Sender.Username, upd.Message.Body.Text)

			// Отправка сообщения с именем пользователя
			userName := upd.Message.Sender.Username
			if userName == "" {
				userName = "друг" // если имя пользователя пустое
			}

			messageText := fmt.Sprintf("Привет, %s! Ты написал: \"%s\"", userName, upd.Message.Body.Text)
			_, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(messageText))
			if err != nil {
				log.Printf("Error sending message: %v", err)
			} else {
				log.Printf("Message sent successfully to user %s", userName)
			}

		default:
			log.Printf("Unhandled update type: %T", upd)
		}
	}

	log.Println("Bot stopped")
}
