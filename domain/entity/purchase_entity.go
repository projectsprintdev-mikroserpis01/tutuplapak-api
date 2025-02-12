package entity

import "time"

// Purchase represents the "purchase" table
type Purchase struct {
	ID                  int            `db:"id"`
	PurchasedItems      []PurchaseItem `db:"purchased_items"`
	SenderName          string         `db:"sender_name"`
	SenderContactType   string         `db:"sender_contact_type"`
	SenderContactDetail string         `db:"sender_contact_detail"`
	PaymentProofIDs     []string       `db:"payment_proof_ids"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
}

// PurchaseItem represents an item in the "purchased_items" JSONB array in the "purchase" table
type PurchaseItem struct {
	ProductID        int     `db:"product_id"`
	Name             string  `db:"name"`
	Category         int     `db:"category"`
	Quantity         int     `db:"qty"`
	Price            float64 `db:"price"`
	SKU              string  `db:"sku"`
	FileID           string  `db:"file_id"`
	FileURL          string  `db:"file_url"`
	FileThumbnailURL string  `db:"file_thumbnail_url"`
}
