package spaserver

// маршрутизация приложения
func (s *Server) Routes() error {
	s.loadViews()
	// s.protected.GET("/page", s.Page) // переход/загрузка на текущую страницу
	s.Echo().GET("/sse", s.Sse)
	return nil
}
