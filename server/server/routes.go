package server

import (
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/websocket/v2"
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
	api.Use("/polling/subscribe", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	api.Get("/polling/subscribe", handler.BuildPollerSubscriberWSHandler())
}

func registerFrontendRoutes(server *Server) {
	server.Static("/", "./frontend/build")
	server.Get("/_dashboard", monitor.New())
}
