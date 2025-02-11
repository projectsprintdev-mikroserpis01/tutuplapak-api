package contracts

import (
	"context"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindByEmailOrPhone(ctx context.Context, email, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type UserService interface {
	GetUser(ctx context.Context, id int) (*dto.GetUserResponse, error)
	UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	LinkEmail(ctx context.Context, id int, req *dto.LinkEmailRequest) (*dto.LinkEmailResponse, error)
	LinkPhone(ctx context.Context, id int, req *dto.LinkPhoneRequest) (*dto.LinkPhoneResponse, error)
}
