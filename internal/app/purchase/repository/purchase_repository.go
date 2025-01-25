package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"
)

type purchaseRepository struct {
	db *sqlx.DB
}

func NewPurchaseRepository(db *sqlx.DB) contracts.PurchaseRepository {
	return &purchaseRepository{
		db: db,
	}
}

func (r *purchaseRepository) CreatePurchase(ctx context.Context, purchasedItems []entity.PurchaseItem, senderName string, senderContactType string, senderContactDetail string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO purchase (purchased_items, sender_name, sender_contact_type, sender_contact_detail) VALUES ($1, $2, $3, $4) RETURNING id", purchasedItems, senderName, senderContactType, senderContactDetail)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *purchaseRepository) DecreaseQuantity(ctx context.Context, productId int, quantity int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE products SET qty = qty - $1 WHERE id = $2", quantity, productId)
	if err != nil {
		return err
	}
	return nil
}

func (r *purchaseRepository) GetProductById(ctx context.Context, productId int) (entity.Product, error) {
	var product entity.Product
	err := r.db.GetContext(ctx, &product, "SELECT * FROM products WHERE id=$1", productId)
	return product, err
}

func (r *purchaseRepository) GetSellerById(ctx context.Context, sellerId int) (entity.DummyUser, error) {
	var seller entity.DummyUser
	err := r.db.GetContext(ctx, &seller, "SELECT * FROM users WHERE id=$1", sellerId)
	return seller, err
}

func (r *purchaseRepository) GetPurchaseById(ctx context.Context, purchaseId int) (entity.Purchase, error) {
	var purchase entity.Purchase
	err := r.db.GetContext(ctx, &purchase, "SELECT * FROM purchase WHERE id=$1", purchaseId)
	return purchase, err
}
