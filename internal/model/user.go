package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id" example:"1"`
	Email     string    `json:"email" example:"ken@ken.me"`
	FirstName string    `json:"firstName,omitempty" example:"Ken"`
	LastName  string    `json:"lastName,omitempty" example:"Chen"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt" example:"2025-07-07 21:51:47"`
	UpdatedAt time.Time `json:"updatedAt" example:"2025-07-07 21:51:47"`
}
