package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fakhripraya/book-service/config"
	"github.com/fakhripraya/book-service/data"
	"github.com/fakhripraya/book-service/database"
	"github.com/fakhripraya/book-service/entities"
	"gorm.io/gorm"
)

// OwnerApprovalBookTransaction is a method to approve the book transaction info by the owner
func (bookHandler *BookHandler) OwnerApprovalBookTransaction(rw http.ResponseWriter, r *http.Request) {

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
		var dbErr error

		// look for the requested book
		if dbErr := config.DB.Where("id = ?", approvalReq.BookID).First(&targetBook).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		bookedKost := &database.DBKost{}
		if dbErr := config.DB.Where("id = ?", targetBook.KostID).First(&bookedKost).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		// occurs when transaction already been approved by the owner
		if targetBook.Status != 0 {
			rw.WriteHeader(http.StatusForbidden)

			return fmt.Errorf("Status booking tidak valid untuk di approve")
		}

		// only owner can approve the book transaction in this method
		if currentUser.ID != bookedKost.OwnerID {
			rw.WriteHeader(http.StatusForbidden)

			return fmt.Errorf("Hanya owner kost yang bisa approve book ini")
		}

		// TODO: buat dokumentasi
		// Status 1 = approved by owner
		// Status 3 = rejected
		if approvalReq.FlagApproval == true {
			targetBook.Status = 1
			targetBook.Modified = time.Now().Local()
			targetBook.ModifiedBy = currentUser.Username
		} else {
			targetBook.Status = 3
			targetBook.Modified = time.Now().Local()
			targetBook.ModifiedBy = currentUser.Username
		}

		// update the room book
		dbErr = config.DB.Save(targetBook).Error

		if dbErr != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		return nil

	})

	// if transaction error
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// TODO: send notification

	if approvalReq.FlagApproval == true {
		data.ToJSON(&GenericError{Message: "Sukses Approve booking"}, rw)
	} else {
		data.ToJSON(&GenericError{Message: "Sukses Reject booking"}, rw)
	}

	return

}

// TenantApprovalBookTransaction is a method to approve the book transaction info by the tenant
func (bookHandler *BookHandler) TenantApprovalBookTransaction(rw http.ResponseWriter, r *http.Request) {

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

		// look for the requested book
		if dbErr := config.DB.Where("id = ?", approvalReq.BookID).First(&targetBook).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		// occurs when transaction already been approved by the tenant
		if targetBook.Status != 1 {
			rw.WriteHeader(http.StatusForbidden)

			return fmt.Errorf("Status booking tidak valid untuk di approve")
		}

		// only tenant can approve the book transaction in this method
		if currentUser.ID != targetBook.BookerID {
			rw.WriteHeader(http.StatusForbidden)

			return fmt.Errorf("Hanya tenant kost yang bisa approve book ini")
		}

		// look for the base transaction
		if dbErr := config.DB.Where("trx_reference_id = ?", targetBook.ID).First(&targetTransaction).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		// look for the base transaction detail
		if dbErr := config.DB.Where("trx_id = ?", targetTransaction.ID).First(&targetTransactionDetail).Error; dbErr != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		// TODO: buat dokumentasi
		// Status 2 = approved by user
		// Status 3 = rejected
		if approvalReq.FlagApproval == true {
			targetBook.Status = 2
			targetBook.Modified = time.Now().Local()
			targetBook.ModifiedBy = currentUser.Username
		} else {
			targetBook.Status = 3
			targetBook.Modified = time.Now().Local()
			targetBook.ModifiedBy = currentUser.Username
		}

		// update the room book
		dbErr = config.DB.Save(targetBook).Error

		if dbErr != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		// TODO: tanyain lagi booknya gmn
		if approvalReq.FlagApproval == true {
			targetTransaction.PaidOff = targetTransaction.PaidOff + targetTransactionDetail.Payment
		}

		// add the base transaction to the database
		dbErr = bookHandler.book.UpdateTransaction(currentUser, &targetTransaction)

		if dbErr != nil {
			rw.WriteHeader(http.StatusBadRequest)

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
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		return nil

	})

	// if transaction error
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// TODO: send notification

	// send status ok if reach this point
	rw.WriteHeader(http.StatusOK)
	if approvalReq.FlagApproval == true {
		data.ToJSON(&GenericError{Message: "Sukses Approve booking"}, rw)
	} else {
		data.ToJSON(&GenericError{Message: "Sukses Reject booking"}, rw)
	}

	return

}
