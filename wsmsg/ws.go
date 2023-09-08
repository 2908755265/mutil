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

type WsPusher struct {
	msync.MapInterface[string, *websocket.Conn]
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

	return conn.WriteMessage(websocket.TextMessage, bts)
}

func NewWSPusher() *WsPusher {
	return &WsPusher{
		MapInterface: msync.NewSyncMap[string, *websocket.Conn](),
	}
}
