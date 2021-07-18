package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brandenc40/green-mountain-grill/server/respository/model"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second
	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// SubscribeToPoller -
func (h *Handler) SubscribeToPoller(ctx *fasthttp.RequestCtx) {
	err := h.webSocket.Upgrade(ctx, func(ws *websocket.Conn) {
		channel, unsubscribe := h.poller.Subscribe()
		defer unsubscribe()
		go h.withRecover(func() {
			h.subWriter(ws, channel)
		})
		h.withRecover(func() {
			h.subReader(ws)
		})
		h.logger.Info("unsubscribe() called")
	})
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			h.logger.Error("websocket handshake error", zap.Error(err))
		}
		h.logger.Error("unable to upgrade to websocket", zap.Error(err))
		return
	}
}

func (h *Handler) subReader(ws *websocket.Conn) {
	ws.SetReadLimit(512)
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { _ = ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
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
	for {
		select {
		case m, ok := <-channel:
			if !ok {
				// The poller closed the channel.
				h.logger.Info("The poller closed the channel")
				return
			}
			b, err := json.Marshal(m)
			if err != nil {
				h.logger.Error("unable to marshal", zap.Error(err))
				return
			}
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
			err = ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				h.logger.Error("websocket write error", zap.Error(err))
				return
			}
		case <-pingTicker.C:
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				h.logger.Error("websocket write ping error", zap.Error(err))
				return
			}
		}
	}
}

func (h *Handler) withRecover(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
			h.logger.Error("panic recovered: ", zap.Error(err))
		}
	}()

	fn()
}
