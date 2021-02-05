package handlers

import (
	"crypto/tls"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/marc0u/myfinsapi/api/models"

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
	items, err := item.FindAllStocks(server.DB, c.Query("desc"))
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
	balance, err := item.FindLastRecord(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	items := []models.StockHolding{}
	items = append(items, models.StockHolding{
		Date:        balance.Date,
		Ticker:      "CASH",
		TotalAmount: balance.Balance,
		Country:     balance.Country,
		Currency:    balance.Currency,
	})
	for _, ticker := range tickers {
		result, err := item.FindStocksByTicker(server.DB, ticker)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err.Error()})
			return
		}
		holding := models.ReduceStocksAmount(*result)
		if holding.StocksAmount > 0 {
			prices, err := models.FetchDailyPrices(ticker)
			if err != nil {
				c.Status(500).JSON(fiber.Map{"error": err.Error()})
				return
			}
			lastPrice := prices[len(prices)-1]
			holding.Date = lastPrice.Date
			holding.StockPrice = lastPrice.Price
			holding.TotalAmount = math.Round(float64(holding.StocksAmount)*holding.StockPrice*100) / 100
			items = append(items, holding)
		}
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetPortfolioDaily(c *fiber.Ctx) {
	// Getting Portfolio Tickers
	item := models.Stock{}
	tickers, err := item.FindTickers(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Fetch Stocks Prices
	stocksPrices, err := models.FetchStocksPrices(tickers)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Get Stocks records
	items, err := item.FindAllStocks(server.DB, "false")
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Process Detailed Balance
	balance := models.MakeBalance(*items)
	balance, err = models.FillMissedDays(balance)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	balance, err = models.SetPrices(balance, stocksPrices)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response Balance Detailed
	if strings.ToLower(c.Query("detail")) == "detailed" {
		c.JSON(balance)
		return
	}
}

func (server *Server) MirrorProductionTables() error {
	trans := []models.Transaction{}
	stocks := []models.Stock{}
	urlTrans := "http://192.168.1.15:7001/api/myfins/v2/transactions"
	urlStocks := "http://192.168.1.15:7001/api/myfins/v2/stocks"
	client := resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetRetryCount(3).SetRetryWaitTime(3 * time.Second).SetRetryMaxWaitTime(5 * time.Second).SetTimeout(10 * time.Second)
	resp, err := client.R().Get(urlTrans)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &trans)
	if err != nil {
		return err
	}
	resp, err = client.R().Get(urlStocks)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &stocks)
	if err != nil {
		return err
	}
	for i, j := 0, len(trans)-1; i < j; i, j = i+1, j-1 {
		trans[i], trans[j] = trans[j], trans[i]
	}
	for _, item := range trans {
		// Saving data
		_, err := item.SaveTransaction(server.DB)
		if err != nil {
			return err
		}
	}
	for i, j := 0, len(stocks)-1; i < j; i, j = i+1, j-1 {
		stocks[i], stocks[j] = stocks[j], stocks[i]
	}
	for _, item := range stocks {
		// Saving data
		_, err := item.SaveAStock(server.DB)
		if err != nil {
			return err
		}
	}
	return nil
}
