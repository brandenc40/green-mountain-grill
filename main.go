package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/brandenc40/gmg/internal/handler"
	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

type config struct {
	GrillIP   string `yaml:"grill_ip"`
	GrillPort int    `yaml:"grill_port"`
}

func main() {
	cfg := readConfig()
	logger := newLogger()

	h := handler.New(handler.Params{
		GrillIP:   net.ParseIP(cfg.GrillIP),
		GrillPort: cfg.GrillPort,
		Logger:    logger,
	})

	r := router.New()
	r.GET("/state", withMiddleware(h.GetGrillState, logger))
	r.GET("/id", withMiddleware(h.GetGrillID, logger))
	r.GET("/firmware", withMiddleware(h.GetGrillFirmware, logger))

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

func withMiddleware(h fasthttp.RequestHandler, logger *logrus.Logger) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.Info("[REQUEST ", string(ctx.Method()), " ", string(ctx.Path()), "]")
		t := time.Now()
		h(ctx)
		logger.Info("[RESPONSE ", ctx.Response.StatusCode(), " ", time.Since(t).String(), "]")
	}
}
