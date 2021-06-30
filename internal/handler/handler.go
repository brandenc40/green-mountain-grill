package handler

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/brandenc40/gmg/grillclient"
	repo "github.com/brandenc40/gmg/internal/respository"
	"github.com/brandenc40/gmg/internal/respository/mapper"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	contentTypeJSON = "application/json"
)

// Handler -
type Handler struct {
	grill              grillclient.Client
	logger             *logrus.Logger
	repo               repo.Repository
	currentSessionUUID uuid.UUID
}

// Params -
type Params struct {
	GrillIP    net.IP
	GrillPort  int
	Logger     *logrus.Logger
	Repository repo.Repository
}

// New -
func New(p Params) *Handler {
	h := &Handler{
		grill: grillclient.New(grillclient.Params{
			GrillIP:   p.GrillIP,
			GrillPort: p.GrillPort,
			Logger:    p.Logger,
		}),
		logger:             p.Logger,
		repo:               p.Repository,
		currentSessionUUID: uuid.Nil,
	}
	if h.logger == nil {
		h.logger = logrus.New()
	}
	return h
}

// GetGrillState -
func (c *Handler) GetGrillState(ctx *fasthttp.RequestCtx) {
	state, err := c.grill.GetState()
	if err != nil {
		if _, ok := err.(grillclient.GrillUnreachableErr); ok {
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
func (c *Handler) GetGrillID(ctx *fasthttp.RequestCtx) {
	id, err := c.grill.GetID()
	if err != nil {
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
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(id)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// NewSession -
func (c *Handler) NewSession(ctx *fasthttp.RequestCtx) {
	if !c.grill.IsAvailable() {
		ctx.Error("grill is not available", fasthttp.StatusServiceUnavailable)
		return
	}
	u, err := uuid.NewRandom()
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	c.currentSessionUUID = u
	c.monitorGrillAsync()
	ctx.SetBodyString(c.currentSessionUUID.String())
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (c *Handler) GetSessionData(ctx *fasthttp.RequestCtx) {
	hist, err := c.repo.GetStateHistory(c.currentSessionUUID)
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

// monitorGrillAsync - start monitoring and report errors asynchronously
func (c *Handler) monitorGrillAsync() {
	go func() {
		if err := c.monitorGrill(); err != nil {
			c.logger.WithError(err).Error("error monitoring grill")
		}
	}()
}

// monitorGrill - check grill state and store to database once every minute
func (c *Handler) monitorGrill() error {
	for true {
		if _, err := c.storeGrillState(); err != nil {
			return err
		}
		time.Sleep(time.Minute)
	}
	return nil
}

// storeGrillState - get current state and store to db
func (c *Handler) storeGrillState() (*grillclient.State, error) {
	if c.currentSessionUUID == uuid.Nil {
		return nil, errors.New("no grill session UUID has been set")
	}
	state, err := c.grill.GetState()
	if err != nil {
		return nil, err
	}
	stateModel := mapper.GrillStateEntityToModel(state, c.currentSessionUUID)
	if err := c.repo.InsertStateData(stateModel); err != nil {
		return nil, err
	}
	return state, nil
}
