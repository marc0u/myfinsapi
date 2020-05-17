package seed

import (
	"log"

	"gitlab.com/marco.urriola/apifinances/api/models"

	"github.com/jinzhu/gorm"
)

var incomes = []models.Income{
	models.Income{
		Amount:        111,
		Detail:        "First income April",
		Category:      "Test Category",
		PaymentMethod: "sc",
		MadeBy:        "marco",
	},
	models.Income{
		Amount:        222,
		Detail:        "Second income April",
		Category:      "Test Category",
		PaymentMethod: "sc",
		MadeBy:        "marco",
	},
	models.Income{
		Amount:        111,
		Detail:        "First income April",
		Category:      "Test Category",
		PaymentMethod: "sc",
		MadeBy:        "marco",
	},
	models.Income{
		Amount:        222,
		Detail:        "Second income April",
		Category:      "Test Category",
		PaymentMethod: "sc",
		MadeBy:        "marco",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&models.Income{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	// for i, _ := range incomes {
	// 	err = db.Debug().Model(&models.Income{}).Create(&incomes[i]).Error
	// 	if err != nil {
	// 		log.Fatalf("cannot seed posts table: %v", err)
	// 	}
	// }
}
