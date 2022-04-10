package module

import (
	"strings"
	"fmt"
)

type ModuleDescribe struct {
	name string
}

func(m* ModuleDescribe) SetModuleName(name string) {
	m.name = strings.ToLower(name)
}

func(m* ModuleDescribe) GetTitleName() string {
	return strings.Title(m.name)
}

func(m* ModuleDescribe) GetUpperName() string {
	return strings.ToUpper(m.name)
}

func(m* ModuleDescribe) GetName() string {
	return m.name
}

func (m* ModuleDescribe)GetClassName() string {
	temp := strings.Split(m.name, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])

		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i]) // + string(vv[i+1])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

func(m *ModuleDescribe) GetModuleHName() string {
	return fmt.Sprintf("player_%s_module.h", m.name)
}

func(m *ModuleDescribe) GetModuleCppName() string {
	return fmt.Sprintf("player_%s_module.cpp", m.name)
}

func(m *ModuleDescribe) GetModuleMsgHName() string {
	return fmt.Sprintf("player_%s_msg_handle.h", m.name)
}

func(m *ModuleDescribe) GetModuleMsgCppName() string {
	return fmt.Sprintf("player_%s_msg_handle.cpp", m.name)
}

func(m *ModuleDescribe) GetModulePmHName() string {
	return fmt.Sprintf("player_%s_pm_command.h", m.name)
}

func(m *ModuleDescribe) GetModulePMCppName() string {
	return fmt.Sprintf("player_%s_pm_command.cpp", m.name)
}

func(m *ModuleDescribe) GetModuleInterfaceName() string {
	return fmt.Sprintf("player_%s_module_interface.h", m.name)
}

func(m *ModuleDescribe) GetModuleMkdir() string {
	return fmt.Sprintf("player_%v\\", m.name)
}