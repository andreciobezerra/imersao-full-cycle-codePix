package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Bank struct {
	Base     `valid:"required"`
	Name     string     `json:"name" valid:"notnull"`
	Code     string     `json:"code" valid:"notnull"`
	Accounts []*Account `valid:"-"`
}

func (bank *Bank) isValid() error {
	_, err := govalidator.ValidateStruct(bank)

	if err != nil {
		return err
	}

	return nil
}

func NewBank(code, name string) (*Bank, error) {
	bank := Bank{Code: code, Name: name}
	bank.ID = uuid.NewV4().String()
	bank.CreatedAt = time.Now()

	err := bank.isValid()

	if err != nil {
		return nil, err
	}

	return &bank, nil
}
