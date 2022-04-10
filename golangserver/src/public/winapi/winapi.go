package winapi

import (
	"syscall"
	"public/common"
)

func SetConsoleMode() {
	//C.SetWinConsoleMode()

	dll, err := syscall.LoadLibrary("conf/SetConsoleMode.dll")
	if err != nil{
		common.GStdout.Error("load SetConsoleMode error :%v", err)
		return
	}

	defer syscall.FreeLibrary(dll)

	proc, err := syscall.GetProcAddress(dll, "SetWinConsoleMode")
	if err!=nil {
		common.GStdout.Error("get SetWinConsoleMode error :%v", err)
		return
	}

	r, _, _ := syscall.Syscall(uintptr(proc), 0, 0, 0,0)

	if int(r) != 1 {
		common.GStdout.Error("Syscall SetWinConsoleMode error :%v", err)
		return
	}

	common.GStdout.Success("SetWinConsoleMode success")
	return
}