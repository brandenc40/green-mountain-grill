package server

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Config     *Config
	Logger     *zap.Logger
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
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
						p.Logger.Error("unable to start server", zap.Error(err))
						if err := p.Shutdowner.Shutdown(); err != nil {
							p.Logger.Error("could not shutdown, exiting with os.Exit(1)", zap.Error(err))
							os.Exit(1)
						}
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
	logger     *zap.Logger
	httpServer *fasthttp.Server
	router     *router.Router
}

func (s *Server) Run() error {
	s.logger.Info("server running on port " + s.config.ServerPort)
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
		if s.logger.Core().Enabled(zap.DebugLevel) {
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
		if s.logger.Core().Enabled(zap.DebugLevel) {
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