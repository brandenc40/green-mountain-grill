package server

import "github.com/brandenc40/green-mountain-grill/internal/handler"

func RegisterRoutes(server *Server, handler *handler.Handler) {
	server.router.GET("/api/state", server.WithLogging(handler.GetGrillState))
	server.router.GET("/api/id", server.WithLogging(handler.GetGrillID))
	server.router.GET("/api/firmware", server.WithLogging(handler.GetGrillFirmware))
	server.router.GET("/api/session", server.WithLogging(handler.GetSessionData))
	server.router.POST("/api/session/new", server.WithLogging(handler.NewSession))
}
