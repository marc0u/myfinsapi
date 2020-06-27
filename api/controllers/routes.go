package controllers

func (s *Server) initializeRoutes() {
	// Transactions routes
	s.Router.Post("/api/myfins/v1/transactions", s.CreateTransaction)
	s.Router.Get("/api/myfins/v1/transactions", s.GetTransactions)
	s.Router.Get("/api/myfins/v1/transactions/last", s.GetLastTransaction)
	s.Router.Get("/api/myfins/v1/transactions/:id", s.GetTransaction)
	s.Router.Put("/api/myfins/v1/transactions/:id", s.UpdateTransaction)
	s.Router.Delete("/api/myfins/v1/transactions/:id", s.DeleteTransaction)
}
