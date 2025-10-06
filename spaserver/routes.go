package spaserver

// маршрутизация приложения
func (s *Server) Routes() error {
	s.loadViews()
	return nil
}
