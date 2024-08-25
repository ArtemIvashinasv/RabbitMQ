# RabbitMQ

## Описание

Сервис по созданию заказов, сохранению их в БД , так же отправляет не большой JSON обьект для с статусами обработки заказов

API имеет следующие операции:
- Добавить заказ
- Получить заказ по id

# Запуск
Для запуск потребуется СУБД Postgresql и RabbitMQ
## 1.Команда для установки образа/запуска POSTGRES
- docker run --name my-postgres -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=mydb -p 5432:5432 -v postgres-data:/var/lib/postgresql/data -d postgres:15

## 2.Для запуска приложения необходимо находится в корневой директории app  (команда для перехода `cd app`)
- для запуска приложения в консоли необходимо запустить команду go run `cmd/main.go`

## 3.Проверка функционала добавления заказа/обработка/смена статуса
- Открыть второй окно консоли и сделать запрос 
    curl -X POST http://localhost:8080/create-order \
    -H "Content-Type: application/json" \
    -d '{
        "name": "John Doe",
        "count": 123,
        "price": 2
    }'


