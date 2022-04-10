package loadconfig

import "public/config"

var (
	SqlAddr  string
)

func init() {
	Load()
}

func Load() {
	SqlAddr = config.Cfg.MustValue("db", "sql_addr")
}