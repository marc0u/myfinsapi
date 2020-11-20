package handlers

import (
	"fmt"
	"strconv"

	"github.com/marc0u/myfinsapi/api/models"

	"github.com/gofiber/fiber"
)

func (server *Server) CreateTransaction(c *fiber.Ctx) {
	// Reading body http request
	item := models.Transaction{}
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
	itemCreated, err := item.SaveTransaction(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.Status(201).JSON(itemCreated)
}

func (server *Server) GetTransactions(c *fiber.Ctx) {
	// Getting data
	item := models.Transaction{}
	items, err := item.FindAllTransactions(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetLastTransaction(c *fiber.Ctx) {
	// Getting data
	item := models.Transaction{}
	items, err := item.FindLastTransaction(server.DB)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetTransactionByID(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Getting data
	item := models.Transaction{}
	itemByID, err := item.FindTransactionByID(server.DB, id)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(itemByID)
}

func (server *Server) GetTransactionByDate(c *fiber.Ctx) {
	// Getting URL parameter ID
	year := c.Params("year")
	if year == "" {
		c.Status(400).JSON(fiber.Map{"error": "Year parameter is missing."})
		return
	}
	month := c.Params("month")
	if month == "" {
		c.Status(400).JSON(fiber.Map{"error": "Month parameter is missing"})
		return
	}
	date := fmt.Sprintf("%s-%s", year, month)
	// Getting data
	item := models.Transaction{}
	itemsByDate, err := item.FindTransactionsByDate(server.DB, date)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(itemsByDate)
}

func (server *Server) UpdateTransaction(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Reading body http request
	item := models.Transaction{}
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
	itemUpdated, err := item.UpdateATransaction(server.DB, id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	itemUpdated.ID = uint32(id)
	c.JSON(itemUpdated)
}

func (server *Server) DeleteTransaction(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Deleting data
	item := models.Transaction{}
	_, err = item.DeleteATransaction(server.DB, id)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.Status(204)
}
