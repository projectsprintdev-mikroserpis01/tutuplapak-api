package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/middlewares"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/jwt"
)

type userController struct {
	service contracts.UserService
}

func InitUserController(router fiber.Router, service contracts.UserService, middleware *middlewares.Middleware) {
	controller := &userController{
		service,
	}

	userRouter := router.Group("/user")

	userRouter.Get("/", middleware.RequireAuth(), controller.getUser)
	userRouter.Put("/", middleware.RequireAuth(), controller.updateUser)
	userRouter.Post("/link/email", middleware.RequireAuth(), controller.linkEmail)
	userRouter.Post("/link/phone", middleware.RequireAuth(), controller.linkPhone)
}

func (c *userController) getUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("claims").(jwt.Claims).UserID

	res, err := c.service.GetUser(ctx.Context(), userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *userController) updateUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("claims").(jwt.Claims).UserID

	var req dto.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.UpdateUser(ctx.Context(), userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *userController) linkEmail(ctx *fiber.Ctx) error {
	userID := ctx.Locals("claims").(jwt.Claims).UserID

	var req dto.LinkEmailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.LinkEmail(ctx.Context(), userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *userController) linkPhone(ctx *fiber.Ctx) error {
	userID := ctx.Locals("claims").(jwt.Claims).UserID

	var req dto.LinkPhoneRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := c.service.LinkPhone(ctx.Context(), userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
