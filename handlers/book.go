package handlers

import (
	"github.com/fakhripraya/book-service/data"

	"github.com/hashicorp/go-hclog"
	"github.com/srinathgs/mysqlstore"
)

// KeyBook is a key used for the Book object in the context
type KeyBook struct{}

// BookHandler is a handler struct for book changes
type BookHandler struct {
	logger hclog.Logger
	book   *data.Book
	store  *mysqlstore.MySQLStore
}

// NewBookHandler returns a new book handler with the given logger
func NewBookHandler(newLogger hclog.Logger, newBook *data.Book, newStore *mysqlstore.MySQLStore) *BookHandler {
	return &BookHandler{newLogger, newBook, newStore}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}
