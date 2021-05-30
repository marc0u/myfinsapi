package handlers

import (
	"strconv"
	"strings"

	"github.com/marc0u/myfinsapi/api/models"

	"github.com/gofiber/fiber/v2"
)

func (server *Server) GetStack(c *fiber.Ctx) error {
	return c.JSON(server.Router.Stack())
}

func (server *Server) CreateStock(c *fiber.Ctx) error {
	// Reading body http request
	item := models.Stock{}
	err := c.BodyParser(&item)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Preparing and validating data
	item.Prepare()
	err = item.Validate()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Saving data
	itemCreated, err := item.SaveAStock(server.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.Status(201).JSON(itemCreated)
}

func (server *Server) UpdateStock(c *fiber.Ctx) error {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Reading body http request
	item := models.Stock{}
	err = c.BodyParser(&item)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Preparing and validating data
	item.Prepare()
	err = item.Validate()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Saving data
	itemUpdated, err := item.UpdateAStock(server.DB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	itemUpdated.ID = uint32(id)
	return c.JSON(itemUpdated)
}

func (server *Server) DeleteStock(c *fiber.Ctx) error {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Deleting data
	item := models.Stock{}
	_, err = item.DeleteAStock(server.DB, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	c.Status(204)
	return nil
}

func (server *Server) GetStocks(c *fiber.Ctx) error {
	// Getting data
	item := models.Stock{}
	items, err := item.FindAllStocks(server.DB, c.Query("desc"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(items)
}

func (server *Server) GetStockByID(c *fiber.Ctx) error {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Getting data
	item := models.Stock{}
	itemByID, err := item.FindStockByID(server.DB, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(itemByID)
}

func (server *Server) GetHoldings(c *fiber.Ctx) error {
	// Getting data
	item := models.Stock{}
	holdings, err := item.FindHoldings(server.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(holdings)
}

func (server *Server) GetSummary(c *fiber.Ctx) error {
	item := models.Stock{}
	credit, err := item.FindStocksByTransType(server.DB, "CREDIT")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	withdrawal, err := item.FindStocksByTransType(server.DB, "WITHDRAWAL")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	holdings, err := item.FindHoldings(server.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	totalAssets := models.ReduceHoldings(holdings)
	totalInvested := models.ReduceTotalAmount(*credit) - models.ReduceTotalAmount(*withdrawal)
	totalGained := totalAssets - totalInvested
	return c.JSON(fiber.Map{"total_assets": totalAssets, "total_invested": totalInvested, "total_gained": totalGained})
}

func (server *Server) GetPortfolioDaily(c *fiber.Ctx) error {
	// Getting Portfolio Tickers
	item := models.Stock{}
	tickers, err := item.FindTickers(server.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Fetch Stocks Prices
	stocksPrices, err := models.FetchStocksPrices(tickers)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Get Stocks records
	items, err := item.FindAllStocks(server.DB, "false")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Process Detailed Balance
	balance := models.MakeBalance(*items)
	balance, err = models.FillMissedDays(balance)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	balance, err = models.SetStocksPrices(balance, stocksPrices)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response Balance Detailed
	if strings.ToLower(c.Query("detailed")) == "true" {
		return c.JSON(balance)
	}
	// Process Compact Balance
	compactBalance, err := models.CompactBalance(balance)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(compactBalance)
}
