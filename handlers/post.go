package handlers

import (
	"net/http"
	"time"

	"github.com/fakhripraya/book-service/config"
	"github.com/fakhripraya/book-service/data"
	"github.com/fakhripraya/book-service/database"
	"github.com/fakhripraya/book-service/entities"
	"gorm.io/gorm"
)

// AddBook is a method to add the new given book info to the database
func (bookHandler *BookHandler) AddBook(rw http.ResponseWriter, r *http.Request) {

	// get the book via context
	bookReq := r.Context().Value(KeyBook{}).(*entities.TransactionRoomBook)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := bookHandler.book.GetCurrentUser(rw, r, bookHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// proceed to create the new book with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var newBook database.DBTransactionRoomBook
		var kostTarget database.DBKost
		var dbErr error

		if dbErr = config.DB.Where("kost_id = ?", newBook.KostID).First(&kostTarget).Error; err != nil {
			return dbErr
		}

		newBook.BookerID = currentUser.ID
		newBook.KostID = bookReq.KostID
		newBook.RoomID = bookReq.RoomID
		newBook.RoomDetailID = bookReq.RoomDetailID
		newBook.PeriodID = bookReq.PeriodID
		newBook.Status = 0 // status 0 = baru // TODO: create a documented status later
		newBook.BookCode, dbErr = bookHandler.book.GenerateCode("K", kostTarget.Country[0:1], kostTarget.City[0:1])

		if dbErr != nil {
			return dbErr
		}

		newBook.BookDate = bookReq.BookDate
		newBook.IsActive = true
		newBook.Created = time.Now().Local()
		newBook.CreatedBy = currentUser.Username
		newBook.Modified = time.Now().Local()
		newBook.ModifiedBy = currentUser.Username

		if dbErr = tx.Create(&newBook).Error; dbErr != nil {
			return dbErr
		}

		// add the room book member to the database
		dbErr = bookHandler.book.AddRoomBookMember(currentUser, newBook.ID, &bookReq.TrxBookMember)

		if dbErr != nil {
			return dbErr
		}

		// add the base transaction to the database
		dbErr = bookHandler.book.AddTransaction(currentUser, newBook.ID, 0, bookReq.PaymentMethodID, bookReq.Payment, bookReq.MustPay) // TODO: 0 adalah kategori, bikin dokumentasi ntr

		if dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
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
