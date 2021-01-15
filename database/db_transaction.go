package database

import "time"

// DBTransaction is an entity that directly communicate with the Transaction table in the database
type DBTransaction struct {
	ID             uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	TrxReferenceID uint      `gorm:"not null" json:"trx_reference_id"`
	CategoryID     uint      `gorm:"not null" json:"category_id"`
	Created        time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy      string    `json:"created_by"`
	Modified       time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy     string    `json:"modified_by"`
}

// DBTransactionDetail is an entity that directly communicate with the TransactionDetail table in the database
type DBTransactionDetail struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	TrxID      uint      `gorm:"not null" json:"trx_id"`
	StatusID   uint      `gorm:"not null" json:"status_id"`
	PaidOff    uint64    `gorm:"not null" json:"paid_off"`
	MustPay    uint64    `gorm:"not null" json:"must_pay"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}
