package model

import (
	"time"
)

type Contact struct {
	ID             string     `json:"id" db:"id" sql:"id"`
	Email          *string    `json:"email" db:"email" sql:"email"`
	PhoneNumber    *string    `json:"phone_number" db:"phone_number" sql:"phone_number"`
	LinkedID       *string    `json:"linked_id" db:"linked_id" sql:"linked_id"`
	LinkPrecedence string     `json:"link_precedence" db:"link_precedence" sql:"link_precedence"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at" sql:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at" sql:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at" sql:"deleted_at"`
}
