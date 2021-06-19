package handler

import (
	"encoding/json"
	"net"

	"github.com/brandenc40/gmg/grillclient"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	contentTypeJSON = "application/json"
)

// Handler -
type Handler struct {
	grill  grillclient.Client
	logger *logrus.Logger
}

// Params -
type Params struct {
	GrillIP   net.IP
	GrillPort int
	Logger    *logrus.Logger
}

// New -
func New(p Params) *Handler {
	return &Handler{
		grill: grillclient.New(grillclient.Params{
			IP:        p.GrillIP,
			GrillPort: p.GrillPort,
			Logger:    p.Logger,
		}),
		logger: p.Logger,
	}
}

// GetGrillState -
func (c *Handler) GetGrillState(ctx *fasthttp.RequestCtx) {
	state, err := c.grill.GetState()
	if err != nil {
		c.logger.WithError(err).Error("unable to get grill state")
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	marshalled, err := json.Marshal(state)
	if err != nil {
		c.logger.WithError(err).Error("error unmarshalling grill state")
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType(contentTypeJSON)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(marshalled)
}

// GetGrillID -
func (c *Handler) GetGrillID(ctx *fasthttp.RequestCtx) {
	id, err := c.grill.GetID()
	if err != nil {
		c.logger.WithError(err).Error("unable to get grill id")
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(id)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// GetGrillFirmware -
func (c *Handler) GetGrillFirmware(ctx *fasthttp.RequestCtx) {
	id, err := c.grill.GetFirmware()
	if err != nil {
		c.logger.WithError(err).Error("unable to get grill firmware")
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(id)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
