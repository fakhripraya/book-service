package database

import "time"

// DBTransaction is an entity that directly communicate with the Transaction table in the database
type DBTransaction struct {
	ID             uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	TrxReferenceID uint      `gorm:"not null" json:"trx_reference_id"`
	TrxCategory    uint      `gorm:"not null" json:"trx_category"` // kategori transaksi (bayar kost, bayar perpanjang, dll)
	PaidOff        float64   `gorm:"not null" json:"paid_off"`
	MustPay        float64   `gorm:"not null" json:"must_pay"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	Created        time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy      string    `json:"created_by"`
	Modified       time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy     string    `json:"modified_by"`
}

// DBTransactionDetail is an entity that directly communicate with the TransactionDetail table in the database
type DBTransactionDetail struct {
	ID              uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	TrxID           uint      `gorm:"not null" json:"trx_id"`
	PaymentMethodID uint      `gorm:"not null" json:"payment_method_id"`
	Payment         float64   `gorm:"not null" json:"Payment"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	Created         time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy       string    `json:"created_by"`
	Modified        time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy      string    `json:"modified_by"`
}
