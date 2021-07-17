package handler

import (
	"encoding/json"
	"time"

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
	err := h.upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		closeChan := make(chan bool)
		go h.subWriter(ws, closeChan)
		h.subReader(ws, closeChan)
	})
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			h.logger.Error("websocket handshake error", zap.Error(err))
		}
		h.logger.Error("unable to upgrade to websocket", zap.Error(err))
		return
	}
}

func (h *Handler) subReader(ws *websocket.Conn, closeWriterChan chan bool) {
	ws.SetReadLimit(512)
	if err := ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		h.logger.Error("unable to set read deadline", zap.Error(err))
		closeWriterChan <- true
		return
	}
	ws.SetPongHandler(func(string) error {
		if err := ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			h.logger.Error("unable to set read deadline in pong handler", zap.Error(err))
			return err
		}
		return nil
	})
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
	closeWriterChan <- true
}

func (h *Handler) subWriter(ws *websocket.Conn, closeWriterChan chan bool) {
	pingTicker := time.NewTicker(pingPeriod)
	channel, unsubscribe := h.poller.Subscribe()
	defer func() {
		pingTicker.Stop()
		unsubscribe()
		func() { // close in the writer since this will run in a goroutine
			if err := ws.Close(); err != nil {
				h.logger.Error("error closing websocket conn", zap.Error(err))
			}
		}()
	}()
	if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		h.logger.Error("unable to set write deadline", zap.Error(err))
		return
	}
	err := ws.WriteMessage(websocket.TextMessage, []byte("hello there"))
	if err != nil {
		h.logger.Error("websocket write error", zap.Error(err))
		return
	}
	for {
		select {
		case m := <-channel:
			b, err := json.Marshal(m)
			if err != nil {
				h.logger.Error("unable to marshal", zap.Error(err))
				return
			}
			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				h.logger.Error("unable to set write deadline", zap.Error(err))
				return
			}
			err = ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				h.logger.Error("websocket write error", zap.Error(err))
				return
			}
		case <-pingTicker.C:
			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				h.logger.Error("unable to set write deadline", zap.Error(err))
				return
			}
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				h.logger.Error("websocket write ping error", zap.Error(err))
				return
			}
		case <-closeWriterChan:
			return
		}
	}
}
