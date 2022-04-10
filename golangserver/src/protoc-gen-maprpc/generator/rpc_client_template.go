package generator

const RpcClientTemplateH = `#pragma once
#include <chrono>
#include "map/region_id.h"
#include "{{$.ModelFileName}}"
#include "session_address.h"
#include "rpc_client.h"
#include "{{$.GeneratedRpcTemplateDefineFilename}}"
{{range $serviceIndex, $service := $.Services}}
namespace {{$service.NamespaceName}}
{
  {{range $i, $method := $service.Methods}}
  async::Future<ZResult<{{$.ContainerNamespace}}{{$method.OutputType}}>> {{$method.MethodName}}(RegionId region_id,
   {{$.ContainerNamespace}}{{$method.InputType}} request,
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
async::Future<ZResult<{{$.ContainerNamespace}}{{$method.OutputType}}>> {{$method.MethodName}}(RegionId region_id,
    {{$.ContainerNamespace}}{{$method.InputType}} request,
	std::chrono::seconds timeout)
{
  request.set_region_id(region_id.GetId());
  co_return co_await RpcClientInst.Call<{{$.ContainerNamespace}}{{$method.OutputType}}>(
    region_id.GetServerAddress(), {{$service.MethodsEnumName}}::k{{$service.Name}}{{$method.EnumName}}, request, timeout);
}
{{end}}
};
{{end}}

`