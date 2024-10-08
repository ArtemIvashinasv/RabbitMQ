package order

import (
	"time"
)

type Notification struct {
	OrderId string `json:"id"`
	Status  string `json:"new_status"`
	Message string `json:"message"`
}

type Order struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Count     int       `json:"count"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Price     float64   `json:"price"`
}
