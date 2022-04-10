package connect

import (
	"encoding/binary"
	//"flag"
	"fmt"
	"public/common"

	"public/link"
)

func GetConnectSessionByMsgHander(ip string, port string, server_type uint64, msghandler MessageHandler) (*Session, error) {
	addr := fmt.Sprintf("%v:%v", ip, port)

	lsession, err := link.Connect("tcp", addr, link.Packet(4, 1024*1024, 1024*1024, binary.LittleEndian, ConnectCodec{}))
	if err != nil {
		common.GStdout.Error("%v", err)
		return nil, err
	}

	return NewSessionByMsgHander(lsession, server_type, msghandler), err
}
