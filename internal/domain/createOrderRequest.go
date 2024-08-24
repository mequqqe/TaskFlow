package domain

import (
	"time"
)

// Order представляет заказ
type Order struct {
	ID               uint                 `gorm:"primaryKey"`
	EntrepreneurName string               `json:"entrepreneur_name"`
	Theme            string               `json:"theme"`
	Amount           float64              `json:"amount"`
	Deadline         time.Time            `json:"deadline"`
	Requirements     string               `json:"requirements"`
	Status           string               `json:"status"`
	StatusHistories  []OrderStatusHistory `json:"status_histories"`
}

type OrderStatusHistory struct {
	ID        uint      `gorm:"primaryKey"`
	OrderID   uint      `gorm:"index"`
	Status    string    `json:"status"`
	ChangedAt time.Time `json:"changed_at"`
}
