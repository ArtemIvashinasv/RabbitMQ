package order

import "errors"

func Valid(order *Order) error {
	if order.Name == "" {
		return errors.New("Имя не может быть пустым")
	}

	if order.Count <= 0 {
		return errors.New("Количество товаров не может быть <= 0")
	}

	if order.Status == "" {
		return errors.New("Статус не может быть пустым")
	}

	return nil
}
