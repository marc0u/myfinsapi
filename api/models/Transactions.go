package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID           uint32 `gorm:"primary_key; auto_increment" json:"id"`
	Date         string `gorm:"not null" json:"date"`
	Amount       int32  `gorm:"not null" json:"amount"`
	Type         string `gorm:"size:20; not null" json:"type"`
	DetailOrigin string `gorm:"size:50; not null" json:"detail_origin"`
	DetailCustom string `gorm:"size:50;" json:"detail_custom"`
	Category     string `gorm:"size:20;" json:"category"`
	Method       string `gorm:"size:20; not null" json:"method"`
	Bank         string `gorm:"size:20;" json:"bank"`
	Account      string `gorm:"size:20;" json:"account"`
	MadeBy       string `gorm:"size:20; not null" json:"made_by"`
	Balance      int32  `json:"balance"`
	// UserID			uint32		`sql:"type:int REFERENCES users(id)" json:"user_id"`
}

func (t *Transaction) Prepare() {
	t.ID = 0
	t.Type = html.EscapeString(strings.ToUpper(strings.TrimSpace(t.Type)))
	t.DetailOrigin = html.EscapeString(strings.TrimSpace(t.DetailOrigin))
	t.Category = html.EscapeString(strings.ToUpper(strings.TrimSpace(t.Category)))
	t.Method = html.EscapeString(strings.ToUpper(strings.TrimSpace(t.Method)))
	t.Bank = html.EscapeString(strings.ToUpper(strings.TrimSpace(t.Bank)))
	t.Account = html.EscapeString(strings.ToUpper(strings.TrimSpace(t.Account)))
	t.MadeBy = html.EscapeString(strings.ToUpper(strings.TrimSpace(t.MadeBy)))
	if t.DetailCustom == "" {
		t.DetailCustom = t.DetailOrigin
		return
	}
	t.DetailCustom = html.EscapeString(strings.TrimSpace(t.DetailCustom))
}

func (t *Transaction) Validate() error {
	if t.Date == "" {
		return errors.New("Date field is required.")
	}
	if t.Amount == 0 {
		return errors.New("Amount field is required.")
	}
	if t.Type == "" {
		return errors.New("Type field is required.")
	}
	if t.DetailOrigin == "" {
		return errors.New("Detail field is required.")
	}
	if t.Method == "" {
		return errors.New("Method field is required.")
	}
	if t.MadeBy == "" {
		return errors.New("Made By field is required.")
	}
	if len(t.Type) > 20 {
		return errors.New("Detail field must be under 50 characters.")
	}
	if len(t.DetailOrigin) > 50 {
		return errors.New("Detail field must be under 50 characters.")
	}
	if len(t.DetailCustom) > 50 {
		return errors.New("Detail field must be under 50 characters.")
	}
	if len(t.Category) > 20 {
		return errors.New("Detail field must be under 20 characters.")
	}
	if len(t.Method) > 20 {
		return errors.New("Method field must be under 20 characters.")
	}
	if len(t.Bank) > 20 {
		return errors.New("Bank field must be under 20 characters.")
	}
	if len(t.Account) > 20 {
		return errors.New("Account Number field must be under 20 characters.")
	}
	if len(t.MadeBy) > 20 {
		return errors.New("Made By field must be under 20 characters.")
	}
	return nil
}

func (t *Transaction) SaveTransaction(db *gorm.DB) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	return t, nil
}

func (t *Transaction) FindAllTransactions(db *gorm.DB) (*[]Transaction, error) {
	var err error
	transactions := []Transaction{}
	err = db.Debug().Model(&Transaction{}).Order("date desc").Order("id desc").Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	return &transactions, nil
}

func (t *Transaction) FindTransactionByID(db *gorm.DB, id uint64) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Where("id = ?", id).Take(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	return t, nil
}

func (t *Transaction) FindTransactionsBetweenDates(db *gorm.DB, from string, to string) (*[]Transaction, error) {
	transactions := []Transaction{}
	err := db.Debug().Model(&Transaction{}).Where("date BETWEEN ? AND ?", from, to).Order("date desc").Order("id desc").Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	return &transactions, nil
}

func (t *Transaction) FindLastTransaction(db *gorm.DB) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Last(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	return t, nil
}

func (t *Transaction) UpdateATransaction(db *gorm.DB, id uint64) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Where("id = ?", id).Updates(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	return t, nil
}

func (t *Transaction) DeleteATransaction(db *gorm.DB, id uint64) (int64, error) {
	db = db.Debug().Model(&Transaction{}).Where("id = ?", id).Take(&Transaction{}).Delete(&Transaction{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Transaction not found.")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
