package generator

const RpcClientTemplateH = `#pragma once
#include <chrono>
#include "{{$.ModelFileName}}"
#include "session_address.h"
#include "rpc_client.h"
#include "{{$.GeneratedRpcTemplateDefineFilename}}"
{{range $serviceIndex, $service := $.Services}}
namespace {{$service.NamespaceName}}
{
  {{range $i, $method := $service.Methods}}
  async::Future<ZResult<{{$.ContainerNamespace}}{{$method.OutputType}}>> {{$method.MethodName}}(SessionAddress target,
     const {{$.ContainerNamespace}}{{$method.InputType}} request,
     std::chrono::seconds timeout =  std::chrono::seconds::zero());
  {{end}}
};
{{end}}
`

const RpcClientTemplateCpp = `#include "stdafx.h"
#include "rpc/{{$.GeneratedRpcClientTemplateHFilename}}"
{{range $serviceIndex, $service := $.Services}}
namespace {{$service.NamespaceName}}
{
{{range $i, $method := $service.Methods}}
async::Future<ZResult<{{$.ContainerNamespace}}{{$method.OutputType}}>> {{$method.MethodName}}(SessionAddress target,
    const {{$.ContainerNamespace}}{{$method.InputType}} request,
	std::chrono::seconds timeout)
{
  co_return co_await RpcClientInst.Call<{{$.ContainerNamespace}}{{$method.OutputType}}>(
    target, {{$service.MethodsEnumName}}::k{{$service.Name}}{{$method.EnumName}}, request, timeout);
}
{{end}}
};
{{end}}

`