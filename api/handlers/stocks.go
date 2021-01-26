package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/marc0u/myfinsapi/api/models"
	"github.com/marc0u/myfinsapi/api/utils"

	"github.com/gofiber/fiber"

	"github.com/go-resty/resty"
)

func (server *Server) CreateStock(c *fiber.Ctx) {
	// Reading body http request
	item := models.Stock{}
	err := c.BodyParser(&item)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Preparing and validating data
	item.Prepare()
	err = item.Validate()
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Saving data
	itemCreated, err := item.SaveAStock(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.Status(201).JSON(itemCreated)
}

func (server *Server) UpdateStock(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Reading body http request
	item := models.Stock{}
	err = c.BodyParser(&item)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Preparing and validating data
	item.Prepare()
	err = item.Validate()
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Saving data
	itemUpdated, err := item.UpdateAStock(server.DB, id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	itemUpdated.ID = uint32(id)
	c.JSON(itemUpdated)
}

func (server *Server) DeleteStock(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Deleting data
	item := models.Stock{}
	_, err = item.DeleteAStock(server.DB, id)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.Status(204)
}

func (server *Server) GetStocks(c *fiber.Ctx) {
	// Getting data
	item := models.Stock{}
	items, err := item.FindAllStocks(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetStockByID(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Getting data
	item := models.Stock{}
	itemByID, err := item.FindStockByID(server.DB, id)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(itemByID)
}

func (server *Server) GetHoldings(c *fiber.Ctx) {
	// Getting data
	item := models.Stock{}
	tickers, err := item.FindTickers(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	items := []models.StockHolding{}
	for _, ticker := range tickers {
		result, err := item.FindStocksByTicker(server.DB, ticker)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err.Error()})
			return
		}
		holding := models.ReduceStockHolding(*result)
		if holding.StocksAmount != 0 {
			items = append(items, holding)
		}
	}
	// Http response
	c.JSON(items)
}

type StockBalance struct {
	Ticker       string
	StocksAmount int32
}

type Balance struct {
	Cash   int32
	Stocks []StockBalance
}

type DayBalance struct {
	Date   string
	Amount int32
}

type Prices struct {
	Date  string  `json:"Date"`
	Price float32 `json:"Close"`
}

type StockPrices struct {
	Ticker string
	Prices []Prices
}

func FetchDailyPrices(ticker string, result interface{}) (*resty.Response, error) {
	urlBase := fmt.Sprintf("http://rancher.loc:7002/api/stocks/v2/cl/day/%v", ticker)
	client := resty.New()
	resp, err := client.
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetRetryCount(3).
		SetRetryWaitTime(3 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetTimeout(10 * time.Second).
		R().
		SetResult(&result).
		ForceContentType("application/json").
		Get(urlBase)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (server *Server) GetPortfolioDaily(c *fiber.Ctx) {
	// Getting URL parameters
	from, to, err := utils.ParseFromToDates(c.Query("from"), c.Query("to"))
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Getting Portfolio Tickers
	item := models.Stock{}
	tickers, err := item.FindTickers(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Fetch Stocks Prices
	stocksPrices := []StockPrices{}
	for _, ticker := range tickers {
		resp, err := FetchDailyPrices(ticker, stocksPrices)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err.Error()})
			return
		}
		prices := []Prices{}
		err = json.Unmarshal(resp.Body(), &prices)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err.Error()})
			return
		}
		stock := StockPrices{ticker, prices}
		stocksPrices = append(stocksPrices, stock)
	}
	// Get Stocks records
	items, err := item.FindStocksBetweenDates(server.DB, from, to)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Data Processing
	daysBalance := []DayBalance{}
	currentDay := DayBalance{}
	stocksBalance := []StockBalance{}
	date := ""
	itemsLength := len(*items) - 1
	var balance int32
	for index, record := range *items {
		if date != record.Date || index == itemsLength {
			if date != "" {
				currentDay.Amount = currentDay.Amount + balance
			}
		}
		switch record.TransType {
		case "CREDIT", "DIVIDEND":
			balance = balance + int32(math.Abs(float64(record.TotalAmount)))
		case "WITHDRAWAL":
			balance = balance - int32(math.Abs(float64(record.TotalAmount)))
		case "BUY", "SELL":
			changed := false
			if record.TransType == "SELL" {
				balance = balance + int32(math.Abs(float64(record.TotalAmount)))
			} else {
				balance = balance - int32(math.Abs(float64(record.TotalAmount)))
			}
			for index, stock := range stocksBalance {
				if stock.Ticker == record.Ticker {
					if record.TransType == "SELL" {
						stocksBalance[index].StocksAmount = stocksBalance[index].StocksAmount - int32(math.Abs(float64(record.StocksAmount)))
					} else {
						stocksBalance[index].StocksAmount = stocksBalance[index].StocksAmount + int32(math.Abs(float64(record.StocksAmount)))
					}
					changed = true
					break
				}
			}
			if changed {
				break
			}
			stockBalance := StockBalance{Ticker: record.Ticker, StocksAmount: record.StocksAmount}
			stocksBalance = append(stocksBalance, stockBalance)
		}
		if date != record.Date || index == itemsLength {
			if date != "" {
				if index == itemsLength {
					currentDay.Amount = balance
				}
				for _, stock := range stocksBalance {
					if stock.StocksAmount < 0 {
						continue
					}
					for _, stockPrice := range stocksPrices {
						if stock.Ticker == stockPrice.Ticker {
							for _, price := range stockPrice.Prices {
								if price.Date == currentDay.Date {
									currentDay.Amount = currentDay.Amount + (int32(price.Price) * stock.StocksAmount)
									break
								}
							}
							break
						}
					}
				}
				daysBalance = append(daysBalance, currentDay)
				currentDay = DayBalance{}
			}
		}
		currentDay.Date = record.Date
		date = record.Date
	}
	// Http response
	c.JSON(daysBalance)
}
