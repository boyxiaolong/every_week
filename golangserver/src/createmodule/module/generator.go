package module

import "text/template"

import (
	"bytes"
	"os"
	"fmt"
)

type FileDescribe struct {
	name string
	content string
}

type Generator struct{
	files []*FileDescribe
}

func (m* Generator) GeneratorAllFile(module *ModuleDescribe) {
	result := m.existFile(module.GetModuleMkdir())

	if !result {
		os.Mkdir(module.GetModuleMkdir(), 0777)
	}

	m.files = make([]*FileDescribe,0)

	err, content , name := m.GeneratorModuleH(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	err, content , name = m.GeneratorModuleCpp(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	err, content , name = m.GeneratorModuleInterface(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	err, content , name = m.GeneratorMsgH(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	err, content , name = m.GeneratorMsgCpp(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	err, content , name = m.GeneratorPmH(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	err, content , name = m.GeneratorPmCpp(module)
	if err != nil {
		return
	}

	m.files = append(m.files,&FileDescribe{name,content})

	m.WriteAllFile(module)
}

func(m *Generator) WriteAllFile(module *ModuleDescribe) {
	for _, v := range m.files {
		m.WriteFile(module,v)
	}
}

func(m *Generator) WriteFile(module *ModuleDescribe,describe *FileDescribe) {
	path := module.GetModuleMkdir() + describe.name
	exist := m.existFile(path)

	if exist {
		fmt.Printf("filename is exist:%v\n", describe.name)
		return
	}

	file, _ := os.Create(path)
	file.WriteString(describe.content)
	file.Close()

	fmt.Println(fmt.Sprintf("write file success ,path %v, name %v", path, describe.name))
}

func (m *Generator) GeneratorModuleH(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("module_h").Parse(ModuleHTemplate)
	if err != nil {
		fmt.Println("template module_h new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template module_h Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModuleHName()
}

func (m *Generator) GeneratorModuleCpp(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("module_cpp").Parse(ModuleCppTemplate)
	if err != nil {
		fmt.Println("template module_cpp new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template module_cpp Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModuleCppName()
}

func (m *Generator) GeneratorModuleInterface(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("module_interface").Parse(ModuleInterfaceTemplate)
	if err != nil {
		fmt.Println("template module_interface new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template module_interface Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModuleInterfaceName()
}


func (m *Generator) GeneratorMsgH(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("msg_h").Parse(MsgHTemplate)
	if err != nil {
		fmt.Println("template msg_h new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template msg_h Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModuleMsgHName()
}

func (m *Generator) GeneratorMsgCpp(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("msg_cpp").Parse(MsgCppTemplate)
	if err != nil {
		fmt.Println("template msg_cpp new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template msg_cpp Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModuleMsgCppName()
}

func (m *Generator) GeneratorPmH(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("pm_h").Parse(PmHTemplate)
	if err != nil {
		fmt.Println("template pm_h new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template pm_h Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModulePmHName()
}

func (m *Generator) GeneratorPmCpp(module *ModuleDescribe) (error,string,string) {
	t, err := template.New("pm_cpp").Parse(PmCppTemplate)
	if err != nil {
		fmt.Println("template pm_cpp new error %v", err)
		return err,"", ""
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, module)
	if err != nil {
		fmt.Println("template pm_cpp Execute error %v", err)
		return err,"", ""
	}

	return nil,w.String(), module.GetModulePMCppName()
}

func (m *Generator)existFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}