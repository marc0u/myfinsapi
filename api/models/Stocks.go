package models

import (
	"errors"
	"html"
	"math"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/marc0u/myfinsapi/api/utils"
)

type Stock struct {
	ID           uint32  `gorm:"primary_key; auto_increment" json:"id"`
	Date         string  `gorm:"not null" json:"date"`
	Ticker       string  `gorm:"size:10" json:"ticker"`
	TransType    string  `gorm:"not null; size:20" json:"trans_type"`
	StocksAmount int32   `json:"stocks_amount"`
	StockPrice   float64 `json:"stock_price"`
	TotalAmount  float64 `gorm:"not null" json:"total_amount"`
	Balance      float64 `gorm:"not null" json:"balance"`
	Country      string  `gorm:"not null; size:20" json:"country"`
	Currency     string  `gorm:"not null; size:10" json:"currency"`
}

type StockHolding struct {
	Ticker       string  `gorm:"size:10" json:"ticker"`
	StocksAmount int32   `json:"stocks_amount"`
	StockPrice   float64 `json:"stock_price"`
	TotalAmount  float64 `gorm:"not null" json:"total_amount"`
	Country      string  `gorm:"not null; size:20" json:"country"`
	Currency     string  `gorm:"not null; size:10" json:"currency"`
}

func (r *Stock) Prepare() {
	r.ID = 0
	r.Ticker = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Ticker)))
	r.TransType = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.TransType)))
	r.Country = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Country)))
	r.Currency = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Currency)))
	if r.TransType == "SELL" || r.TransType == "BUY" {
		if r.TransType == "SELL" {
			r.StocksAmount = int32(math.Abs(float64(r.StocksAmount))) * -1
		}
		if r.TransType == "BUY" {
			r.StocksAmount = int32(math.Abs(float64(r.StocksAmount)))
		}
		r.StockPrice = math.Abs(r.TotalAmount / float64(r.StocksAmount))
		return
	}
	r.Ticker = ""
	r.StocksAmount = 0
	r.StockPrice = 0
}

func (r *Stock) Validate() error {
	if r.Date == "" {
		return errors.New("Date field is required.")
	}
	if r.TransType == "" {
		return errors.New("TransType field is required.")
	}
	if r.TotalAmount == 0 {
		return errors.New("TotalAmount field must not be 0.")
	}
	if r.Country == "" {
		return errors.New("Country field is required.")
	}
	if r.Currency == "" {
		return errors.New("Currency field is required.")
	}
	if len(r.TransType) > 20 {
		return errors.New("TransType field must be under 20 characters.")
	}
	if len(r.Country) > 20 {
		return errors.New("Country field must be under 20 characters.")
	}
	if len(r.Currency) > 10 {
		return errors.New("Currency field must be under 10 characters.")
	}
	if r.TransType == "SELL" || r.TransType == "BUY" {
		if r.Ticker == "" {
			return errors.New("Ticker field is required.")
		}
		if r.StocksAmount == 0 {
			return errors.New("StocksAmount field must not be 0.")
		}
		if len(r.Ticker) > 10 {
			return errors.New("Ticker field must be under 10 characters.")
		}
	}
	return nil
}

func (r *Stock) SaveAStock(db *gorm.DB) (*Stock, error) {
	var err error
	item := Stock{}
	err = db.Debug().Model(&Stock{}).Last(&item).Error
	if r.TransType == "BUY" || r.TransType == "WITHDRAWAL" {
		r.Balance = item.Balance - r.TotalAmount
	} else {
		r.Balance = item.Balance + r.TotalAmount
	}
	err = db.Debug().Model(&Stock{}).Create(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) UpdateAStock(db *gorm.DB, id uint64) (*Stock, error) {
	var err error
	err = db.Debug().Model(&Stock{}).Where("id = ?", id).Updates(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) DeleteAStock(db *gorm.DB, id uint64) (int64, error) {
	db = db.Debug().Model(&Stock{}).Where("id = ?", id).Take(&Stock{}).Delete(&Stock{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Stock not found.")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (r *Stock) FindAllStocks(db *gorm.DB) (*[]Stock, error) {
	var err error
	stocks := []Stock{}
	err = db.Debug().Model(&Stock{}).Order("date desc").Order("id desc").Find(&stocks).Error
	if err != nil {
		return &[]Stock{}, err
	}
	return &stocks, nil
}

func (r *Stock) FindStockByID(db *gorm.DB, id uint64) (*Stock, error) {
	var err error
	err = db.Debug().Model(&Stock{}).Where("id = ?", id).Take(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) FindStocksByTicker(db *gorm.DB, ticker string) (*[]Stock, error) {
	var err error
	stocks := []Stock{}
	err = db.Debug().Model(&Stock{}).Where("ticker = ?", ticker).Find(&stocks).Error
	if err != nil {
		return &[]Stock{}, err
	}
	return &stocks, nil
}

func (r *Stock) FindTickers(db *gorm.DB) ([]string, error) {
	var err error
	stocks := []Stock{}
	err = db.Debug().Model(&Stock{}).Select("ticker").Not("ticker = ?", "").Find(&stocks).Error
	if err != nil {
		return []string{}, err
	}
	tickers := []string{}
	for _, stock := range stocks {
		tickers = append(tickers, stock.Ticker)
	}
	return utils.RemoveDuplicateStrings(tickers), nil
}

func ReduceStockHolding(stocks []Stock) StockHolding {
	var (
		stocksAmount int32
		totalAmount  float64
		stockPrice   float64
	)
	for _, stock := range stocks {
		stocksAmount = stocksAmount + stock.StocksAmount
		if stocksAmount == 0 {
			totalAmount = 0.0
			continue
		}
		totalAmount = totalAmount + stock.TotalAmount
	}
	if stocksAmount != 0 {
		stockPrice = totalAmount / float64(stocksAmount)
	}
	return StockHolding{Ticker: stocks[0].Ticker,
		StocksAmount: stocksAmount,
		StockPrice:   stockPrice,
		TotalAmount:  totalAmount,
		Country:      stocks[0].Country,
		Currency:     stocks[0].Currency}
}

// func (r *Stock) FindStocksHolings(db *gorm.DB) (*[]Stock, error) {
// 	var err error
// 	stocks := []Stock{}
// 	err = db.Debug().Model(&Stock{}).Order("date").Not("ticker = ?", "").Find(&stocks).Error
// 	if err != nil {
// 		return &[]Stock{}, err
// 	}
// 	return &stocks, nil
// }
