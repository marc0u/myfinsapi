package handlers

import (
	"strconv"

	"github.com/marc0u/myfinsapi/api/models"
	"github.com/marc0u/myfinsapi/api/utils"

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

func (server *Server) GetTransactionsLastMonth(c *fiber.Ctx) {
	from, to := utils.GetFirstLastDateCurrentMonth()
	// Getting data
	item := models.Transaction{}
	items, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetTransactionsBetweenDates(c *fiber.Ctx) {
	// Getting URL parameter ID
	from, to, err := utils.ParseFromToDates(c.Query("from"), c.Query("to"))
	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Getting data
	item := models.Transaction{}
	items, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetSummary(c *fiber.Ctx) {
	var (
		from string
		to   string
		err  error
	)
	if c.Query("from") == "" {
		from, to = utils.GetFirstLastDateCurrentMonth()
	} else {
		from, to, err = utils.ParseFromToDates(c.Query("from"), c.Query("to"))
		if err != nil {
			c.Status(400).JSON(fiber.Map{"error": err.Error()})
			return
		}
	}
	// Getting data
	item := models.Transaction{}
	itemsByDate, err := item.FindTransactionsBetweenDates(server.DB, from, to)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	// Processing general data
	summary := models.Summary{}
	summary.StartDate = from
	summary.EndDate = to
	summary.Incomes = models.ReduceAmountsByType(*itemsByDate, "INCOME")
	summary.Expenses = models.ReduceAmountsByType(*itemsByDate, "EXPENSE")
	// Processiong categories data
	categories, err := item.FindAllCategories(server.DB)
	if err != nil {
		c.Status(404).JSON(fiber.Map{"error": err.Error()})
		return
	}
	summary.CategoriesSummary = models.ReduceAmountsByCategories(*itemsByDate, categories)
	// Http response
	c.JSON(summary)
}
