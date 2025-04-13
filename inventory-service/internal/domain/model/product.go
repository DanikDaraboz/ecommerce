package model

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}