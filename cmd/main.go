package main

import (
	"log"
	"net/http"
	db "taskflow/data"
	"taskflow/internal/controller"
	"taskflow/internal/delivery"
	"taskflow/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gorilla/mux"
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	botToken := "7081061730:AAGhXTl_9trEkrY95m9JaGFZcskEUqHOzrQ"
	chatID := int64(2099795903)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Не удалось создать бот: %v", err)
	}

	message := "Проверка работы Telegram бота"
	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Не удалось отправить сообщение: %v", err)
	}

	log.Println("Сообщение успешно отправлено")

	orderRepo := repository.NewOrderRepository(db)
	orderService := delivery.NewOrderService(orderRepo, botToken, chatID)
	orderController := controller.NewOrderController(orderService)

	r := mux.NewRouter()
	r.HandleFunc("/orders", orderController.CreateOrder).Methods("POST")

	log.Println("Registering routes...")
	http.Handle("/", r)
	log.Println("Server starting on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
