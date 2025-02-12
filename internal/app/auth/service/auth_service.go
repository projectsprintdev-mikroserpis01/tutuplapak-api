package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/jwt"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/validator"
)

type authService struct {
	repo      contracts.AuthRepository
	validator validator.ValidatorInterface
	bcrypt    bcrypt.BcryptInterface
	jwt       jwt.JwtInterface
}

func NewAuthService(repo contracts.AuthRepository, validator validator.ValidatorInterface, bcrypt bcrypt.BcryptInterface, jwt jwt.JwtInterface) contracts.AuthService {
	return &authService{
		repo,
		validator,
		bcrypt,
		jwt,
	}
}

// LoginWithEmail is a method to login with email
func (s *authService) LoginWithEmail(ctx context.Context, req *dto.LoginWithEmailRequest) (*dto.LoginWithEmailResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "email not found")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	isCorrect := s.bcrypt.Compare(req.Password, user.Password)
	if !isCorrect {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid email or password")
	}

	token, err := s.jwt.Create(user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	email := ""
	if user.Email.Valid {
		email = user.Email.String
	}

	phone := ""
	if user.Phone.Valid {
		phone = user.Phone.String
	}

	res := &dto.LoginWithEmailResponse{
		Email: email,
		Phone: phone,
		Token: token,
	}

	return res, nil
}

// LoginWithPhone is a method to login with phone
func (s *authService) LoginWithPhone(ctx context.Context, req *dto.LoginWithPhoneRequest) (*dto.LoginWithPhoneResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	user, err := s.repo.FindByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "phone not found")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	isCorrect := s.bcrypt.Compare(req.Password, user.Password)
	if !isCorrect {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid phone or password")
	}

	token, err := s.jwt.Create(user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	email := ""
	if user.Email.Valid {
		email = user.Email.String
	}

	phone := ""
	if user.Phone.Valid {
		phone = user.Phone.String
	}

	res := &dto.LoginWithPhoneResponse{
		Email: email,
		Phone: phone,
		Token: token,
	}

	return res, nil
}

// RegisterWithEmail is a method to register with email
func (s *authService) RegisterWithEmail(ctx context.Context, req *dto.RegisterWithEmailRequest) (*dto.RegisterWithEmailResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return nil, fiber.NewError(fiber.StatusConflict, "email already exists")
	}

	hashedPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user := &entity.User{
		Email:    sql.NullString{String: req.Email, Valid: true},
		Password: hashedPassword,
	}

	err = s.repo.RegisterWithEmail(ctx, user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	token, err := s.jwt.Create(user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &dto.RegisterWithEmailResponse{
		Email: user.Email.String,
		Phone: "",
		Token: token,
	}

	return res, nil
}

// RegisterWithPhone is a method to register with phone
func (s *authService) RegisterWithPhone(ctx context.Context, req *dto.RegisterWithPhoneRequest) (*dto.RegisterWithPhoneResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, valErr.Error())
	}

	_, err := s.repo.FindByPhone(ctx, req.Phone)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return nil, fiber.NewError(fiber.StatusConflict, "phone already exists")
	}

	hashedPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user := &entity.User{
		Phone:    sql.NullString{String: req.Phone, Valid: true},
		Password: hashedPassword,
	}

	err = s.repo.RegisterWithPhone(ctx, user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	token, err := s.jwt.Create(user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := &dto.RegisterWithPhoneResponse{
		Email: "",
		Phone: user.Phone.String,
		Token: token,
	}

	return res, nil
}
