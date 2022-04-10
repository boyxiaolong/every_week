package generator

const RpcServerTemplateH = `#pragma once
#include "rpc_service.h"
#include "{{$.ModelFileName}}"
#include "rpc/{{$.GeneratedRpcTemplateDefineFilename}}"

{{range $serviceIndex, $service := $.Services}}
class {{$service.Name}}RpcServiceBase : public core::RpcDispatcher
{
public:
  virtual ~{{$service.Name}}RpcServiceBase() = default;

private:
  bool DispatchRequest(uint32_t method, 
	core::RpcRequestContext ctx, 
    std::string_view request) override;

protected:{{range $i, $method := $service.Methods}}
  virtual void {{$method.MethodName}}(
    core::RpcRequestContext ctx, 
    const {{$.ContainerNamespace}}{{$method.InputType}}& request);
  {{end}}
};

{{end}}

`

const RpcServerTemplateCpp = `#include "stdafx.h"
#include "rpc/{{$.GeneratedRpcServerTemplateHFilename}}"
#include "check.h"
#include "server/server_status.h"
#include "error_code.pb.h"
#include "applog.h"

{{range $serviceIndex, $service := $.Services}}
bool {{$service.Name}}RpcServiceBase::DispatchRequest(
 uint32_t method, core::RpcRequestContext ctx, std::string_view request)
{
  if (!ServerStatusInst::Singleton()->IsReady())
  {
    ctx.ResponseErr(ErrorCode::kECGeneralServiceNotReady);
    return true;
  }

  auto server_status = ServerStatusInst::Singleton()->GetServiceStatus();
  if (server_status != ServiceStatus::kRunning 
    && server_status != ServiceStatus::kPrepareShutdown)
  {
    ctx.ResponseErr(ErrorCode::kECGeneralServerShutingdown);
    return true;
  }

  switch (method) 
  {
  {{range $i, $method := $service.Methods}}
	case {{$service.MethodsEnumName}}::k{{$service.Name}}{{$method.EnumName}}:
    {
      {{$.ContainerNamespace}}{{$method.InputType}} msg;
      CHECKF(msg.ParseFromArray(request.data(), static_cast<int>(request.size())));
      {{$method.MethodName}}(ctx, msg);
    }
    break;{{end}}
  default:
     LOGERR("Invalid rpc service method: {}", static_cast<uint32_t>(method));
     return false;
     break;
  }

  return true;
}
{{end}}


{{range $serviceIndex, $service := $.Services}}
{{range $i, $method := $service.Methods}}
void MapRpcServiceBase::{{$method.MethodName}}(
  core::RpcRequestContext ctx, 
  const {{$.ContainerNamespace}}{{$method.InputType}}& request)
{
  ctx.ResponseErr(ErrorCode::kECGeneralInvalidRequest);
}
{{end}}
{{end}}

`

const RpcServerTemplateDefine = `#pragma once

#include <cstdint>
{{range $serviceIndex, $service := $.Services}}
enum {{$service.MethodsEnumName}}:uint32_t 
{
{{range $i, $method := $service.Methods}} k{{$service.Name}}{{$method.EnumName}} = {{$method.MethodID}},
{{end}}
};
{{end}}
`
