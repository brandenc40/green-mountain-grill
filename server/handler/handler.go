package handler

import (
	"encoding/json"
	"strconv"
	"time"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/brandenc40/green-mountain-grill/server/poller"
	repo "github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	contentTypeJSON = "application/json"
)

// Params -
type Params struct {
	fx.In

	Logger      *zap.Logger
	GrillClient gmg.Client
	Repository  repo.Repository
	Poller      *poller.Poller
}

// New -
func New(p Params) *Handler {
	h := &Handler{
		grill:  p.GrillClient,
		logger: p.Logger,
		repo:   p.Repository,
		poller: p.Poller,
	}
	return h
}

// Handler -
type Handler struct {
	grill    gmg.Client
	logger   *zap.Logger
	repo     repo.Repository
	poller   *poller.Poller
	upgrader websocket.FastHTTPUpgrader
}

// GetGrillState -
func (h *Handler) GetGrillState(ctx *fasthttp.RequestCtx) {
	state, err := h.grill.GetState()
	if err != nil {
		if _, ok := err.(gmg.GrillUnreachableErr); ok {
			ctx.Error(err.Error(), fasthttp.StatusServiceUnavailable)
			return
		}
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	marshalled, err := json.Marshal(state)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType(contentTypeJSON)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(marshalled)
}

// GetGrillID -
func (h *Handler) GetGrillID(ctx *fasthttp.RequestCtx) {
	id, err := h.grill.GetID()
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(id)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// GetGrillFirmware -
func (h *Handler) GetGrillFirmware(ctx *fasthttp.RequestCtx) {
	id, err := h.grill.GetFirmware()
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(id)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// StartPolling -
func (h *Handler) StartPolling(ctx *fasthttp.RequestCtx) {
	err := h.poller.StartPolling(time.Minute)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBodyString(h.poller.CurrentSession().String())
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// StopPolling -
func (h *Handler) StopPolling(ctx *fasthttp.RequestCtx) {
	h.poller.StopPolling()
	ctx.SetBodyString("OK")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// ViewSubscribers -
func (h *Handler) ViewSubscribers(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString(strconv.Itoa(h.poller.Subscribers()))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// GetSessionData -
func (h *Handler) GetSessionData(ctx *fasthttp.RequestCtx) {
	hist, err := h.repo.GetStateHistory(h.poller.CurrentSession())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	marshalled, err := json.Marshal(hist)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(marshalled)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
