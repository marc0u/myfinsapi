package controllers

func (s *Server) initializeRoutes() {
	// Incomes routes
	s.Router.Post("/api/finances/v1/incomes", s.CreateIncome)
	s.Router.Get("/api/finances/v1/incomes", s.GetIncomes)
	s.Router.Get("/api/finances/v1/incomes/:id", s.GetIncome)
	s.Router.Put("/api/finances/v1/incomes/:id", s.UpdateIncome)
	s.Router.Delete("/api/finances/v1/incomes/:id", s.DeleteIncome)
}
