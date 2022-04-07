package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transactions []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	AccountFromID     string   `gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	Amount            float64  `json:"amount gorm:"type:float" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	PixKeyIdTo        string   `gorm:"column:pix_key_id_to;type:uuid;" valid:"notnull"`
	Status            string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description       string   `json:"description" gorm:"type:varchar(255)" valid:"-"`
	CancelDescription string   `json:"cancel_descrition" gorm:"type:varchar(255)" valid:"-"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("The amount must be greater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("Invalid status for the transaction")
	}

	if t.PixKeyTo.AccountID == t.AccountFromID {
		return errors.New("The source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}

	return nil

}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdateAt = time.Now()

	err := t.isValid()

	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.CancelDescription = description
	t.UpdateAt = time.Now()

	err := t.isValid()

	return err
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string, id string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom:   accountFrom,
		AccountFromID: accountFrom.ID,
		Amount:        amount,
		PixKeyTo:      pixKeyTo,
		PixKeyIdTo:    pixKeyTo.ID,
		Status:        TransactionPending,
		Description:   description,
	}

	if id == "" {
		transaction.ID = uuid.NewV4().String()
	} else {
		transaction.ID = id
	}
	transaction.CreatedAt = time.Now()
	err := transaction.isValid()
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
