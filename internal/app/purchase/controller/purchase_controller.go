package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
)

type purchaseController struct {
	purchaseService contracts.PurchaseService
}

func InitPurchaseController(router fiber.Router, purchaseService contracts.PurchaseService) {
	controller := purchaseController{
		purchaseService: purchaseService,
	}

	// jwt := jwt.Jwt

	// middleware := middlewares.NewMiddleware(jwt)

	purchaseRoute := router.Group("/v1/purchase")
	purchaseRoute.Post("/", controller.Purchase)
	purchaseRoute.Post("/:purchaseId", controller.UploadPayment)
}

func (mc *purchaseController) Purchase(ctx *fiber.Ctx) error {
	var req dto.PurchaseRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := mc.purchaseService.Purchase(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (mc *purchaseController) UploadPayment(ctx *fiber.Ctx) error {
	var requestBody dto.UploadPaymentRequest
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	purchaseId := ctx.Params("purchaseId")

	err := mc.purchaseService.UploadPayment(ctx.Context(), requestBody, purchaseId)
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusCreated)
}
