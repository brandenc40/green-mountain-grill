package server

import (
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func RegisterRoutes(server *Server, handler *handler.Handler) {
	registerAPIRoutes(server, handler)
	registerFrontendRoutes(server)
}

func registerAPIRoutes(server *Server, handler *handler.Handler) {
	api := server.Group("/api")
	api.Get("/state", handler.GetGrillState)
	api.Get("/id", handler.GetGrillID)
	api.Get("/firmware", handler.GetGrillFirmware)
	api.Get("/session", handler.GetSessionData)
	api.Get("/polling/start", handler.StartPolling)
	api.Get("/polling/stop", handler.StopPolling)
	api.Get("/polling/subscribers", handler.ViewSubscribers)
	api.Get("/polling/subscribe", handler.SubscribeToPoller)
}

func registerFrontendRoutes(server *Server) {
	server.Static("/", "./frontend/build")
	server.Get("/_dashboard", monitor.New())
}
