package main

import (
	"encoding/binary"
	"fmt"
	"leafclient/msg"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3563")
	login := msg.ReqLogin{
		Account: "allen",
	}
	data, err := proto.Marshal(&login)
	if err != nil {
		log.Fatal("序列化失败 error:", err)
	}

	m := make([]byte, 4+len(data))
	binary.BigEndian.PutUint16(m, uint16(len(data)))
	binary.BigEndian.PutUint16(m[2:], 0)
	copy(m[4:], data)
	num, err := conn.Write(m)
	if err != nil {
		log.Fatal("消息发送失败 error:", err)
	}
	fmt.Println("leafclient send finish ", num)
	time.Sleep(1 * time.Second)
}
