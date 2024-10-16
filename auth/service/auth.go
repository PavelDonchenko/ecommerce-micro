package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/PavelDonchenko/ecommerce-micro/common/database/sqldb"
	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
	"github.com/PavelDonchenko/ecommerce-micro/common/mid"
)

var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

type UserRepo interface {
	NewWithTx(tx sqldb.CommitRollbacker, log *logger.Logger) (UserRepo, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) error
	GetClaimsByType(ctx context.Context, claimType string) ([]Claims, error)
	CreateUserClaims(ctx context.Context, claims []int, uid uuid.UUID) error
}

type Auth struct {
	userRepo UserRepo
	log      *logger.Logger
}

func NewAuth(userRepo UserRepo, log *logger.Logger) *Auth {
	return &Auth{userRepo: userRepo, log: log}
}

func (a *Auth) newWithTX(ctx context.Context) (*Auth, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	ur, err := a.userRepo.NewWithTx(tx, a.log)
	if err != nil {
		return nil, err
	}

	return &Auth{
		userRepo: ur,
		log:      a.log,
	}, nil
}

func (a *Auth) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return a.userRepo.GetByEmail(ctx, email)
}

func (a *Auth) CreateUser(ctx context.Context, nu NewUser) (User, error) {
	a, err := a.newWithTX(ctx)

	now := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	if len(nu.Claims) == 0 {
		usersClaims, err := a.userRepo.GetClaimsByType(ctx, USERROLE)
		if err != nil {
			return User{}, fmt.Errorf("get users claims: %w", err)
		}
		nu.Claims = getClaimsIDs(usersClaims)
	}

	u := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = a.userRepo.Create(ctx, u); err != nil {
		if errors.Is(err, ErrUniqueEmail) {
			return User{}, ErrUniqueEmail
		}
		return User{}, err
	}

	if err = a.userRepo.CreateUserClaims(ctx, nu.Claims, u.ID); err != nil {
		return User{}, err
	}

	return u, nil
}

func getClaimsIDs(claims []Claims) []int {
	if len(claims) == 0 {
		return nil
	}
	result := make([]int, len(claims))

	for i := range claims {
		result[i] = claims[i].ID
	}

	return result
}
