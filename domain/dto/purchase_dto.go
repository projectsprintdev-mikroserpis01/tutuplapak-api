package dto

import "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/domain/entity"

// PurchaseItem represents an item in the purchase
type PurchaseItem struct {
	ProductID        int     `json:"productId"`
	Name             string  `json:"name"`
	Category         int     `json:"category"`
	Quantity         int     `json:"qty"`
	Price            float64 `json:"price"`
	SKU              string  `json:"sku"`
	FileID           string  `json:"fileId"`
	FileURL          string  `json:"fileUri"`
	FileThumbnailURL string  `json:"fileThumbnailUri"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

// PaymentDetail represents payment details for a seller
type PaymentDetail struct {
	BankAccountName   string  `json:"bankAccountName"`
	BankAccountHolder string  `json:"bankAccountHolder"`
	BankAccountNumber string  `json:"bankAccountNumber"`
	TotalPrice        float64 `json:"totalPrice"`
}

type PurchaseRequest struct {
	PurchasedItems []struct {
		ProductID string `json:"product_id" validate:"required"`
		Qty       int    `json:"qty" validate:"required,min=2"`
	} `json:"purchased_items" validate:"required,min=1"`
	SenderName          string `json:"sender_name" validate:"required,min=4,max=55"`
	SenderContactType   string `json:"sender_contact_type" validate:"required,oneof=email phone"`
	SenderContactDetail string `json:"sender_contact_detail" validate:"required"`
}

// PurchaseResponse represents the response for a purchase
type PurchaseResponse struct {
	PurchaseID     string                `json:"purchaseId"`
	PurchasedItems []entity.PurchaseItem `json:"purchasedItems"`
	TotalPrice     float64               `json:"totalPrice"`
	PaymentDetails []PaymentDetail       `json:"paymentDetails"`
}

type UploadPaymentRequest struct {
	FileIDs []string `json:"file_ids" validate:"required,min=1"`
}
