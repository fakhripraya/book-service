package database

import "time"

// DBTransactionKostBook is an entity that directly communicate with the TransactionKostBook table in the database
type DBTransactionKostBook struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	BookerID   uint      `gorm:"not null" json:"booker_id"`
	StatusID   uint      `gorm:"not null" json:"status_id"`
	BookCode   string    `gorm:"not null" json:"book_code"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}
