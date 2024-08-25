package main

import (
	"Sprint2/internal/consumer"
	"Sprint2/internal/handler"
	"Sprint2/internal/rabbitmq"
	"Sprint2/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {

	//Инициализация БД
	db, err := storage.New()
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v", err)
	}

	// Подключение к RabbitMQ
	conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Не удалось установить соединение с RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Создание канала
	channel, err := rabbitmq.CreateChannel(conn)
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %v", err)
	}
	defer channel.Close()

	// Создание очереди для новых заказов
	_, err = rabbitmq.DeclareQueue(channel, rabbitmq.NewOrderQueue)
	
	// Создание очереди для уведомлений
	_, err = rabbitmq.DeclareQueue(channel, rabbitmq.NotificationQueue)
	

	// Создание очереди для заказов в обработке
	_, err = rabbitmq.DeclareQueue(channel, rabbitmq.ProcessingQueue)
	

	//Запуск обработчика и отправки уведомлений
	go consumer.ProcessOrders(channel)
	go consumer.SendNotification(channel)

	r := chi.NewRouter()
	r.Post("/api/create-order", handler.CreateOrder(channel,db))
	r.Get("/api/order/{id}", handler.GetOrder(db))

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер %v", err)
		return
	}

	log.Printf("Сервер запущен")

}
