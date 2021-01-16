package entities

// ApprovalRoomBook is an entity to communicate with the ApprovalRoomBook client side
type ApprovalRoomBook struct {
	BookID       uint `json:"book_id"`
	FlagApproval bool `json:"flag_approval"`
}
