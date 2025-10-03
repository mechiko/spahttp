package spaserver

import (
	"spahttp/spaserver/views/footer"
	"spahttp/spaserver/views/header"
	"spahttp/spaserver/views/home"
	"spahttp/spaserver/views/login"
)

// загружаем все виды
func (s *Server) loadViews() {
	// view header
	view1 := header.New(s)
	s.views[view1.Model()] = view1
	view1.InitData()
	view1.Routes()
	// view footer
	view2 := footer.New(s)
	s.views[view2.Model()] = view2
	view2.Routes()
	view2.InitData()

	view3 := login.New(s)
	s.views[view3.Model()] = view3
	view3.Routes()
	view3.InitData()

	view4 := home.New(s, s.protected)
	s.views[view4.Model()] = view4
	view4.Routes()
	view4.InitData()
}
