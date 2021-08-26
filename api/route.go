package api

func (s *Server) route() {
	v1 := s.engine.Group("/v1")
	processor := v1.Group("/processor")
	{
		processor.POST("/file", s.mustToken, s.processFile)
		processor.GET("/workers", s.mustToken, s.workers)
		processor.GET("/tasks/pending", s.mustToken, s.pendingList)
		processor.GET("/tasks/failed", s.mustToken, s.failedList)

		processor.GET("/info", s.mustToken, s.supportClassifies)
	}

	v1.GET("/version", s.version)
	s.engine.GET("/", s.health)
}
