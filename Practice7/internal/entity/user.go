package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"ID" type:"uuid" default:"uuid_generate_v4()" primaryKey:"true"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	DeletedAt *time.Time `json:"DeletedAt" gorm:"index"`
	Username  string     `json:"Username"`
	Email     string     `json:"Email"`
	Password  string     `json:"Password"`
	Role      string     `json:"Role"`
	Verified  bool       `json:"Verified"`
}
