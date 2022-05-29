package handler

import (
	"encoding/json"
	"strconv"
	"time"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/brandenc40/green-mountain-grill/server/poller"
	repo "github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	_contentTypeJSON = "application/json"
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
	grill  gmg.Client
	logger *zap.Logger
	repo   repo.Repository
	poller *poller.Poller
}

// GetGrillState -
func (h *Handler) GetGrillState(ctx *fiber.Ctx) error {
	state, err := h.grill.GetState()
	if err != nil {
		if _, ok := err.(gmg.GrillUnreachableErr); ok {
			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(state)
}

// GetGrillID -
func (h *Handler) GetGrillID(ctx *fiber.Ctx) error {
	id, err := h.grill.GetID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).Send(id)
}

// GetGrillFirmware -
func (h *Handler) GetGrillFirmware(ctx *fiber.Ctx) error {
	firmware, err := h.grill.GetFirmware()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).Send(firmware)
}

// StartPolling -
func (h *Handler) StartPolling(ctx *fiber.Ctx) error {
	err := h.poller.StartPolling(time.Minute)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	sessionUUID := h.poller.CurrentSession().String()
	return ctx.Status(fiber.StatusOK).SendString(sessionUUID)
}

// StopPolling -
func (h *Handler) StopPolling(ctx *fiber.Ctx) error {
	h.poller.StopPolling()
	return ctx.Status(fiber.StatusOK).SendString("OK")
}

// ViewSubscribers -
func (h *Handler) ViewSubscribers(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString(strconv.Itoa(h.poller.Subscribers()))
}

// GetSessionData -
func (h *Handler) GetSessionData(ctx *fiber.Ctx) error {
	hist, err := h.repo.GetStateHistory(h.poller.CurrentSession())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	marshalled, err := json.Marshal(hist)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).Send(marshalled)
}
