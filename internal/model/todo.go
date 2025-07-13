package model

import (
	"time"
)

type Todo struct {
	ID        int       `json:"id" example:"1"`
	UserID    int       `json:"userId" example:"1"`
	Title     string    `json:"title" example:"114"`
	Content   string    `json:"content" example:"514"`
	CreatedAt time.Time `json:"createdAt,omitempty" example:"2025-07-07 21:51:47"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" example:"2025-07-07 21:51:47"`
	Done      bool      `json:"done,omitempty" example:"false"`
}
