package models

import "time"

type OrderEvent struct {
    OrderID     string    `json:"order_id"`
    UserID      string    `json:"user_id"`
    UserName    string    `json:"user_name"`
    UserEmail   string    `json:"user_email"`
    Items       []Item    `json:"items"`
    TotalAmount float64   `json:"total_amount"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

type Item struct {
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

type NotificationEvent struct {
    Type      string    `json:"type"`
    UserID    string    `json:"user_id"`
    UserEmail string    `json:"user_email"`
    Message   string    `json:"message"`
    OrderID   string    `json:"order_id"`
    CreatedAt time.Time `json:"created_at"`
}