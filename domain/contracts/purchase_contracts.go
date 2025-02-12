package contracts

import (
	"context"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
)

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, purchasedItems []entity.PurchaseItem, senderName string, senderContactType string, senderContactDetail string) (int64, error)
	DecreaseQuantity(ctx context.Context, productId int, quantity int) error
	GetProductById(ctx context.Context, productId int) (entity.Product, error)
	GetSellerById(ctx context.Context, sellerId int) (entity.DummyUser, error)
	GetPurchaseById(ctx context.Context, purchaseId int) (entity.Purchase, error)
}

type PurchaseService interface {
	Purchase(ctx context.Context, req dto.PurchaseRequest) (dto.PurchaseResponse, error)
	UploadPayment(ctx context.Context, req dto.UploadPaymentRequest, purchaseId string) error
}
