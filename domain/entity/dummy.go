package entity

import "time"

// User represents the "users" table
type DummyUser struct {
	ID                int       `db:"id"`
	Email             string    `db:"email"`
	Password          string    `db:"password"`
	Name              string    `db:"name"`
	Phone             int       `db:"phone"`
	FileID            string    `db:"file_id"`
	FileURL           string    `db:"file_url"`
	FileThumbnailURL  string    `db:"file_thumbnail_url"`
	BankAccountName   string    `db:"bank_account_name"`
	BankAccountHolder string    `db:"bank_account_holder"`
	BankAccountNumber string    `db:"bank_account_number"`
	CreatedAt         time.Time `db:"created_at"`
}

// Product represents the "products" table
type Product struct {
	ID               int       `db:"id"`
	Name             string    `db:"name"`
	Category         int       `db:"category"`
	Quantity         int       `db:"qty"`
	Price            float64   `db:"price"`
	SKU              string    `db:"sku"`
	FileID           string    `db:"file_id"`
	FileURL          string    `db:"file_url"`
	FileThumbnailURL string    `db:"file_thumbnail_url"`
	UserID           int       `db:"user_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
