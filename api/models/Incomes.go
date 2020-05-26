package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Income struct {
	ID uint32 `gorm:"primary_key; auto_increment" json:"id"`
	// UserID			uint32		`sql:"type:int REFERENCES users(id)" json:"user_id"`
	DateTime      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_time"`
	Amount        int32     `gorm:"not null" json:"amount"`
	Detail        string    `gorm:"size:50; not null" json:"detail"`
	Category      string    `gorm:"size:20; not null" json:"category"`
	PaymentMethod string    `gorm:"size:20; not null" json:"payment_method"`
	MadeBy        string    `gorm:"size:20; not null" json:"made_by"`
}

func (i *Income) Prepare() {
	i.ID = 0
	i.DateTime = time.Now()
	i.Detail = html.EscapeString(strings.TrimSpace(i.Detail))
	i.Category = html.EscapeString(strings.ToUpper(strings.TrimSpace(i.Category)))
	i.PaymentMethod = html.EscapeString(strings.ToUpper(strings.TrimSpace(i.PaymentMethod)))
	i.MadeBy = html.EscapeString(strings.Title(strings.TrimSpace(i.MadeBy)))
}

func (i *Income) Validate() error {
	if i.Amount < 1 {
		return errors.New("Amount field is required.")
	}
	if i.Detail == "" {
		return errors.New("Detail field is required.")
	}
	if i.Category == "" {
		return errors.New("Category field is required.")
	}
	if i.PaymentMethod == "" {
		return errors.New("Payment Method field is required.")
	}
	if i.MadeBy == "" {
		return errors.New("Made By field is required.")
	}
	if len(i.Detail) > 50 {
		return errors.New("Detail field must be under 50 characters.")
	}
	if len(i.Category) > 20 {
		return errors.New("Detail field must be under 20 characters.")
	}
	if len(i.PaymentMethod) > 20 {
		return errors.New("Payment Method field must be under 20 characters.")
	}
	if len(i.MadeBy) > 20 {
		return errors.New("Made By field must be under 20 characters.")
	}
	return nil
}

func (i *Income) SaveIncome(db *gorm.DB) (*Income, error) {
	var err error
	err = db.Debug().Model(&Income{}).Create(&i).Error
	if err != nil {
		return &Income{}, err
	}
	return i, nil
}

func (i *Income) FindAllIncomes(db *gorm.DB) (*[]Income, error) {
	var err error
	incomes := []Income{}
	err = db.Debug().Model(&Income{}).Limit(100).Find(&incomes).Error
	if err != nil {
		return &[]Income{}, err
	}
	return &incomes, nil
}

func (i *Income) FindIncomeByID(db *gorm.DB, id uint64) (*Income, error) {
	var err error
	err = db.Debug().Model(&Income{}).Where("id = ?", id).Take(&i).Error
	if err != nil {
		return &Income{}, err
	}
	return i, nil
}

func (i *Income) UpdateAIncome(db *gorm.DB, id uint64) (*Income, error) {
	var err error
	err = db.Debug().Model(&Income{}).Where("id = ?", id).Updates(&i, true).Error
	if err != nil {
		return &Income{}, err
	}
	return i, nil
}

func (i *Income) DeleteAIncome(db *gorm.DB, id uint64) (int64, error) {
	db = db.Debug().Model(&Income{}).Where("id = ?", id).Take(&Income{}).Delete(&Income{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Income not found.")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
