package models

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

type StockPrices struct {
	Ticker string
	Prices []Price
}

type Price struct {
	Date  string  `json:"date"`
	Price float64 `json:"close"`
}

type StockBalance struct {
	Ticker       string
	StocksAmount int32
	StockPrice   float64
	TotalAmount  float64
}

type Balance struct {
	Date   string
	Cash   int32
	Stocks []StockBalance
}

type DayBalance struct {
	Date   string
	Amount int32
}

func (b *Balance) RemoveEmptyStocks() {
	for i := 0; i < len(b.Stocks); i++ {
		if b.Stocks[i].StocksAmount < 1 {
			b.Stocks = append(b.Stocks[:i], b.Stocks[i+1:]...)
			i = i - 1
		}
		if len(b.Stocks) == 1 {
			if b.Stocks[0].StocksAmount < 1 {
				b.Stocks = []StockBalance{}
			}
		}
	}
}

func (b *Balance) SetStocksPrices(stocksPrices []StockPrices) error {
	for i := 0; i < len(b.Stocks); i++ {
		for _, stockPrice := range stocksPrices {
			if b.Stocks[i].Ticker == stockPrice.Ticker {
				price, err := FindStockPricesByDate(b.Date, stockPrice.Prices)
				if err != nil {
					return err
				}
				b.Stocks[i].StockPrice = price.Price
				b.Stocks[i].TotalAmount = math.Floor(price.Price * float64(b.Stocks[i].StocksAmount))
				break
			}
		}
	}
	return nil
}

func (b *Balance) PrepareCash(record Stock) {
	b.Date = record.Date
	switch record.TransType {
	case "CREDIT", "DIVIDEND":
		b.Cash = b.Cash + int32(math.Abs(float64(record.TotalAmount)))
	case "WITHDRAWAL":
		b.Cash = b.Cash - int32(math.Abs(float64(record.TotalAmount)))
	case "BUY", "SELL":
		changed := false
		if record.TransType == "SELL" {
			b.Cash = b.Cash + int32(math.Abs(float64(record.TotalAmount)))
		} else {
			b.Cash = b.Cash - int32(math.Abs(float64(record.TotalAmount)))
		}
		for i, stock := range b.Stocks {
			if stock.Ticker == record.Ticker {
				if record.TransType == "SELL" {
					b.Stocks[i].StocksAmount = b.Stocks[i].StocksAmount - int32(math.Abs(float64(record.StocksAmount)))
				} else {
					b.Stocks[i].StocksAmount = b.Stocks[i].StocksAmount + int32(math.Abs(float64(record.StocksAmount)))
				}
				changed = true
				break
			}
		}
		if changed {
			break
		}
		stockBalance := StockBalance{Ticker: record.Ticker, StocksAmount: record.StocksAmount}
		b.Stocks = append(b.Stocks, stockBalance)
	}
}

func FetchStocksPrices(tickers []string) ([]StockPrices, error) {
	// Fetch Stocks Prices
	stocksPrices := []StockPrices{}
	for _, ticker := range tickers {
		prices, err := FetchDailyPrices(ticker)
		if err != nil {
			return nil, err
		}
		stock := StockPrices{ticker, prices}
		stocksPrices = append(stocksPrices, stock)
	}
	return stocksPrices, nil
}

func FetchDailyPrices(ticker string) ([]Price, error) {
	urlBase := fmt.Sprintf("http://192.168.1.15:7002/api/stocks/v2/cl/day/%v", ticker)
	client := resty.New()
	resp, err := client.
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetRetryCount(3).
		SetRetryWaitTime(3 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetTimeout(10 * time.Second).
		R().
		Get(urlBase)
	if err != nil {
		return nil, err
	}
	prices := []Price{}
	err = json.Unmarshal(resp.Body(), &prices)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func FindStockPricesByDate(date string, prices []Price) (Price, error) {
	for i := 0; i < 10; i++ {
		for _, price := range prices {
			if price.Date == date {
				return price, nil
			}
		}
		year, _ := strconv.Atoi(strings.Split(date, "-")[0])
		month, _ := strconv.Atoi(strings.Split(date, "-")[1])
		day, _ := strconv.Atoi(strings.Split(date, "-")[2])
		if day == 1 {
			if month == 1 {
				month = 12
				year = year - 1
			} else {
				month = month - 1
			}
			day = 31
		} else {
			day = day - 1
		}
		if month < 10 && day < 10 {
			date = fmt.Sprintf("%v-0%v-0%v", year, month, day)
			continue
		}
		if month < 10 {
			date = fmt.Sprintf("%v-0%v-%v", year, month, day)
			continue
		}
		if day < 10 {
			date = fmt.Sprintf("%v-%v-0%v", year, month, day)
			continue
		}
		date = fmt.Sprintf("%v-%v-%v", year, month, day)
	}
	return Price{}, errors.New("Stock price not found")
}
