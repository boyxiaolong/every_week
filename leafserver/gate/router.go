package gate

import (
	"server/game"
	"server/msg"
)

func init() {
	msg.Processor.Route(&msg.ReqLogin{}, game.ChanRPC)
}
