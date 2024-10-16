package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/PavelDonchenko/ecommerce-micro/auth/service"
	"github.com/PavelDonchenko/ecommerce-micro/common/database/sqldb"
	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
)

type UserRepo struct {
	db  sqlx.ExtContext
	log *logger.Logger
}

func NewUser(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) NewWithTx(tx sqldb.CommitRollbacker, log *logger.Logger) (service.UserRepo, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := UserRepo{
		log: log,
		db:  ec,
	}

	return &store, nil
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (service.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}
	q := `SELECT id,
			name,
			email,
			password_hash,
			created_at,
			updated_at,
			deleted_at
          FROM users.users WHERE email = :email`
	var userDB User
	var err error
	if err = sqldb.NamedQueryStruct(ctx, ur.log, ur.db, q, data, &userDB); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return service.User{}, service.ErrNotFound
		}
		return service.User{}, fmt.Errorf("db GetByEmail: %w", err)
	}

	return toServiceUser(userDB), nil
}

func (ur *UserRepo) Create(ctx context.Context, user service.User) error {
	q := `INSERT INTO users.users ( 
            id,
			name,
			email,
			password_hash,
            created_at,
            updated_at       
			) VALUES (
			:id,
			:name,
			:email,
			:password_hash,
			:created_at,
            :updated_at)`

	if err := sqldb.NamedExecContext(ctx, ur.log, ur.db, q, toDBUser(user)); err != nil {
		if errors.Is(err, sqldb.ErrDBDuplicatedEntry) {
			return service.ErrUniqueEmail
		}
		return fmt.Errorf("db Create user: %w", err)
	}

	return nil
}

func (ur *UserRepo) CreateUserClaims(ctx context.Context, claims []int, uid uuid.UUID) error {
	cq := `INSERT INTO users.users_claims(user_id, claim_id) VALUES (:user_id, :claim_id)`
	var userClaim struct {
		UserID  uuid.UUID `db:"user_id"`
		ClaimID int       `db:"claim_id"`
	}

	for _, claim := range claims {
		userClaim.UserID = uid
		userClaim.ClaimID = claim

		if err := sqldb.NamedExecContext(ctx, ur.log, ur.db, cq, userClaim); err != nil {
			return fmt.Errorf("db Create claims: %w", err)
		}
	}

	return nil
}

func (ur *UserRepo) GetClaimsByType(ctx context.Context, claimType string) ([]service.Claims, error) {
	data := struct {
		ClaimType string `db:"type"`
	}{
		ClaimType: claimType,
	}

	q := `SELECT id, type, value FROM users.claims WHERE type = :type`

	var claims []Claims

	if err := sqldb.NamedQuerySlice(ctx, ur.log, ur.db, q, data, &claims); err != nil {
		return nil, fmt.Errorf("db GetClaimsByType: %w", err)
	}

	return toServiceClaims(claims), nil
}
