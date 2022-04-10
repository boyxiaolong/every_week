package generator

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/ahmetb/go-linq"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"

	"protoc-gen-maprpc/plugin"
)

type Method struct {
	name       string
	inputType  string
	outputType string
	messageID  int32
}

func (fd *Method) Name() string {
	return strings.Title(fd.name)
}

func (fd *Method) MethodName() string {
	return fd.name
}

func (fd *Method) EnumName() string {
	return fd.name
}

func (fd *Method) InputType() string {
	split := strings.Split(fd.inputType, ".")
	if len(split) == 0 {
		return ""
	}
	return split[len(split)-1]
}

func (fd *Method) OutputType() string {
	split := strings.Split(fd.outputType, ".")
	if len(split) == 0 {
		return ""
	}
	return split[len(split)-1]
}

func (fd *Method) MethodID() int32 {
	return fd.messageID
}

type Service struct {
	name string

	Methods []*Method
}

func (fd *Service) Name() string {
	return strings.Title(fd.name)
}

func (fd *Service) MethodsEnumName() string {
	return fmt.Sprintf("%sRpcMethodID", fd.Name())
}

func (fd *Service) NamespaceName() string {
	return fmt.Sprintf("%srpc", strings.ToLower(fd.name))
}

type FileDescriptor struct {
	fileName    string
	packageName string

	Services []*Service
}

func (fd *FileDescriptor) fileNameWithoutExt() string {
	return strings.TrimSuffix(strings.ToLower(fd.fileName), ".proto")
}

