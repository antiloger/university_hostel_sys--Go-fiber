package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Avaliable string    `json:"avaliable"`
}

type Users struct {
	Users []User `json:"users"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
