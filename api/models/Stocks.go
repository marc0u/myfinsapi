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
	Date              string  `gorm:"not null" json:"date"`
	Ticker            string  `gorm:"size:10" json:"ticker"`
	StocksAmount      int32   `json:"stocks_amount"`
	StockPrice        float64 `json:"stock_price"`
	TotalAmount       float64 `gorm:"not null" json:"current_total_amount"`
	BoughtTotalAmount float64 `gorm:"not null" json:"bought_total_amount"`
	Country           string  `gorm:"not null; size:20" json:"country"`
	Currency          string  `gorm:"not null; size:10" json:"currency"`
}

func (r *Stock) Prepare() {
	r.ID = 0
	r.TransType = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.TransType)))
	r.Country = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Country)))
	r.Currency = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Currency)))
	r.TotalAmount = math.Round(r.TotalAmount*100) / 100
	if r.TransType == "SELL" || r.TransType == "BUY" {
		r.Ticker = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Ticker)))
		if r.TransType == "SELL" {
			r.StocksAmount = int32(math.Abs(float64(r.StocksAmount))) * -1
			r.TotalAmount = math.Abs(r.TotalAmount)
		}
		if r.TransType == "BUY" {
			r.StocksAmount = int32(math.Abs(float64(r.StocksAmount)))
			r.TotalAmount = math.Abs(r.TotalAmount) * -1
		}
		r.StockPrice = math.Round(math.Abs(r.TotalAmount/float64(r.StocksAmount))*100) / 100
		return
	}
	if r.TransType == "DIVIDEND" {
		r.Ticker = html.EscapeString(strings.ToUpper(strings.TrimSpace(r.Ticker)))
	} else {
		r.Ticker = ""
	}
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
	lastItem := Stock{}
	err = db.Model(&Stock{}).Order("date desc").Order("id desc").Last(&lastItem).Error
	if err != nil && err.Error() != "record not found" {
		return &Stock{}, err
	}
	r.SetBalance(lastItem)
	err = db.Model(&Stock{}).Create(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) SetBalance(lastItem Stock) {
	switch r.TransType {
	case "CREDIT", "DIVIDEND", "SELL":
		r.Balance = lastItem.Balance + math.Abs(r.TotalAmount)
	case "WITHDRAWAL", "BUY":
		r.Balance = lastItem.Balance - math.Abs(r.TotalAmount)
	}
}

