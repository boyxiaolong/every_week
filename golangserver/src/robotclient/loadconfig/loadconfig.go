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
	ChargeIp          atomic.Value
	ChargePort        atomic.Value
)

func init() {
	command.GCommand.RegCommand("setip", SetIp, "setip")
	command.GCommand.RegCommand("showip", ShowIp, "showip")
	command.GCommand.RegCommand("showchargeip", ShowChargeIp, "showchargeip")
	Load()
}

func Load() {
	Ip.Store(config.Cfg.MustValue("net", "ip"))
	Port.Store(config.Cfg.MustValue("net", "port"))
	Account = uint64(config.Cfg.MustInt64("test", "account", 1))
	TestAllType = config.Cfg.MustInt("test", "testall_type", 1)

	ChargeIp.Store(config.Cfg.MustValue("charge", "ip"))
	ChargePort.Store(config.Cfg.MustValue("charge", "port"))
}

//Login Login
func SetIp(str *common.StringParse) (err error) {
	if str.Len() < 3 {
		common.GStdout.Error("login param error")
		return
	}

	Ip.Store(str.GetString(1))
	Port.Store(str.GetString(2))

	common.GStdout.Success("set ip success")

	return
}

//Login Login
func ShowIp(str *common.StringParse) (err error) {
	common.GStdout.Info("ip:%v, port:%v", Ip.Load().(string), Port.Load().(string))
	return
}

//ShowChargeIp ShowChargeIp
func ShowChargeIp(str *common.StringParse) (err error) {
	common.GStdout.Info("charge ip:%v, port:%v", ChargeIp.Load().(string), ChargePort.Load().(uint32))
	return
}

func GetIp() string {
	return Ip.Load().(string)
}

func GetPort() string {
	return Port.Load().(string)
}

func GetChargeIp() string {
	return ChargeIp.Load().(string)
}

func GetChargePort() string {
	return ChargePort.Load().(string)
}