func (fd *FileDescriptor) GeneratedFilename() string {
	return fmt.Sprintf("%s.gayrpc.h", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) GeneratedRpcClientTemplateHFilename() string {
	return fmt.Sprintf("%s_client.h", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) GeneratedRpcClientTemplateCppFilename() string {
	return fmt.Sprintf("%s_client.cpp", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) GeneratedRpcServerTemplateHFilename() string {
	return fmt.Sprintf("%s_service_base.h", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) GeneratedRpcServerTemplateCppFilename() string {
	return fmt.Sprintf("%s_service_base.cpp", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) GeneratedRpcTemplateDefineFilename() string {
	return fmt.Sprintf("%s_define.h", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) MacroName() string {
	packetMacro := strings.Join(strings.Split(fd.packageName, "."), "_")
	return fmt.Sprintf("%s_%s_H",
		strings.ToUpper(packetMacro),
		strings.ToUpper(fd.fileNameWithoutExt()))
}

func (fd *FileDescriptor) ModelFileName() string {
	return fmt.Sprintf("%s.pb.h", fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) Namespace() string {
	return strings.Title(fd.fileNameWithoutExt())
}

func (fd *FileDescriptor) ContainerNamespace() string {
	return fmt.Sprintf("%s::", strings.Join(strings.Split(fd.packageName, "."), "::"))
}

func (fd *FileDescriptor) PackageNames() []string {
	return strings.Split(fd.packageName, ".")
}

// ---------------------------------------------------------------------------------------------------------------------

type Generator struct {
	*plugin.BaseGenerator
}

func (g *Generator) GenerateAllFiles() {
	for _, fd := range g.Files {
		if !linq.From(g.Request.FileToGenerate).Contains(fd.Descriptor.GetName()) {
			continue
		}
		d := &FileDescriptor{
			fileName:    fd.Descriptor.GetName(),
			packageName: fd.Descriptor.GetPackage(),
		}
		linq.From(fd.Descriptor.GetService()).SelectT(func(fd *descriptor.ServiceDescriptorProto) *Service {
			d := &Service{
				name: fd.GetName(),
			}
			linq.From(fd.GetMethod()).SelectT(func(fd *descriptor.MethodDescriptorProto) *Method {
				var options = make(map[string]string)
				linq.From(strings.Fields(proto.CompactTextString(fd.GetOptions()))).
					WhereT(func(s string) bool {
						return len(strings.TrimSpace(s)) > 0
					}).
					SelectT(func(s string) linq.KeyValue {
						lines := strings.Split(strings.TrimSpace(s), ":")
						if len(lines) != 2 {
							g.Fail("miss option in:", s)
						}
						return linq.KeyValue{
							Key:   lines[0],
							Value: lines[1],
						}
					}).ToMap(&options)

				var id int64
				if op, ok := options["51002"]; ok && op != "" {
					var err error
					id, err = strconv.ParseInt(op, 10, 32)
					if err != nil {
						g.Error(err, "method id illegal")
					}
				} else {
					g.Fail("miss message id of service function name:", *fd.Name)
				}

				return &Method{
					name:       fd.GetName(),
					inputType:  fd.GetInputType(),
					outputType: fd.GetOutputType(),
					messageID:  int32(id),
				}
			}).ToSlice(&d.Methods)

			idMap := make(map[int32]string)
			for _, v := range d.Methods {
				if existName, ok := idMap[v.messageID]; ok {
					g.Fail(v.name,
						"'s message_id :",
						strconv.Itoa(int(v.messageID)),
						" already exists of service function name:", existName)
				}
				idMap[v.messageID] = v.name
			}
			return d
		}).ToSlice(&d.Services)
		content, name := g.printRpcClientTemplateHFile(d)
		g.Response.File = append(g.Response.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(name),
			Content: proto.String(content),
		})

		content, name = g.printRpcClientTemplateCppFile(d)
		g.Response.File = append(g.Response.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(name),
			Content: proto.String(content),
		})

		content, name = g.printRpcServerTemplateHFile(d)
		g.Response.File = append(g.Response.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(name),
			Content: proto.String(content),
		})

		content, name = g.printRpcServerTemplateCppFile(d)
		g.Response.File = append(g.Response.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(name),
			Content: proto.String(content),
		})

		content, name = g.printRpcTemplateDefineFile(d)
		g.Response.File = append(g.Response.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(name),
			Content: proto.String(content),
		})
	}
}

func (g *Generator) printRpcClientTemplateHFile(model *FileDescriptor) (content string, name string) {
	t, err := template.New("rpc_client_template_h").Parse(RpcClientTemplateH)
	if err != nil {
		g.Error(err, "rpc_client_template_h parse failed")
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, model)
	if err != nil {
		g.Error(err, "execute rpc_client_template_h")
	}

	return w.String(), model.GeneratedRpcClientTemplateHFilename()
}

func (g *Generator) printRpcClientTemplateCppFile(model *FileDescriptor) (content string, name string) {
	t, err := template.New("rpc_client_template_cpp").Parse(RpcClientTemplateCpp)
	if err != nil {
		g.Error(err, "rpc_client_template_cpp parse failed")
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, model)
	if err != nil {
		g.Error(err, "execute rpc_client_template_cpp")
	}

	return w.String(), model.GeneratedRpcClientTemplateCppFilename()
}

func (g *Generator) printRpcServerTemplateHFile(model *FileDescriptor) (content string, name string) {
	t, err := template.New("rpc_server_template_h").Parse(RpcServerTemplateH)
	if err != nil {
		g.Error(err, "rpc_server_template_h parse failed")
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, model)
	if err != nil {
		g.Error(err, "execute rpc_server_template_h")
	}

	return w.String(), model.GeneratedRpcServerTemplateHFilename()
}

func (g *Generator) printRpcServerTemplateCppFile(model *FileDescriptor) (content string, name string) {
	t, err := template.New("rpc_server_template_cpp").Parse(RpcServerTemplateCpp)
	if err != nil {
		g.Error(err, "rpc_server_template_cpp parse failed")
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, model)
	if err != nil {
		g.Error(err, "execute rpc_server_template_cpp")
	}

	return w.String(), model.GeneratedRpcServerTemplateCppFilename()
}

func (g *Generator) printRpcTemplateDefineFile(model *FileDescriptor) (content string, name string) {
	t, err := template.New("rpc_template_define").Parse(RpcServerTemplateDefine)
	if err != nil {
		g.Error(err, "rpc_template_define parse failed")
	}

	w := bytes.NewBuffer(make([]byte, 0, 1024))
	err = t.Execute(w, model)
	if err != nil {
		g.Error(err, "execute rpc_template_define")
	}

	return w.String(), model.GeneratedRpcTemplateDefineFilename()
}

func NewGenerator(name string) *Generator {
	g := new(Generator)
	g.BaseGenerator = plugin.NewBaseGenerator(name)
	return g
}
