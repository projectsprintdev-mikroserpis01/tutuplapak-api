package entity

import "database/sql"

type User struct {
	ID                int            `db:"id" json:"id"`
	Email             sql.NullString `db:"email" json:"email"`
	Password          string         `db:"password" json:"-"`
	Phone             sql.NullString `db:"phone" json:"phone"`
	BankAccountName   sql.NullString `db:"bank_account_name" json:"bankAccoutnName"`
	BankAccountHolder sql.NullString `db:"bank_account_holder" json:"bankAccountHolder"`
	BankAccountNumber sql.NullString `db:"bank_account_number" json:"bankAccountNumber"`
	FileID            sql.NullInt16  `db:"file_id" json:"fileId"`
	FileURI           sql.NullString `db:"file_uri" json:"fileUri"`
	FileThumbnailURI  sql.NullString `db:"file_thumbnail_uri" json:"fileThumbnailUri"`
	CreatedAt         string         `db:"created_at" json:"createdAt"`
}
