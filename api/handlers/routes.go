package handlers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber"
)

func (s *Server) initializeRoutes() {
	version := os.Getenv("API_VERSION")[0:1]
	// Help
	s.Router.Get("/help", s.Help)
	// Handle Transactions
	s.Router.Post(fmt.Sprintf("/api/myfins/v%v/transactions", version), s.CreateTransaction)
	s.Router.Put(fmt.Sprintf("/api/myfins/v%v/transactions/:id", version), s.UpdateTransaction)
	s.Router.Delete(fmt.Sprintf("/api/myfins/v%v/transactions/:id", version), s.DeleteTransaction)
	// GET Transactions
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions", version), s.GetTransactions)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/last", version), s.GetLastTransaction)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/:id", version), s.GetTransactionByID)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/date/:year/:month", version), s.GetTransactionByDate)
	// Handle My Stocks
	s.Router.Post(fmt.Sprintf("/api/myfins/v%v/stocks", version), s.CreateStock)
	s.Router.Put(fmt.Sprintf("/api/myfins/v%v/stocks/:id", version), s.UpdateStock)
	s.Router.Delete(fmt.Sprintf("/api/myfins/v%v/stocks/:id", version), s.DeleteStock)
	// GET My Stocks
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks", version), s.GetStocks)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks/:id", version), s.GetStockByID)
}

func (server *Server) Help(c *fiber.Ctx) {
	version := os.Getenv("API_VERSION")[0:1]
	var msg = `Handle Transations
POST:/api/myfins/v%[1]v/transactions
PUT:/api/myfins/v%[1]v/transactions/:id
DELETE:/api/myfins/v%[1]v/transactions/:id

Get Transactions
/api/myfins/v%[1]v/transactions
/api/myfins/v%[1]v/transactions/last
/api/myfins/v%[1]v/transactions/:id
/api/myfins/v%[1]v/transactions/date/:year/:month

Handle Stocks
POST:/api/myfins/v%[1]v/stocks
PUT:/api/myfins/v%[1]v/stocks/:id
DELETE:/api/myfins/v%[1]v/stocks/:id

Get Stocks
GET:/api/myfins/v%[1]v/stocks
GET:/api/myfins/v%[1]v/stocks/:id`

	msg = fmt.Sprintf(msg, version)
	c.SendString(msg)
}
