package handler

import (
	"encoding/json"
	"errors"

	gmg "github.com/brandenc40/green-mountain-grill"
	repo "github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/brandenc40/green-mountain-grill/server/respository/mapper"
	"github.com/google/uuid"
	"github.com/jasonlvhit/gocron"
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
	Scheduler   *gocron.Scheduler
}

// New -
func New(p Params) *Handler {
	h := &Handler{
		grill:              p.GrillClient,
		logger:             p.Logger,
		repo:               p.Repository,
		currentSessionUUID: uuid.Nil,
		scheduler:          p.Scheduler,
	}
	return h
}

// Handler -
type Handler struct {
	grill              gmg.Client
	logger             *zap.Logger
	repo               repo.Repository
	currentSessionUUID uuid.UUID
	scheduler          *gocron.Scheduler
	stopChannel        chan bool
	isMonitoring       bool
}

// GetGrillState -
func (c *Handler) GetGrillState(ctx *fasthttp.RequestCtx) {
	state, err := c.grill.GetState()
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
			c.logger.Error("error monitoring grill", zap.Error(err))
		}
	}()
}

// monitorGrill - check grill state and store to database once every minute
func (c *Handler) monitorGrill() error {
	if c.isMonitoring {
		c.stopMonitoringGrill()
		c.scheduler.Clear()
	}
	err := c.scheduler.Every(1).Minute().Do(c.storeGrillState)
	if err != nil {
		return err
	}
	c.stopChannel = c.scheduler.Start()
	c.isMonitoring = true
	return nil
}

func (c *Handler) stopMonitoringGrill() {
	c.stopChannel <- true
	c.isMonitoring = false
	close(c.stopChannel)
	c.stopChannel = nil
}

// storeGrillState - get current state and store to db
func (c *Handler) storeGrillState() (*gmg.State, error) {
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