func (r *Stock) UpdateAStock(db *gorm.DB, id uint64) (*Stock, error) {
	var err error
	err = db.Model(&Stock{}).Where("id = ?", id).Updates(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) DeleteAStock(db *gorm.DB, id uint64) (int64, error) {
	db = db.Model(&Stock{}).Where("id = ?", id).Take(&Stock{}).Delete(&Stock{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Stock not found.")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (r *Stock) FindAllStocks(db *gorm.DB, desc string) (*[]Stock, error) {
	var err error
	stocks := []Stock{}
	if desc == "false" {
		err = db.Model(&Stock{}).Order("date").Order("id").Find(&stocks).Error
	} else {
		err = db.Model(&Stock{}).Order("date desc").Order("id desc").Find(&stocks).Error
	}
	if err != nil {
		return &[]Stock{}, err
	}
	return &stocks, nil
}

func (r *Stock) FindStockByID(db *gorm.DB, id uint64) (*Stock, error) {
	var err error
	err = db.Model(&Stock{}).Where("id = ?", id).Take(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) FindStocksByTicker(db *gorm.DB, ticker string) (*[]Stock, error) {
	var err error
	stocks := []Stock{}
	err = db.Model(&Stock{}).Where("ticker = ?", ticker).Find(&stocks).Error
	if err != nil {
		return &[]Stock{}, err
	}
	return &stocks, nil
}

func (r *Stock) FindStocksByTransType(db *gorm.DB, trans_type string) (*[]Stock, error) {
	var err error
	stocks := []Stock{}
	err = db.Model(&Stock{}).Where("trans_type = ?", trans_type).Find(&stocks).Error
	if err != nil {
		return &[]Stock{}, err
	}
	return &stocks, nil
}

func (r *Stock) FindTickers(db *gorm.DB) ([]string, error) {
	var err error
	stocks := []Stock{}
	err = db.Model(&Stock{}).Select("ticker").Not("ticker = ?", "").Find(&stocks).Error
	if err != nil {
		return []string{}, err
	}
	tickers := []string{}
	for _, stock := range stocks {
		tickers = append(tickers, stock.Ticker)
	}
	return utils.RemoveDuplicateStrings(tickers), nil
}

func (r *Stock) FindStocksBetweenDates(db *gorm.DB, from string, to string) (*[]Stock, error) {
	stocks := []Stock{}
	err := db.Model(&Stock{}).Where("date BETWEEN ? AND ?", from, to).Order("date").Order("id").Find(&stocks).Error
	if err != nil {
		return &[]Stock{}, err
	}
	return &stocks, nil
}

func (r *Stock) FindLastRecord(db *gorm.DB) (*Stock, error) {
	var err error
	err = db.Model(&Transaction{}).Order("date desc").Order("id desc").Last(&r).Error
	if err != nil {
		return &Stock{}, err
	}
	return r, nil
}

func (r *Stock) FindHoldings(db *gorm.DB) ([]StockHolding, error) {
	tickers, err := r.FindTickers(db)
	if err != nil {
		return nil, err
	}
	balance, err := r.FindLastRecord(db)
	if err != nil {
		return nil, err
	}
	stocksHolding := []StockHolding{}
	stocksHolding = append(stocksHolding, StockHolding{
		Date:        balance.Date,
		Ticker:      "CASH",
		TotalAmount: balance.Balance,
		Country:     balance.Country,
		Currency:    balance.Currency,
	})
	for _, ticker := range tickers {
		result, err := r.FindStocksByTicker(db, ticker)
		if err != nil {
			return nil, err
		}
		holding := ReduceStocksAmount(*result)
		if holding.StocksAmount > 0 {
			prices, err := FetchDailyPrices(ticker)
			if err != nil {
				return nil, err
			}
			lastPrice := prices[len(prices)-1]
			holding.Date = lastPrice.Date
			holding.StockPrice = lastPrice.Price
			holding.TotalAmount = math.Round(float64(holding.StocksAmount)*holding.StockPrice*100) / 100
			stocksHolding = append(stocksHolding, holding)
		}
	}
	return stocksHolding, nil
}

func ReduceStocksAmount(stocks []Stock) StockHolding {
	var stocksAmount float64
	var totalAmount float64
	for _, stock := range stocks {
		if stock.TransType == "SELL" {
			stocksAmount = stocksAmount - math.Abs(float64(stock.StocksAmount))
			totalAmount = totalAmount - math.Abs(float64(stock.TotalAmount))
		}
		if stock.TransType == "BUY" {
			stocksAmount = stocksAmount + math.Abs(float64(stock.StocksAmount))
			totalAmount = totalAmount + math.Abs(float64(stock.TotalAmount))
		}
		if stock.StocksAmount < 1 {
			totalAmount = 0
			continue
		}
	}
	return StockHolding{
		Ticker:            stocks[0].Ticker,
		StocksAmount:      int32(stocksAmount),
		BoughtTotalAmount: totalAmount,
		Country:           stocks[0].Country,
		Currency:          stocks[0].Currency}
}

func ReduceTotalAmount(stocks []Stock) float64 {
	var totalAmount float64
	for _, stock := range stocks {
		totalAmount += stock.TotalAmount
	}
	return totalAmount
}

func ReduceHoldings(holdings []StockHolding) float64 {
	var totalHoldings float64
	for _, holding := range holdings {
		totalHoldings += holding.TotalAmount
	}
	return totalHoldings
}
