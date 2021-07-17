package server

import (
	"github.com/brandenc40/green-mountain-grill/server/handler"
	"github.com/fasthttp/router"
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
		{fasthttp.MethodGet, "/api/polling/start", handler.StartPolling},
		{fasthttp.MethodGet, "/api/polling/stop", handler.StopPolling},
		{router.MethodWild, "/api/polling/subscribe", handler.SubscribeToPoller},
		{fasthttp.MethodGet, "/api/polling/subscribers", handler.ViewSubscribers},
	}
	for _, r := range routes {
		server.router.Handle(r.Method, r.Path, r.Handler)
	}
}
