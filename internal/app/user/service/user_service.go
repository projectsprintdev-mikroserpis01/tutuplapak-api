package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/validator"
)

type userService struct {
	repo      contracts.UserRepository
	validator validator.ValidatorInterface
}

func NewUserService(repo contracts.UserRepository, validator validator.ValidatorInterface) contracts.UserService {
	return &userService{
		repo,
		validator,
	}
}

// GetUser implements contracts.UserService.
func (u *userService) GetUser(ctx context.Context, id int) (*dto.GetUserResponse, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &dto.GetUserResponse{
		Email: func() string {
			if user.Email.Valid {
				return user.Email.String
			}
			return ""
		}(),
		Phone: func() string {
			if user.Phone.Valid {
				return user.Phone.String
			}
			return ""
		}(),
		FileID: func() string {
			if user.FileID.Valid {
				return strconv.Itoa(int(user.FileID.Int16))
			}
			return ""
		}(),
		FileURI: func() string {
			if user.FileURI.Valid {
				return user.FileURI.String
			}
			return ""
		}(),
		FileThumbnailURI: func() string {
			if user.FileThumbnailURI.Valid {
				return user.FileThumbnailURI.String
			}
			return ""
		}(),
		BankAccountName: func() string {
			if user.BankAccountName.Valid {
				return user.BankAccountName.String
			}
			return ""
		}(),
		BankAccountHolder: func() string {
			if user.BankAccountHolder.Valid {
				return user.BankAccountHolder.String
			}
			return ""
		}(),
		BankAccountNumber: func() string {
			if user.BankAccountNumber.Valid {
				return user.BankAccountNumber.String
			}
			return ""
		}(),
	}

	return res, nil
}

// LinkEmail implements contracts.UserService.
func (u *userService) LinkEmail(ctx context.Context, id int, req *dto.LinkEmailRequest) (*dto.LinkEmailResponse, error) {
	valErr := u.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user.Email = sql.NullString{
		String: req.Email,
		Valid:  true,
	}

	err = u.repo.Update(ctx, user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &dto.LinkEmailResponse{
		Email: func() string {
			if user.Email.Valid {
				return user.Email.String
			}
			return ""
		}(),
		Phone: func() string {
			if user.Phone.Valid {
				return user.Phone.String
			}
			return ""
		}(),
		FileID: func() string {
			if user.FileID.Valid {
				return strconv.Itoa(int(user.FileID.Int16))
			}
			return ""
		}(),
		FileURI: func() string {
			if user.FileURI.Valid {
				return user.FileURI.String
			}
			return ""
		}(),
		FileThumbnailURI: func() string {
			if user.FileThumbnailURI.Valid {
				return user.FileThumbnailURI.String
			}
			return ""
		}(),
		BankAccountName: func() string {
			if user.BankAccountName.Valid {
				return user.BankAccountName.String
			}
			return ""
		}(),
		BankAccountHolder: func() string {
			if user.BankAccountHolder.Valid {
				return user.BankAccountHolder.String
			}
			return ""
		}(),
		BankAccountNumber: func() string {
			if user.BankAccountNumber.Valid {
				return user.BankAccountNumber.String
			}
			return ""
		}(),
	}

	return res, nil
}

// LinkPhone implements contracts.UserService.
func (u *userService) LinkPhone(ctx context.Context, id int, req *dto.LinkPhoneRequest) (*dto.LinkPhoneResponse, error) {
	valErr := u.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user.Phone = sql.NullString{
		String: req.Phone,
		Valid:  true,
	}

	err = u.repo.Update(ctx, user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &dto.LinkPhoneResponse{
		Email: func() string {
			if user.Email.Valid {
				return user.Email.String
			}
			return ""
		}(),
		Phone: func() string {
			if user.Phone.Valid {
				return user.Phone.String
			}
			return ""
		}(),
		FileID: func() string {
			if user.FileID.Valid {
				return strconv.Itoa(int(user.FileID.Int16))
			}
			return ""
		}(),
		FileURI: func() string {
			if user.FileURI.Valid {
				return user.FileURI.String
			}
			return ""
		}(),
		FileThumbnailURI: func() string {
			if user.FileThumbnailURI.Valid {
				return user.FileThumbnailURI.String
			}
			return ""
		}(),
		BankAccountName: func() string {
			if user.BankAccountName.Valid {
				return user.BankAccountName.String
			}
			return ""
		}(),
		BankAccountHolder: func() string {
			if user.BankAccountHolder.Valid {
				return user.BankAccountHolder.String
			}
			return ""
		}(),
		BankAccountNumber: func() string {
			if user.BankAccountNumber.Valid {
				return user.BankAccountNumber.String
			}
			return ""
		}(),
	}

	return res, nil
}

// UpdateUser implements contracts.UserService.
func (u *userService) UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	valErr := u.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	if req.FileID == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "fileId cannot be nil")
	}

	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if *req.FileID != "" {
		fileID, err := strconv.Atoi(*req.FileID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "fileId must be a number")
		}

		user.FileID = sql.NullInt16{
			Int16: int16(fileID),
			Valid: true,
		}
	}

	user.BankAccountHolder = sql.NullString{
		String: req.BankAccountHolder,
		Valid:  true,
	}
	user.BankAccountName = sql.NullString{
		String: req.BankAccountName,
		Valid:  true,
	}
	user.BankAccountNumber = sql.NullString{
		String: req.BankAccountNumber,
		Valid:  true,
	}

	err = u.repo.Update(ctx, user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &dto.UpdateUserResponse{
		Email: func() string {
			if user.Email.Valid {
				return user.Email.String
			}
			return ""
		}(),
		Phone: func() string {
			if user.Phone.Valid {
				return user.Phone.String
			}
			return ""
		}(),
		FileID: func() string {
			if user.FileID.Valid {
				return strconv.Itoa(int(user.FileID.Int16))
			}
			return ""
		}(),
		FileURI: func() string {
			if user.FileURI.Valid {
				return user.FileURI.String
			}
			return ""
		}(),
		FileThumbnailURI: func() string {
			if user.FileThumbnailURI.Valid {
				return user.FileThumbnailURI.String
			}
			return ""
		}(),
		BankAccountName: func() string {
			if user.BankAccountName.Valid {
				return user.BankAccountName.String
			}
			return ""
		}(),
		BankAccountHolder: func() string {
			if user.BankAccountHolder.Valid {
				return user.BankAccountHolder.String
			}
			return ""
		}(),
		BankAccountNumber: func() string {
			if user.BankAccountNumber.Valid {
				return user.BankAccountNumber.String
			}
			return ""
		}(),
	}

	return res, nil
}
