package controllers

func (s *Server) initializeRoutes() {
	// Handle Transactions
	s.Router.Post("/api/myfins/v1/transactions", s.CreateTransaction)
	s.Router.Put("/api/myfins/v1/transactions/:id", s.UpdateTransaction)
	s.Router.Delete("/api/myfins/v1/transactions/:id", s.DeleteTransaction)
	// GET Transactions
	s.Router.Get("/api/myfins/v1/transactions", s.GetTransactions)
	s.Router.Get("/api/myfins/v1/transactions/last", s.GetLastTransaction)
	s.Router.Get("/api/myfins/v1/transactions/:id", s.GetTransactionByID)
	s.Router.Get("/api/myfins/v1/transactions/date/:year/:month", s.GetTransactionByDate)
}
