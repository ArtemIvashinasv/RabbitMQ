package storage

import (
	"Sprint2/internal/order"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
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

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Не удалось открыть бд: %v", err)
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка соединения с БД: %v", err)
		return err
	}

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

	// Выполняем SQL запросы
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Ошибка выполнения sql запроса: %v", err)
		return err
	}

	return nil
}

func (s *Storage) CreateOrder(o *order.Order) error {
	_, err := s.db.Exec("INSERT INTO orders (id,Name,count,status,created_at,price) values (:id,:Name,:count,:status,:created_at,:price)",
		sql.Named("id", o.ID),
		sql.Named("product", o.Name),
		sql.Named("count", o.Count),
		sql.Named("status", o.Status),
		sql.Named("created_at", o.CreatedAt),
		sql.Named("price", o.Price))

	if err != nil {
		log.Fatalf("Не добавили заказ: %v", err)
		return err
	}
	log.Println("Успешное добавление заказа в бд с id:", o.ID)

	return nil
}
