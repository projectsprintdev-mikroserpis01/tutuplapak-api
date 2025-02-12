package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
)

type authController struct {
	service contracts.AuthService
}

func InitAuthController(router fiber.Router, service contracts.AuthService) {
	controller := &authController{
		service,
	}

	router.Post("/login/email", controller.loginWithEmail)
	router.Post("/login/phone", controller.loginWithPhone)
	router.Post("/register/email", controller.registerWithEmail)
	router.Post("/register/phone", controller.registerWithPhone)
}

func (c *authController) loginWithEmail(ctx *fiber.Ctx) error {
	var req dto.LoginWithEmailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.LoginWithEmail(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) loginWithPhone(ctx *fiber.Ctx) error {
	var req dto.LoginWithPhoneRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.LoginWithPhone(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *authController) registerWithEmail(ctx *fiber.Ctx) error {
	var req dto.RegisterWithEmailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.RegisterWithEmail(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (c *authController) registerWithPhone(ctx *fiber.Ctx) error {
	var req dto.RegisterWithPhoneRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.RegisterWithPhone(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}
