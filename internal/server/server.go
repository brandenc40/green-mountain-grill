package server

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"go.uber.org/fx"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var Module = fx.Provide(New, NewConfig)

type Params struct {
	fx.In

	Config    *Config
	Logger    *logrus.Logger
	Lifecycle fx.Lifecycle
}

func New(p Params) *Server {
	r := router.New()
	s := &Server{
		config: p.Config,
		logger: p.Logger,
		router: r,
		httpServer: &fasthttp.Server{
			Handler:     r.Handler,
			ReadTimeout: 3 * time.Second,
		},
	}
	p.Lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					err := s.Run()
					if err != nil {
						p.Logger.WithError(err).Error("unable to start server")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown()
			},
		},
	)
	return s
}

type Server struct {
	config     *Config
	logger     *logrus.Logger
	httpServer *fasthttp.Server
	router     *router.Router
}

func (s *Server) Run() error {
	s.logger.Info("Running on port " + s.config.ServerPort)
	return s.httpServer.ListenAndServe(s.config.ServerPort)
}

func (s *Server) Shutdown() error {
	s.logger.Info("shutting down server")
	return s.httpServer.Shutdown()
}

func (s *Server) WithLogging(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var b bytes.Buffer
		// log request
		if s.logger.Level == logrus.DebugLevel {
			b.WriteString("[REQUEST ")
			b.Write(ctx.Method())
			b.WriteString(" ")
			b.Write(ctx.Path())
			b.WriteString("]")
			s.logger.Debug(b.String())
			b.Reset()
		}

		// start timer and handle
		start := time.Now()
		h(ctx)
		dur := time.Since(start).String()

		// log response
		if s.logger.Level == logrus.DebugLevel {
			b.WriteString("[RESPONSE ")
			b.WriteString(strconv.Itoa(ctx.Response.StatusCode()))
			if ctx.Response.StatusCode() != fasthttp.StatusOK {
				b.WriteString(" (")
				b.Write(ctx.Response.Body())
				b.WriteString(")")
			}
			b.WriteString(" ")
			b.WriteString(dur)
			b.WriteString("]")
			s.logger.Debug(b.String())
		}
	}
}
