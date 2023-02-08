package internal

import (
	"fmt"
	"server/conf"
	"server/game"
	"server/msg"

	"github.com/name5566/leaf/gate"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	fmt.Println("TCPAddr", conf.Server.TCPAddr, conf.MaxMsgLen, conf.LenMsgLen)
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       msg.Processor,
		AgentChanRPC:    game.ChanRPC,
	}
}
