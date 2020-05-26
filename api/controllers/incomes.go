package controllers

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	"strconv"

	"gitlab.com/marco.urriola/apifinances/api/models"
	// "gitlab.com/marco.urriola/apifinances/api/responses"
	// "gitlab.com/marco.urriola/apifinances/api/utils/formaterror"

	"github.com/gofiber/fiber"
	// "github.com/gorilla/mux"
)

func (server *Server) CreateIncome(c *fiber.Ctx) {
	// Reading body http request
	item := models.Income{}
	err := c.BodyParser(&item)
	if err != nil {
		c.Status(400)
		return
	}
	// Preparing and validating data
	item.Prepare()
	err = item.Validate()
	if err != nil {
		c.Status(400)
		return
	}
	// Saving data
	itemCreated, err := item.SaveIncome(server.DB)
	if err != nil {
		c.Status(500)
		return
	}
	// Http response
	c.Status(201).JSON(itemCreated)
}

func (server *Server) GetIncomes(c *fiber.Ctx) {
	// Getting data
	item := models.Income{}
	items, err := item.FindAllIncomes(server.DB)
	if err != nil {
		c.Status(500)
		return
	}
	// Http response
	c.JSON(items)
}

func (server *Server) GetIncome(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400)
		return
	}
	// Getting data
	item := models.Income{}
	itemByID, err := item.FindIncomeByID(server.DB, id)
	if err != nil {
		c.Status(404)
		return
	}
	// Http response
	c.JSON(itemByID)
}

func (server *Server) UpdateIncome(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(err)
		return
	}
	// Reading body http request
	item := models.Income{}
	err = c.BodyParser(&item)
	if err != nil {
		c.Status(400).JSON(err)
		return
	}
	// Preparing and validating data
	item.Prepare()
	err = item.Validate()
	if err != nil {
		c.Status(400).JSON(err)
		return
	}
	// Saving data
	itemUpdated, err := item.UpdateAIncome(server.DB, id)
	if err != nil {
		c.Status(500)
		return
	}
	// Http response
	itemUpdated.ID = uint32(id)
	c.JSON(itemUpdated)
}

func (server *Server) DeleteIncome(c *fiber.Ctx) {
	// Getting URL parameter ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400)
		return
	}
	// Deleting data
	item := models.Income{}
	_, err = item.DeleteAIncome(server.DB, id)
	if err != nil {
		c.Status(404)
		return
	}
	// Http response
	c.Status(204)
}
