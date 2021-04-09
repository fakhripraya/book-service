package data

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/book-service/config"
	"github.com/fakhripraya/book-service/database"
	"github.com/fakhripraya/book-service/entities"
	"github.com/hashicorp/go-hclog"
	"github.com/srinathgs/mysqlstore"
	"gorm.io/gorm"
)

// Claims determine the current user token holder
type Claims struct {
	Username string
	jwt.StandardClaims
}

// Book defines a struct for book flow
type Book struct {
	logger hclog.Logger
}

// NewBook is a function to create new Book struct
func NewBook(newLogger hclog.Logger) *Book {
	return &Book{newLogger}
}

// GetCurrentUser will get the current user login info
func (book *Book) GetCurrentUser(rw http.ResponseWriter, r *http.Request, store *mysqlstore.MySQLStore) (*database.MasterUser, error) {

	// Get a session (existing/new)
	session, err := store.Get(r, "session-name")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return nil, err
	}

	// check the logged in user from the session
	// if user available, get the user info from the session
	if session.Values["userLoggedin"] == nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return nil, fmt.Errorf("Error 401")
	}

	// work with database
	// look for the current user logged in in the db
	var currentUser database.MasterUser
	if err := config.DB.Where("username = ?", session.Values["userLoggedin"].(string)).First(&currentUser).Error; err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return nil, err
	}

	return &currentUser, nil

}

// GenerateCode will generate the new given type code
func (book *Book) GenerateCode(codeType, country, city string) (string, error) {

	// generate 8 random crypted number
	var max int = 8
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	// returns the crypted random 8 number
	var crypted string = string(b)

	var finalCode string = codeType +
		"/" + country +
		"-" + city +
		"/" + strconv.Itoa(time.Now().UTC().Year()) + "-" + time.Now().UTC().Month().String()[0:1] +
		"/" + crypted

	return finalCode, nil

}

// AddTransaction is a function to add transaction based on the given transaction entry, this transaction is not scoped
func (book *Book) AddTransaction(currentUser *database.MasterUser, ReferenceID, TrxCategory uint, mustPay float64) (uint, error) {

	// set variables
	var newTransaction database.DBTransaction
	var dbErr error

	newTransaction.TrxReferenceID = ReferenceID
	newTransaction.TrxCategory = TrxCategory // kategori booking
	newTransaction.PaidOff = 0               // paid off masih 0 karna transaksi baru
	newTransaction.MustPay = mustPay
	newTransaction.IsActive = true
	newTransaction.Created = time.Now().Local()
	newTransaction.CreatedBy = currentUser.Username
	newTransaction.Modified = time.Now().Local()
	newTransaction.ModifiedBy = currentUser.Username

	// insert the new transaction to database
	if dbErr = config.DB.Create(&newTransaction).Error; dbErr != nil {
		return 0, dbErr
	}

	return newTransaction.ID, nil

}

// AddTransactionDetail is a function to add transaction detail based on the given transaction entry, this transaction is not scoped
func (book *Book) AddTransactionDetail(currentUser *database.MasterUser, status, trxID, PaymentMethodID uint, payment float64) error {

	// set variables
	var newTransactionDetail database.DBTransactionDetail
	var dbErr error

	newTransactionDetail.TrxID = trxID
	newTransactionDetail.PaymentMethodID = PaymentMethodID
	newTransactionDetail.Status = status
	newTransactionDetail.Payment = payment
	newTransactionDetail.IsActive = true
	newTransactionDetail.Created = time.Now().Local()
	newTransactionDetail.CreatedBy = currentUser.Username
	newTransactionDetail.Modified = time.Now().Local()
	newTransactionDetail.ModifiedBy = currentUser.Username

	// insert the new transaction detail to database
	if dbErr = config.DB.Create(&newTransactionDetail).Error; dbErr != nil {
		return dbErr
	}

	return nil

}

// UpdateTransaction is a function to update transaction based on the given transaction entry
func (book *Book) UpdateTransaction(currentUser *database.MasterUser, targetTransaction *database.DBTransaction) error {

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var dbErr error

		targetTransaction.Modified = time.Now().Local()
		targetTransaction.ModifiedBy = currentUser.Username

		// update the transaction
		dbErr = config.DB.Where("id = ?", targetTransaction.ID).Save(&targetTransaction).Error
		if dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {

		return err
	}

	return nil

}

// UpdateTransactionDetail is a function to update transaction detail based on the given transaction entry
func (book *Book) UpdateTransactionDetail(currentUser *database.MasterUser, targetTransactionDetail *database.DBTransactionDetail) error {

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var dbErr error

		targetTransactionDetail.Modified = time.Now().Local()
		targetTransactionDetail.ModifiedBy = currentUser.Username

		// update the transaction detail
		dbErr = config.DB.Save(&targetTransactionDetail).Error

		if dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {

		return err
	}

	return nil

}

// AddRoomBookMember is a function to add book member based on the given book entity
func (book *Book) AddRoomBookMember(currentUser *database.MasterUser, roomBookID uint, targetRoomBookMember []database.DBTransactionRoomBookMember) error {

	// add the room book member to the database with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variable
		var dbErr error
		var newRoomBookMember = targetRoomBookMember

		// add the room book id to the slices
		for i := range newRoomBookMember {
			(&newRoomBookMember[i]).RoomBookID = roomBookID
			(&newRoomBookMember[i]).IsActive = true
			(&newRoomBookMember[i]).Created = time.Now().Local()
			(&newRoomBookMember[i]).CreatedBy = currentUser.Username
			(&newRoomBookMember[i]).Modified = time.Now().Local()
			(&newRoomBookMember[i]).ModifiedBy = currentUser.Username
		}

		// insert the new room book member to database
		if dbErr = tx.Create(&newRoomBookMember).Error; dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {

		return err
	}

	return nil
}

// AddVerificationPhoto is a function to add verification photo based on the given book entity
func (book *Book) AddVerificationPhoto(currentUser *database.MasterUser, referenceID uint, targetVerification entities.TransactionVerification) error {

	// add the new transaction verification photo to the database with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variable
		var dbErr error
		var newVerification database.DBTransactionVerification

		newVerification.ReferenceID = referenceID
		newVerification.PictDesc = targetVerification.PictDesc
		newVerification.URL = targetVerification.URL
		newVerification.IsActive = true
		newVerification.Created = time.Now().Local()
		newVerification.CreatedBy = currentUser.Username
		newVerification.Modified = time.Now().Local()
		newVerification.ModifiedBy = currentUser.Username

		// insert the new transaction verification photo to database
		if dbErr = tx.Create(&newVerification).Error; dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {

		return err
	}

	return nil
}
