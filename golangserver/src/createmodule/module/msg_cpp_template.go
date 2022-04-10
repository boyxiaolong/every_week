package module

const MsgCppTemplate = `#include "stdafx.h"
#include "player_{{$.GetName}}_msg_handle.h"
#include "gamelogic/module_event/module_event_define.h"
#include "gamelogic/module_event/module_event_mgr.h"
#include "gamelogic/player_manager.h"
#include "gamelogic/player_module_type.h"
#include "server_type.h"
#include "msg_type.pb.h"
#include "{{$.GetName}}.pb.h"

using namespace protomsg;

bool Player{{$.GetClassName}}MsgHandler::Init(BaseServer* server)
{
//module_event
//REG_MODULE_EVENT(GSModuleEventType::kGSModuleEventBuildingUpgraded,
//PlayerModuleType::kPlayerModule{{$.GetClassName}});

//client message
//REG_PLAYER_MSG(MsgCL2GS{{$.GetClassName}}InfoRequest,
//PlayerModuleType::kPlayerModule{{$.GetClassName}});
return true;
}
`
