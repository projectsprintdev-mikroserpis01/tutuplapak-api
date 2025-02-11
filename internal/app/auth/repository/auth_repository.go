package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) contracts.AuthRepository {
	return &authRepository{db}
}

// FindByEmail is a method to find a user by email
func (r *authRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByPhone is a method to find a user by phone
func (r *authRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE phone = ?", phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// RegisterWithEmail is a method to register a user with email
func (r *authRepository) RegisterWithEmail(ctx context.Context, user *entity.User) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO users (email, password) VALUES (:email, :password)", user)
	if err != nil {
		return err
	}

	return nil
}

// RegisterWithPhone is a method to register a user with phone
func (r *authRepository) RegisterWithPhone(ctx context.Context, user *entity.User) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO users (phone, password) VALUES (:phone, :password)", user)
	if err != nil {
		return err
	}

	return nil
}
