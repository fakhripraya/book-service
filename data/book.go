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
		var updateTransaction database.DBTransaction
		var dbErr error

		updateTransaction.PaidOff = targetTransaction.PaidOff
		updateTransaction.MustPay = targetTransaction.MustPay
		updateTransaction.IsActive = targetTransaction.IsActive
		updateTransaction.Modified = time.Now().Local()
		updateTransaction.ModifiedBy = currentUser.Username

		// update the transaction
		dbErr = config.DB.Save(updateTransaction).Error

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
		var updateTransactionDetail database.DBTransactionDetail
		var dbErr error

		updateTransactionDetail.PaymentMethodID = targetTransactionDetail.PaymentMethodID
		updateTransactionDetail.Status = targetTransactionDetail.Status
		updateTransactionDetail.Payment = targetTransactionDetail.Payment
		updateTransactionDetail.IsActive = targetTransactionDetail.IsActive
		updateTransactionDetail.Modified = time.Now().Local()
		updateTransactionDetail.ModifiedBy = currentUser.Username

		// update the transaction detail
		dbErr = config.DB.Save(updateTransactionDetail).Error

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
func (book *Book) AddRoomBookMember(currentUser *database.MasterUser, roomBookID uint, targetRoomBookMember *entities.TransactionRoomBookMember) error {

	// add the room book member to the database with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var newBookMember database.DBTransactionRoomBookMember
		var dbErr error

		newBookMember.RoomBookID = roomBookID
		newBookMember.TenantID = targetRoomBookMember.TenantID
		newBookMember.IsActive = true
		newBookMember.Created = time.Now().Local()
		newBookMember.CreatedBy = currentUser.Username
		newBookMember.Modified = time.Now().Local()
		newBookMember.ModifiedBy = currentUser.Username

		// insert the new room book member to database
		if dbErr = tx.Create(&newBookMember).Error; dbErr != nil {
			return dbErr
		}

		// add the room book member details to the database with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var newRoomBookMemberDetail = targetRoomBookMember.Members

			// add the room book member id to the slices
			for i := range newRoomBookMemberDetail {
				(&newRoomBookMemberDetail[i]).RoomBookMemberID = newBookMember.ID
				(&newRoomBookMemberDetail[i]).IsActive = true
				(&newRoomBookMemberDetail[i]).Created = time.Now().Local()
				(&newRoomBookMemberDetail[i]).CreatedBy = currentUser.Username
				(&newRoomBookMemberDetail[i]).Modified = time.Now().Local()
				(&newRoomBookMemberDetail[i]).ModifiedBy = currentUser.Username
			}

			// insert the new room book member details to database
			if dbErr2 = tx2.Create(&newRoomBookMemberDetail).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error
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
