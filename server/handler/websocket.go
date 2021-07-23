package handler

import (
	"encoding/json"
	"time"

	"github.com/brandenc40/green-mountain-grill/server/respository/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

const (
	// Time allowed to write to the client.
	writeWait = 2 * time.Second
	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second
	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = 5 * time.Second
)

func (h *Handler) BuildPollerSubscriberWSHandler() fiber.Handler {
	return websocket.New(func(ws *websocket.Conn) {
		channel, unsubscribe := h.poller.Subscribe()
		defer unsubscribe()
		h.subWriter(ws, channel)
		h.subReader(ws)
		h.logger.Info("unsubscribe() called")
	})
}

func (h *Handler) subReader(ws *websocket.Conn) {
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		t, p, err := ws.ReadMessage()
		h.logger.Info("received: ",
			zap.Int("type", t),
			zap.ByteString("p", p))
		if err != nil {
			h.logger.Error("received err: ", zap.Error(err))
			break
		}
	}
	h.logger.Info("subReader() ended")
}

func (h *Handler) subWriter(ws *websocket.Conn, channel chan *model.GrillState) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		func() { // close in the writer since this will run in a goroutine
			if err := ws.Close(); err != nil {
				h.logger.Error("error closing websocket conn", zap.Error(err))
			}
		}()
	}()
	statHist, err := h.repo.GetStateHistory(h.poller.CurrentSession())
	if err != nil {
		h.logger.Error("unable to get state history", zap.Error(err))
		return
	}
	data, err := json.Marshal(statHist)
	if err != nil {
		h.logger.Error("unable to marshal state history", zap.Error(err))
		return
	}
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	err = ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		h.logger.Error("websocket write error", zap.Error(err))
		return
	}
	for {
		select {
		case m, ok := <-channel:
			if !ok {
				h.logger.Info("The poller closed the channel")
				return
			}
			statHist = append(statHist, m)
			b, err := json.Marshal(statHist)
			if err != nil {
				h.logger.Error("unable to marshal state history", zap.Error(err))
				return
			}
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err = ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				h.logger.Error("websocket write error", zap.Error(err))
				return
			}
		case <-pingTicker.C:
			h.logger.Info("tick")
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				h.logger.Error("websocket write ping error", zap.Error(err))
				return
			}
		}
	}
}
