package database

import "time"

// MasterPaymentMethod is an entity that directly communicate with the MasterPaymentMethod table in the database
type MasterPaymentMethod struct {
	ID          uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	PaymentType string    `gorm:"not null" json:"payment_type"` // virtual or physics
	PaymentDesc string    `gorm:"not null" json:"payment_desc"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	Created     time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy   string    `json:"created_by"`
	Modified    time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy  string    `json:"modified_by"`
}

// MasterPaymentMethodTable set the migrated struct table name
func (masterPaymentMethod *MasterPaymentMethod) MasterPaymentMethodTable() string {
	return "dbMasterPaymentMethod"
}
