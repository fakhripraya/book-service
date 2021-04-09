package entities

import (
	"time"

	"github.com/fakhripraya/book-service/database"
)

// TransactionRoomBook is an entity to communicate with the TransactionRoomBook client side
type TransactionRoomBook struct {
	ID               uint                                   `json:"id"`
	BookerID         uint                                   `json:"booker_id"`
	KostID           uint                                   `json:"kost_id"`
	RoomID           uint                                   `json:"room_id"`
	RoomDetailID     uint                                   `json:"room_detail_id"`
	PaymentMethodID  uint                                   `json:"payment_method_id"`
	PeriodID         uint                                   `json:"period_id"`
	Status           uint                                   `json:"status"`
	BookDate         time.Time                              `json:"book_date"`
	Payment          float64                                `json:"Payment"`
	MustPay          float64                                `json:"must_pay"`
	Members          []database.DBTransactionRoomBookMember `json:"members"`
	VerificationData TransactionVerification                `json:"verification_data"`
	IsActive         bool                                   `json:"is_active"`
	Created          time.Time                              `json:"created"`
	CreatedBy        string                                 `json:"created_by"`
	Modified         time.Time                              `json:"modified"`
	ModifiedBy       string                                 `json:"modified_by"`
}

// TransactionRoomBookMember is an entity to communicate with the TransactionRoomBookMember client side
type TransactionRoomBookMember struct {
	ID         uint      `json:"id"`
	RoomBookID uint      `json:"room_book_id"`
	MemberName string    `json:"member_name"`
	Phone      string    `json:"phone"`
	Gender     bool      `json:"gender"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// TransactionVerification is an entity to communicate with the TransactionVerification client side
type TransactionVerification struct {
	ID          uint      `json:"id"`
	ReferenceID uint      `json:"reference_id"`
	PictDesc    string    `json:"pict_desc"`
	URL         string    `json:"url"`
	IsActive    bool      `json:"is_active"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Modified    time.Time `json:"modified"`
	ModifiedBy  string    `json:"modified_by"`
}
