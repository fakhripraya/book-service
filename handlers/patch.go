package handlers

import (
	"net/http"

	"github.com/fakhripraya/book-service/config"
	"github.com/fakhripraya/book-service/data"
	"github.com/fakhripraya/book-service/database"
	"github.com/fakhripraya/book-service/entities"
	"gorm.io/gorm"
)

// ApprovalBookTransaction is a method to approve the book transaction info
func (bookHandler *BookHandler) ApprovalBookTransaction(rw http.ResponseWriter, r *http.Request) {

	// get the approval via context
	approvalReq := r.Context().Value(KeyApproval{}).(*entities.ApprovalRoomBook)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := bookHandler.book.GetCurrentUser(rw, r, bookHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// proceed to create the new approval with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var targetBook database.DBTransactionRoomBook
		var targetTransaction database.DBTransaction
		var targetTransactionDetail database.DBTransactionDetail
		var dbErr error

		if dbErr := config.DB.Where("id = ?", approvalReq.BookID).First(&targetBook).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		if dbErr := config.DB.Where("trx_reference_id = ?", targetBook.ID).First(&targetTransaction).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		if dbErr := config.DB.Where("trx_id = ?", targetTransaction.ID).First(&targetTransactionDetail).Error; dbErr != nil {
			return dbErr
		}

		// Status 1 = approved
		if approvalReq.FlagApproval == true {
			targetBook.Status = 1
		}

		// Status 1 = approved
		if approvalReq.FlagApproval == true {
			targetTransaction.PaidOff = targetTransaction.PaidOff + targetTransactionDetail.Payment
		}

		// add the base transaction to the database
		dbErr = bookHandler.book.UpdateTransaction(currentUser, &targetTransaction)

		if dbErr != nil {
			return dbErr
		}

		// Status 1 = approved
		// Status 2 = reject
		if approvalReq.FlagApproval == true {
			targetTransactionDetail.Status = 1
		} else {
			targetTransactionDetail.Status = 2
		}

		// add the base transaction to the database
		dbErr = bookHandler.book.UpdateTransactionDetail(currentUser, &targetTransactionDetail)

		if dbErr != nil {
			return dbErr
		}

		return nil

	})

	// if transaction error
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericError{Message: "Sukses request booking"}, rw)
	return

}
