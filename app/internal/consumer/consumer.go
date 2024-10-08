package consumer

import (
	"Sprint2/internal/order"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessOrders(channel *amqp.Channel) {
	msgs, err := channel.Consume(
		"new_orders_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистировать consumer %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			//var incomingMsg order.Notification
			var notification order.Notification

			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				log.Println("Не удалось декодировать заказ ", err)
				continue
			}
			log.Println("Заказ получен")

			notification.Status = "process"
			notification.Message = "Заказа обработан"
			log.Println("Заказ обработан")

			notificationJSON, err := json.Marshal(&notification)
			if err != nil {
				log.Println("Не удалось закодировать уведомление ", err)
				continue
			}

			err = channel.Publish(
				"",
				"notification_queue",
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        notificationJSON,
				})
			if err != nil {
				log.Println("Не удалось отправить уведомление", err)
			}
			log.Println("Уведомление успешно отправленно")

		}
	}()
	log.Println("Чтение для выхода Ctrl+C")
	<-forever
}

func SendNotification(channel *amqp.Channel) {
	msgs, err := channel.Consume(
		"notification_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var notification order.Notification
			err = json.Unmarshal(msg.Body, &notification)
			if err != nil {
				log.Printf("Ошибка декодирования JSON: %v", err)
				continue
			}

			log.Printf("Полученно уведомление заказа c id: %v, Новым статусом: %v", notification.OrderId, notification.Status)
		}
	}()

	log.Println("Ожидание уведомлений")
	<-forever
}
