package loadconfig

import (
	"public/command"
	"public/common"
	"public/config"
	"sync/atomic"
)

var (
	Ip          atomic.Value
	Port        atomic.Value
	Account     uint64
	TestAllType int
)

func init() {
	command.GCommand.RegCommand("setip", SetIp, "setip")
	command.GCommand.RegCommand("setport", SetPort, "setport")
	command.GCommand.RegCommand("showip", ShowIp, "showip")
	Load()
}

func Load() {
	Ip.Store(config.Cfg.MustValue("net", "ip"))
	Port.Store(config.Cfg.MustValue("net", "port"))
	Account = uint64(config.Cfg.MustInt64("test", "account", 1))
	TestAllType = config.Cfg.MustInt("test", "testall_type", 1)
}

//SetPort SetPort
func SetPort(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("set port param error")
		return
	}

	Port.Store(str.GetString(1))

	common.GStdout.Success("set port %v success", Port.Load().(string))

	return
}

//SetIp SetIp
func SetIp(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("login param error")
		return
	}
	info := "";

	Ip.Store(str.GetString(1))
	info += Ip.Load().(string);
	if str.Len() >= 3 {
		Port.Store(str.GetString(2))
		info += ":" + Port.Load().(string);
	}

	common.GStdout.Success("set ip %v success", info)

	return
}

//Login Login
func ShowIp(str *common.StringParse) (err error) {
	common.GStdout.Info("ip:%v, port:%v", Ip.Load().(string), Port.Load().(string))
	return
}

func GetIp() string {
	return Ip.Load().(string)
}

func GetPort() string {
	return Port.Load().(string)
}
