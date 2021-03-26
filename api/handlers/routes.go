package handlers

import (
	"fmt"
)

var apiVersion string

func (s *Server) initializeRoutes(version string) {
	apiVersion = version
	version = version[0:1]
	// Notify
	s.Router.Get("/notify", s.Notify)
	// Login Logout
	s.Router.Post("/login", s.GoogleLogin)
	// Help
	s.Router.Get("/api/help", s.Help)
	s.Router.Get("/api/stack", s.GetStack)
	// Handle Transactions
	s.Router.Post(fmt.Sprintf("/api/myfins/v%v/transactions", version), s.CreateTransaction)
	s.Router.Put(fmt.Sprintf("/api/myfins/v%v/transactions/:id", version), s.UpdateTransaction)
	s.Router.Delete(fmt.Sprintf("/api/myfins/v%v/transactions/:id", version), s.DeleteTransaction)
	// GET Transactions
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions", version), s.GetTransactions)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/last", version), s.GetLastTransaction)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/month", version), s.GetTransactionsByMonth)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/dates", version), s.GetTransactionsBetweenDates)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/summary", version), s.GetSummaryByMonth)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/summary/dates", version), s.GetSummaryBetweenDates)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/transactions/:id", version), s.GetTransactionByID)
	// Handle My Stocks
	s.Router.Post(fmt.Sprintf("/api/myfins/v%v/stocks", version), s.CreateStock)
	s.Router.Put(fmt.Sprintf("/api/myfins/v%v/stocks/:id", version), s.UpdateStock)
	s.Router.Delete(fmt.Sprintf("/api/myfins/v%v/stocks/:id", version), s.DeleteStock)
	// GET My Stocks
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks", version), s.GetStocks)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks/holdings", version), s.GetHoldings)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks/summary", version), s.GetSummary)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks/portfolio/daily", version), s.GetPortfolioDaily)
	s.Router.Get(fmt.Sprintf("/api/myfins/v%v/stocks/:id", version), s.GetStockByID)
}
