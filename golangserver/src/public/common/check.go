package common

import (
	"public/config"
	"public/message/error_code"
	"public/message/msgtype"
)

func CHECK_ERROR(code int32, msg_type uint16) {
	if code != 0 {
		GStdout.Error("msg is error,msg_type:%v,msg_name:%v,code:%v,code_name:%v", msg_type, msgtype.MsgType_name[int32(msg_type)], code, error_code.ErrorCode_name[code])
		check_exit()
	} else {
		GStdout.Success("msg is success,msg_type:%v,msg_name:%v", msg_type, msgtype.MsgType_name[int32(msg_type)])
	}
}

func CHECK_RESULT(flag bool, msg_type uint16) {
	if flag == false {
		GStdout.Error("msg is false,msg_type:%v,msg_name:%v", msg_type, msgtype.MsgType_name[int32(msg_type)])
		check_exit()
	} else {
		GStdout.Success("msg is success,msg_type:%v,msg_name:%v", msg_type, msgtype.MsgType_name[int32(msg_type)])
	}
}

func check_exit() {
	if config.Mode == MODE_EXIT {
		QuitClient("check exit", -1)
	} else if config.Mode == MODE_WAIT {
		wait := make(chan int)
		<-wait
	} else if config.Mode == MODE_CONTINUE {
		return
	}
}
