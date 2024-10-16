package repositories

import (
	"net/mail"
	"time"

	"github.com/google/uuid"

	"github.com/PavelDonchenko/ecommerce-micro/auth/service"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	DeletedAt    time.Time `db:"deleted_at"`
}

type Claims struct {
	ID    int    `db:"id"`
	Type  string `db:"type"`
	Value string `db:"value"`
}

func toServiceUser(u User) service.User {
	addr := mail.Address{
		Address: u.Email,
	}

	return service.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        addr,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.In(time.Local),
		UpdatedAt:    u.UpdatedAt.In(time.Local),
		DeletedAt:    u.UpdatedAt.In(time.Local),
	}
}

func toServiceClaim(c Claims) service.Claims {
	return service.Claims{
		ID:    c.ID,
		Type:  c.Type,
		Value: c.Type,
	}
}

func toServiceClaims(claims []Claims) []service.Claims {
	result := make([]service.Claims, len(claims))
	for i := range claims {
		result[i] = toServiceClaim(claims[i])
	}

	return result
}

func toDBUser(u service.User) User {
	return User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email.Address,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
