package upstoxapi

import (
	"github.com/gorilla/websocket"
)

type SocketClient struct {
	socket         *websocket.Conn
	accessToken    string
	instrumentkeys []string
	mode           string
}

const (
	ModeLTPC   = "ltpc"
	ModeFull   = "full"
	ModeOption = "option_greeks"
	ModeD30    = "fulld30"

	MethodSubscribe   = "sub"
	MethodUnsubscribe = "unsub"
	MethodChangeMode  = "change_mode"
)

func NewSocketClient()
