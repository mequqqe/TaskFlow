package delivery

import (
	"fmt"
	"log"
	"taskflow/internal/domain"
	"taskflow/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type OrderService struct {
	repo     repository.OrderRepository
	botToken string
	chatID   int64
}

func NewOrderService(repo repository.OrderRepository, botToken string, chatID int64) *OrderService {
	return &OrderService{
		repo:     repo,
		botToken: botToken,
		chatID:   chatID,
	}
}

func SendTelegramNotification(botToken string, chatID int64, message string, fileData []byte, fileName, fileType string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return fmt.Errorf("failed to create bot: %w", err)
	}

	// Отправка текстового сообщения
	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Отправка файла, если он присутствует
	if len(fileData) > 0 {
		file := tgbotapi.FileBytes{
			Name:  fileName,
			Bytes: fileData,
		}

		var document tgbotapi.Chattable
		if fileType == "image/jpeg" || fileType == "image/png" {
			document = tgbotapi.NewPhotoUpload(chatID, file)
		} else {
			document = tgbotapi.NewDocumentUpload(chatID, file)
		}

		_, err := bot.Send(document)
		if err != nil {
			log.Printf("Failed to send file: %v", err)
			return fmt.Errorf("failed to send file: %w", err)
		}
	}

	log.Printf("Message and files sent successfully to chat ID %d", chatID)
	return nil
}


func (s *OrderService) CreateOrder(order *domain.Order) error {
	if err := s.repo.Create(order); err != nil {
		log.Printf("Ошибка при создании заказа: %v", err)
		return err
	}

	log.Println("Заказ успешно создан")

	// Формируем сообщение и отправляем его вместе с файлами
	message := fmt.Sprintf("Новый заказ!\nИмя предпринимателя: %s\nТема: %s\nСумма: %.2f\nТребования: %s\nДедлайн: %s",
		order.EntrepreneurName, order.Theme, order.Amount, order.Requirements, order.Deadline.Format("02-01-2006"))

	if err := SendTelegramNotification(s.botToken, s.chatID, message, order.FileData, order.FileName, order.FileType); err != nil {
		log.Printf("Ошибка при отправке уведомления: %v", err)
		return fmt.Errorf("failed to send telegram notification: %w", err)
	}

	log.Println("Уведомление и файлы отправлены успешно")
	return nil
}
