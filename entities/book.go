package entities

import (
	"time"
)

// TransactionRoomBook is an entity to communicate with the TransactionRoomBook client side
type TransactionRoomBook struct {
	ID              uint      `json:"id"`
	KostID          uint      `json:"kost_id"`
	RoomID          uint      `json:"room_id"`
	RoomDetailID    uint      `json:"room_detail_id"`
	BookerID        uint      `json:"booker_id"`
	PaymentMethodID uint      `json:"payment_method_id"`
	StatusID        uint      `json:"status_id"`
	BookCode        string    `json:"book_code"`
	Period          uint      `json:"period"`
	BookDate        time.Time `json:"book_date"`
	PaidOff         float64   `json:"paid_off"`
	MustPay         float64   `json:"must_pay"`
	PrevPayment     time.Time `json:"prev_payment"`
	NextPayment     time.Time `json:"next_payment"`
	IsActive        bool      `json:"is_active"`
	Created         time.Time `json:"created"`
	CreatedBy       string    `json:"created_by"`
	Modified        time.Time `json:"modified"`
	ModifiedBy      string    `json:"modified_by"`
}

// TransactionRoomBookMember is an entity to communicate with the TransactionRoomBookMember client side
type TransactionRoomBookMember struct {
	ID         uint                              `json:"id"`
	RoomBookID uint                              `json:"room_book_id"`
	TenantID   uint                              `json:"tenant_id"`
	Members    []TransactionRoomBookMemberDetail `json:"members"`
	IsActive   bool                              `json:"is_active"`
	Created    time.Time                         `json:"created"`
	CreatedBy  string                            `json:"created_by"`
	Modified   time.Time                         `json:"modified"`
	ModifiedBy string                            `json:"modified_by"`
}

// TransactionRoomBookMemberDetail is an entity to communicate with the TransactionRoomBookMemberDetail client side
type TransactionRoomBookMemberDetail struct {
	ID               uint      `json:"id"`
	RoomBookMemberID uint      `json:"room_book_member_id"`
	MemberName       string    `json:"member_name"`
	Phone            string    `json:"phone"`
	Gender           bool      `json:"gender"`
	IsActive         bool      `json:"is_active"`
	Created          time.Time `json:"created"`
	CreatedBy        string    `json:"created_by"`
	Modified         time.Time `json:"modified"`
	ModifiedBy       string    `json:"modified_by"`
}
