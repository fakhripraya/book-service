package handlers

import (
	"net/http"

	"github.com/fakhripraya/book-service/data"
	"github.com/fakhripraya/book-service/database"
	"github.com/fakhripraya/book-service/entities"
)

// AddBook is a method to add the new given book info to the database
func (bookHandler *BookHandler) AddBook(rw http.ResponseWriter, r *http.Request) {

	// get the book via context
	bookReq := r.Context().Value(KeyBook{}).(*entities.TransactionKostBook)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := bookHandler.book.GetCurrentUser(rw, r, bookHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

}
