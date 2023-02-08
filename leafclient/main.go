package main

import (
	"encoding/binary"
	"fmt"
	"leafclient/msg"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3564")
	login := msg.ReqLogin{
		Account: "allen",
	}
	data, err := proto.Marshal(&login)
	if err != nil {
		log.Fatal("序列化失败 error:", err)
	}

	msglen := len(data)
	m := make([]byte, 4+msglen)
	binary.BigEndian.PutUint16(m, uint16(msglen+2))
	binary.BigEndian.PutUint16(m[2:], uint16(0))
	copy(m[4:], data)
	num, err := conn.Write(m)
	if err != nil {
		log.Fatal("消息发送失败 error:", err)
	}
	fmt.Println("leafclient send finish ", num, "msglen", msglen)
	//time.Sleep(1 * time.Second)
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		fmt.Println("read error", err, n)
	}
}
