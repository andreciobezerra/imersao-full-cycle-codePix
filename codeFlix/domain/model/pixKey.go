package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type PixKey struct {
	Base      `valid:"required"`
	Kind      string   `json:"kind" valid:"notnull"`
	Key       string   `json:"key" valid:"notnull"`
	Account   *Account `valid:"-"`
	AccountID string   `valid:"-"`
	Status    string   `json:"status" valid:""notnull`
}

func (p *PixKey) isValid() error {
	_, err := govalidator.ValidateStruct(p)

	if p.Kind != "email" && p.Kind != "cpf" {
		return errors.New("Invalid type of key")
	}

	if p.Status != "active" && p.Status != "inactive" {
		return errors.New("Invalid status")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey{Kind: kind, Account: account, Key: key, AccountID: account.ID, Status: "active"}
	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.isValid()

	if err != nil {
		return err
	}

	return &pixKey, nil
}
