package database

import "time"

// DBTransactionRoomBook is an entity that directly communicate with the TransactionRoomBook table in the database
type DBTransactionRoomBook struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	KostID     uint      `gorm:"not null" json:"kost_id"`
	RoomID     uint      `gorm:"not null" json:"room_id"`
	BookerID   uint      `gorm:"not null" json:"booker_id"`
	StatusID   uint      `gorm:"not null" json:"status_id"`
	BookCode   string    `gorm:"not null" json:"book_code"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}
