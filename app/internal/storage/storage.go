package storage

import (
	"Sprint2/internal/order"
	_ "database/sql"
	"time"

	//"fmt"
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *pgxpool.Pool
}

func New() (*Storage, error) {
	s := &Storage{}

	err := s.initDB()
	if err != nil {
		return nil, err
	}
	return s, err
}

func (s *Storage) initDB() error {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки end: %v", err)
	}

	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// pass := os.Getenv("DB_PASSWORD")
	// dbname := os.Getenv("DB_NAME")
	dsn := "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"

	//dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	dbPool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Не удалось открыть бд: %v", err)
		return err
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ошибка соединения с БД: %v", err)
		return err
	}

	s.db = dbPool

	sqlFile, err := os.Open("internal/SQL/createTable.sql")
	if err != nil {
		log.Fatalf("Не удалось открыть SQL файл: %v", err)
		return err
	}
	defer sqlFile.Close()

	sqlBytes, err := ioutil.ReadAll(sqlFile)
	if err != nil {
		log.Fatalf("Не удалось прочитать SQL файл: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Выполняем SQL запросы
	_, err = s.db.Exec(ctx, string(sqlBytes))
	if err != nil {
		log.Fatalf("Ошибка выполнения sql запроса: %v", err)
		return err
	}

	return nil
}

func (s *Storage) CreateOrder(o *order.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.db.Exec(ctx, "INSERT INTO orders (id,Name,count,status,created_at,price) values ($1,$2,$3,$4,$5,$6)",
		o.ID,
		o.Name,
		o.Count,
		o.Status,
		o.CreatedAt,
		o.Price)

	if err != nil {
		log.Fatalf("Не добавили заказ: %v", err)
		return err
	}
	log.Println("Успешное добавление заказа в бд с id:", o.ID)

	return nil
}

func (s *Storage) GetOrderByID(id string) (*order.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// SQL-запрос с использованием позиционного параметра
	query := `
		SELECT id, name, count, status, created_at, price
		FROM orders
		WHERE id = $1
	`

	// Выполняем запрос
	row := s.db.QueryRow(ctx, query, id)

	// Создаем переменную для хранения результата
	var o order.Order

	// Считываем результат в переменную
	err := row.Scan(
		&o.ID,
		&o.Name,
		&o.Count,
		&o.Status,
		&o.CreatedAt,
		&o.Price,
	)

	if err != nil {
		log.Printf("Ошибка получения заказа: %v", err)
		return nil, err
	}

	return &o, nil
}
