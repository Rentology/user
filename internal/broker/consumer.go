package broker

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"user-service/internal/models"
	userHttp "user-service/internal/user/delivery/http"
	"user-service/lib/sl"
)

type UserConsumer struct {
	service userHttp.UserService
	channel *amqp.Channel
	queue   string
	log     *slog.Logger
}

func NewUserConsumer(service userHttp.UserService, b *Broker, queueName string, log *slog.Logger) (*UserConsumer, error) {
	ch, err := b.CreateChannel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &UserConsumer{service, ch, queueName, log}, nil
}

// Run запускает Consumer для обработки сообщений
func (c *UserConsumer) Run(ctx context.Context) error {

	// Подписываемся на очередь
	messages, err := c.channel.Consume(
		c.queue,         // Имя очереди
		"user_consumer", // Имя консюмера
		false,           // Авто-подтверждение
		false,           // Эксклюзивное подключение
		false,           // Локальная очередь
		false,           // Без ожидания
		nil,             // Дополнительные аргументы
	)
	if err != nil {
		return err
	}

	// Обработка сообщений
	go func() {
		for msg := range messages {
			if err := c.processMessage(ctx, msg); err != nil {
				sl.Infof(c.log, "Ошибка обработки сообщения: %v", err)

				// Сообщение не обработано, отправляем в очередь повторно
				if err := msg.Nack(false, true); err != nil {
					sl.Infof(c.log, "Ошибка повторной отправки сообщения: %v", err)
				}
			}
		}
	}()

	return nil
}

// processMessage обрабатывает одно сообщение
func (c *UserConsumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	// Декодируем сообщение в структуру User
	var user models.User
	if err := json.Unmarshal(msg.Body, &user); err != nil {
		sl.Infof(c.log, "Ошибка декодирования сообщения: %v", err)
		return err
	}

	// Вызываем метод Create из userService
	createdUser, err := c.service.Create(ctx, &user)
	if err != nil {
		sl.Infof(c.log, "Ошибка создания пользователя: %v", err)
		return err
	}

	sl.Infof(c.log, "Пользователь успешно создан: %+v", createdUser)

	// Подтверждаем сообщение
	if err := msg.Ack(false); err != nil {
		sl.Infof(c.log, "Ошибка подтверждения сообщения: %v", err)
		return err
	}

	return nil
}
