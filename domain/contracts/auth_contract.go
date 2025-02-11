package contracts

import (
	"context"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	RegisterWithEmail(ctx context.Context, user *entity.User) error
	RegisterWithPhone(ctx context.Context, user *entity.User) error
}

type AuthService interface {
	LoginWithEmail(ctx context.Context, req *dto.LoginWithEmailRequest) (*dto.LoginWithEmailResponse, error)
	LoginWithPhone(ctx context.Context, req *dto.LoginWithPhoneRequest) (*dto.LoginWithPhoneResponse, error)
	RegisterWithEmail(ctx context.Context, req *dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, error)
	RegisterWithPhone(ctx context.Context, req *dto.RegisterWithPhoneRequest) (*dto.RegisterWithPhoneResponse, error)
}
