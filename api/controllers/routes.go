package controllers

func (s *Server) initializeRoutes() {
	// Transactions routes
	s.Router.Post("/api/finances/v1/transactions", s.CreateTransaction)
	s.Router.Get("/api/finances/v1/transactions", s.GetTransactions)
	s.Router.Get("/api/finances/v1/transactions/:id", s.GetTransaction)
	s.Router.Put("/api/finances/v1/transactions/:id", s.UpdateTransaction)
	s.Router.Delete("/api/finances/v1/transactions/:id", s.DeleteTransaction)
}
