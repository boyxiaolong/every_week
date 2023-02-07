package msg

import (
	"fmt"
	reflect "reflect"

	"github.com/name5566/leaf/network/protobuf"
)

var Processor = protobuf.NewProcessor()

func printMsg(id uint16, t interface{}) {
	fmt.Println("msg", id, t)
}
func init() {
	//Processor.SetByteOrder(true)
	Processor.Register(&ReqLogin{})

	Processor.Range(func(id uint16, t reflect.Type) { fmt.Println(id, t) })
}
