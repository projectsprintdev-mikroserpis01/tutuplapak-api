package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) contracts.UserRepository {
	return &userRepository{db}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmailOrPhone implements contracts.UserRepository.
func (u *userRepository) FindByEmailOrPhone(ctx context.Context, email string, phone string) (*entity.User, error) {
	var user entity.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ? OR phone = ?", email, phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID implements contracts.UserRepository.
func (u *userRepository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByPhone implements contracts.UserRepository.
func (u *userRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE phone = ?", phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update implements contracts.UserRepository.
func (u *userRepository) Update(ctx context.Context, user *entity.User) error {
	_, err := u.db.NamedExecContext(ctx, `
		UPDATE users
		SET email = :email, phone = :phone, password = :password,
			bank_account_number = :bank_account_number, bank_account_name = :bank_account_name, bank_account_holder = :bank_account_holder,
			file_id = :file_id, file_uri = :file_uri, file_thumbnail_uri = :file_thumbnail_uri
		WHERE id = :id
	`, user)
	if err != nil {
		return err
	}

	return nil
}
