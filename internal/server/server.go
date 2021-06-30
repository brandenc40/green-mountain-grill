package server

import (
	"bytes"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/brandenc40/gmg/internal/handler"
	repo "github.com/brandenc40/gmg/internal/respository"
	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

type config struct {
	GrillIP   string `yaml:"grill_ip"`
	GrillPort int    `yaml:"grill_port"`
}

func Run() {
	cfg := readConfig()
	logger := newLogger()

	h := handler.New(handler.Params{
		GrillIP:    net.ParseIP(cfg.GrillIP),
		GrillPort:  cfg.GrillPort,
		Logger:     logger,
		Repository: newRepo(),
	})

	r := router.New()
	r.GET("/api/state", withLogging(h.GetGrillState, logger))
	r.GET("/api/id", withLogging(h.GetGrillID, logger))
	r.GET("/api/firmware", withLogging(h.GetGrillFirmware, logger))
	r.GET("/api/session", withLogging(h.GetSessionData, logger))
	r.POST("/api/session/new", withLogging(h.NewSession, logger))

	logger.Info("Running on port :8080")
	logger.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

func readConfig() *config {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	var c config
	if err = yaml.Unmarshal(file, &c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &c
}

func newLogger() *logrus.Logger {
	level := logrus.DebugLevel
	if os.Getenv("ENVIRONMENT") == "production" {
		level = logrus.InfoLevel
	}
	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
	logger.Debug("logging at DEBUG level")
	return logger
}

func newRepo() repo.Repository {
	r, err := repo.New()
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func withLogging(h fasthttp.RequestHandler, logger *logrus.Logger) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var b bytes.Buffer
		// log request
		b.WriteString("[REQUEST ")
		b.Write(ctx.Method())
		b.WriteString(" ")
		b.Write(ctx.Path())
		b.WriteString("]")
		logger.Info(b.String())
		b.Reset()

		// handle and start timer
		t := time.Now()
		h(ctx)

		// log response
		b.WriteString("[RESPONSE ")
		b.WriteString(strconv.Itoa(ctx.Response.StatusCode()))
		if ctx.Response.StatusCode() != 200 {
			b.WriteString(" (")
			b.Write(ctx.Response.Body())
			b.WriteString(")")
		}
		b.WriteString(" ")
		b.WriteString(time.Since(t).String())
		b.WriteString("]")
		logger.Info(b.String())
	}
}
