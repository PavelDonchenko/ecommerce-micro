package service

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

const (
	USERROLE  = "user"
	ADMINROLE = "admin"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	PasswordHash []byte
	Claims       []Claims
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type Claims struct {
	ID    int
	Type  string
	Value string
}

type NewUser struct {
	ID        uuid.UUID
	Name      string
	Email     mail.Address
	Password  string
	Claims    []int
	CreatedAt time.Time
	UpdatedAt time.Time
}
