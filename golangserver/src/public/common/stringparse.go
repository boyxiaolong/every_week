package common

import (
	"strconv"
	"strings"
	//"fmt"
)

type StringParse struct {
	Strs []string
}

func (m *StringParse) ParseString(str string, sep string) bool {
	m.Strs = make([]string, 0)

	if str == "" {
		return false
	}

	m.Strs = strings.Split(strings.Trim(str, sep), sep)

	//fmt.Printf("%v,%#v,%d",str,m.Strs,len(m.Strs))
	return true
}

func (m *StringParse) Set(str []string) {
	m.Strs = str
}

func (m *StringParse) Len() int {
	return len(m.Strs)
}

func (m *StringParse) GetString(index int) (str string) {
	if index < 0 || index >= len(m.Strs) {
		return
	}

	return m.Strs[index]
}

func (m *StringParse) GetInt(index int) (value int) {
	if index < 0 || index >= len(m.Strs) {
		return
	}

	value, _ = strconv.Atoi(m.Strs[index])
	return
}

func (m *StringParse) GetUInt64(index int) (value uint64) {
	if index < 0 || index >= len(m.Strs) {
		return
	}

	value, _ = strconv.ParseUint(m.Strs[index], 10, 32)
	return
}
