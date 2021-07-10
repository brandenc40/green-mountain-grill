package server

import (
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/valyala/fasthttp"
)

type route struct {
	Method  string
	Path    string
	Handler fasthttp.RequestHandler
}

func RegisterRoutes(server *Server, handler *handler.Handler) {
	routes := []route{
		{fasthttp.MethodGet, "/api/state", handler.GetGrillState},
		{fasthttp.MethodGet, "/api/id", handler.GetGrillID},
		{fasthttp.MethodGet, "/api/firmware", handler.GetGrillFirmware},
		{fasthttp.MethodGet, "/api/session", handler.GetSessionData},
		{fasthttp.MethodPost, "/api/session/new", handler.NewSession},
	}
	for _, r := range routes {
		server.router.Handle(r.Method, r.Path, r.Handler)
	}
}
