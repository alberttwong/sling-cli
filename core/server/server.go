package server

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cast"
)

func NewServer(host, port string) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return cast.ToString(time.Now().UnixNano())
		},
	}))

	wsServer := NewWsServer()
	e.GET(RouteWs.String(), wsServer.NewClient)
	addRoutes(e)

	if port == "" {
		port = "9876"
	}

	return &Server{
		Host:       host,
		Port:       port,
		WsServer:   wsServer,
		EchoServer: e,
	}
}

type Server struct {
	Host       string
	Port       string
	EchoServer *echo.Echo
	WsServer   *WsServer
}

func (s *Server) Start() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	s.EchoServer.Logger.Fatal(s.EchoServer.Start(s.Host + ":" + s.Port))
}
