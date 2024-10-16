package http

import (
	"net/mail"

	"github.com/PavelDonchenko/ecommerce-micro/auth/service"
)

type User struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Claims    []Claims `json:"claims,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	DeletedAt string   `json:"deleted_at,omitempty"`
}

type Claims struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func toAppClaims(claims service.Claims) Claims {
	return Claims{
		ID:    claims.ID,
		Type:  claims.Type,
		Value: claims.Value,
	}
}

type NewUser struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"eqfield=Password"`
	Claims          []int  `json:"claims"`
}

func toServiceNewUser(nu NewUser) (service.NewUser, error) {
	addr, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return service.NewUser{}, err
	}
	return service.NewUser{
		Name:     nu.Name,
		Email:    *addr,
		Password: nu.Password,
		Claims:   nu.Claims,
	}, err
}

func toAppUser(u service.User) User {
	var claims []Claims
	if u.Claims != nil {
		for _, claim := range u.Claims {
			claims = append(claims, toAppClaims(claim))
		}
	}

	return User{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email.Address,
		Claims:    claims,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		DeletedAt: u.DeletedAt.String(),
	}
}
