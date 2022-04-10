package command

import (
	"public/common"
)

type Callback func(str *common.StringParse) (err error)

type KeyCommand struct {
	Command  string
	Callback Callback
	Note     string
}
