package handlers

import (
	"strconv"
	"strings"

	"github.com/marc0u/myfinsapi/api/models"
	"github.com/marc0u/myfinsapi/api/utils"

	"github.com/gofiber/fiber/v2"
)

func (server *Server) CreateTransaction(c *fiber.Ctx) error {
	// Reading body http request
	item := models.Transaction{}
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
	itemCreated, err := item.SaveTransaction(server.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.Status(201).JSON(itemCreated)
}

func (server *Server) UpdateTransaction(c *fiber.Ctx) error {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Reading body http request
	item := models.Transaction{}
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
	itemUpdated, err := item.UpdateATransaction(server.DB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	itemUpdated.ID = uint32(id)
	return c.JSON(itemUpdated)
}

func (server *Server) DeleteTransaction(c *fiber.Ctx) error {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Deleting data
	item := models.Transaction{}
	_, err = item.DeleteATransaction(server.DB, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	c.Status(204)
	return nil
}

func (server *Server) GetTransactions(c *fiber.Ctx) error {
	// Getting data
	order := strings.Split(c.Query("order"), ",")
	desc := strings.Split(c.Query("desc"), ",")
	item := models.Transaction{}
	items, err := item.FindAllTransactions(server.DB, c.Query("limit"), order, desc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(items)
}

func (server *Server) GetLastTransaction(c *fiber.Ctx) error {
	// Getting data
	item := models.Transaction{}
	items, err := item.FindLastTransaction(server.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(items)
}

func (server *Server) GetTransactionByID(c *fiber.Ctx) error {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Getting data
	item := models.Transaction{}
	itemByID, err := item.FindTransactionByID(server.DB, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(itemByID)
}

func (server *Server) GetTransactionsByMonth(c *fiber.Ctx) error {
	// Getting URL parameters
	from, to, err := utils.GetFirstLastDateCurrentMonth(c.Query("change"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Getting data
	item := models.Transaction{}
	items, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(items)
}

func (server *Server) GetTransactionsBetweenDates(c *fiber.Ctx) error {
	// Getting URL parameters
	from, to, err := utils.ParseFromToDates(c.Query("from"), c.Query("to"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Getting data
	item := models.Transaction{}
	items, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Http response
	return c.JSON(items)
}

func (server *Server) GetSummaryByMonth(c *fiber.Ctx) error {
	// Getting URL parameters
	from, to, err := utils.GetFirstLastDateCurrentMonth(c.Query("change"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	// Getting data
	item := models.Transaction{}
	itemsByDate, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	categories, err := item.FindAllCategories(server.DB)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	exclusions := strings.Split(c.Query("exclusions"), ",")
	// Http response
	return c.JSON(models.ProcessSummary(from, to, *itemsByDate, categories, exclusions))
}

func (server *Server) GetSummaryBetweenDates(c *fiber.Ctx) error {
	// Getting URL parameters
	from, to, err := utils.ParseFromToDates(c.Query("from"), c.Query("to"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Getting data
	item := models.Transaction{}
	itemsByDate, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	categories, err := item.FindAllCategories(server.DB)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	exclusions := strings.Split(c.Query("exclusions"), ",")
	// Http response
	return c.JSON(models.ProcessSummary(from, to, *itemsByDate, categories, exclusions))
}
