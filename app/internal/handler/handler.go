package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Sprint2/internal/order"
	"Sprint2/internal/storage"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func CreateOrder(channel *amqp091.Channel, storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var o order.Order

		err := json.NewDecoder(r.Body).Decode(&o)
		if err != nil {
			log.Println("Не удалось закодировать заказ")
			http.Error(w, "Не валидный запрос", http.StatusBadRequest)
			return
		}

		o.ID = uuid.New().String()
		o.CreatedAt = time.Now()
		o.Status = "new"

		err = order.Valid(&o)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Добавляем заказа в БД
		err = storage.CreateOrder(&o)
		if err != nil {
			log.Fatalf("Ошибка добавления заказа в БД: %v", err)
			return
		}

		// Сообщение для RabbitMQ
		msgPublish := map[string]string{
			"ID":     o.ID,
			"Status": o.Status,
		}

		msgRabit, err := json.Marshal(msgPublish)
		if err != nil {
			http.Error(w, "Ошибка сериализации в json", http.StatusInternalServerError)
			return
		}

		//Отправляем заказ  в очередь
		err = channel.Publish(
			"",
			"new_orders_queue",
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        msgRabit,
			})
		if err != nil {
			log.Println("Не удалось опубликовать новый заказ", err)
			http.Error(w, "Не удалось опубликовать заказ в очередь", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(o)
		log.Println("Заказ успешно создан")
	}
}
