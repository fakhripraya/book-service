package entities

import (
	"time"
)

// TransactionKostBook is an entity to communicate with the TransactionKostBook client side
type TransactionKostBook struct {
	ID         uint      `json:"id"`
	BookerID   uint      `json:"booker_id"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}
