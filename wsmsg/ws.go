package wsmsg

import (
	"encoding/json"
	"errors"
	"github.com/2908755265/mutil/msync"
	"github.com/gorilla/websocket"
)

var (
	ConnNotExistErr = errors.New("conn not exist")
)

type MessagePusher interface {
	Push(key string, val any) error
}

type ErrHandler func(key string, val any, conn *websocket.Conn, wrapper *WsPusher)

type WsPusher struct {
	msync.MapInterface[string, *websocket.Conn]
	errHandler ErrHandler
}

func (p *WsPusher) Push(key string, val any) error {
	conn, ok := p.Load(key)
	if !ok {
		return ConnNotExistErr
	}

	bts, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, bts)
	if err != nil {
		p.errHandler(key, val, conn, p)
	}
	return err
}

func defaultErrHandler(key string, _ any, conn *websocket.Conn, wrapper *WsPusher) {
	conn.Close()
	wrapper.Delete(key)
}

func NewWSPusher() *WsPusher {
	return &WsPusher{
		MapInterface: msync.NewSyncMap[string, *websocket.Conn](),
		errHandler:   defaultErrHandler,
	}
}

func WithErrorHandler(p *WsPusher, h ErrHandler) {
	p.errHandler = h
}
