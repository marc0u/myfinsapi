package seed

import (
	"log"

	"gitlab.com/marco.urriola/apifinances/api/models"

	"github.com/jinzhu/gorm"
)

var transactions = []models.Transaction{
	models.Transaction{
		Amount:        111,
		Type:          "income",
		DetailOrigin:  "First transaction April",
		Category:      "Test Category",
		Method:        "debit",
		Bank:          "scotiabank",
		AccountNumber: "000-000-000",
		MadeBy:        "marco",
	},
	models.Transaction{
		Amount:        222,
		Type:          "expense",
		DetailOrigin:  "Second transaction April",
		Category:      "Test Category",
		Method:        "debit",
		Bank:          "falabella",
		AccountNumber: "000-000-000",
		MadeBy:        "marco",
	},
	models.Transaction{
		Amount:        333,
		Type:          "income",
		DetailOrigin:  "Third transaction April",
		Category:      "Test Category",
		Method:        "transfer",
		Bank:          "scotiabank",
		AccountNumber: "000-000-000",
		MadeBy:        "marco",
	},
	models.Transaction{
		Amount:       444,
		Type:         "income",
		DetailOrigin: "Fourth transaction April",
		Category:     "Test Category",
		Method:       "cash",
		MadeBy:       "marco",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	for i, _ := range transactions {
		err = db.Debug().Model(&models.Transaction{}).Create(&transactions[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
