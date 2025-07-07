package server

import (
	"ems/config"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
}

func New(echo *echo.Echo) *Server {
	return &Server{echo: echo}
}
func (s *Server) Start() {
	e := s.echo
	e.Logger.Fatal(e.Start(":" + config.App().Port))
}
