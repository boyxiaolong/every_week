package module

const ModuleCppTemplate= `#include "stdafx.h"
#include "player_{{$.GetName}}_module.h"
#include "error_code.pb.h"
#include "game_server.h"
#include "gamelogic/player_module_type.h"
#include "gsdefine.h"
#include "msg_type.pb.h"
#include "toolkit/string_parse.h"
#include "player_{{$.GetName}}_pm_command.h"
#include "{{$.GetName}}.pb.h"


Player{{$.GetClassName}}Module::Player{{$.GetClassName}}Module(Player& player)
	:PlayerModuleInterface(player, PlayerModuleType::kPlayerModule{{$.GetClassName}})
{
}

Player{{$.GetClassName}}Module::~Player{{$.GetClassName}}Module()
{

}

bool Player{{$.GetClassName}}Module::OnInit()
{
	return true;
}

void Player{{$.GetClassName}}Module::OnRelease()
{

}

void Player{{$.GetClassName}}Module::OnRun()
{

}

async::Future<bool> Player{{$.GetClassName}}Module::OnLoadData()
{
	co_return true;
}

void Player{{$.GetClassName}}Module::OnNewPlayer()
{

}

void Player{{$.GetClassName}}Module::OnDataReady()
{

}

void Player{{$.GetClassName}}Module::OnPlayerOnline()
{

}

void Player{{$.GetClassName}}Module::OnPlayerOffline()
{

}

async::Future<bool> Player{{$.GetClassName}}Module::OnBackup()
{
	co_return true;
}

void Player{{$.GetClassName}}Module::OnMessage(uint16_t msg_type,const google::protobuf::Message& message)
{
	switch (msg_type)
    {
    default:
      LOGERR("Unhandle module mesage module: {}, type : {}",
      static_cast<int>(GetModuleType()), msg_type);
    break;
    }
}

void Player{{$.GetClassName}}Module::OnModuleEvent(const IModuleEvent& base_event)
{

}

int Player{{$.GetClassName}}Module::OnPMCommand(const StringParse& parser)
{
	return Player{{$.GetClassName}}PMCommandInst::Singleton()->HandleCommand(
		GetPlayer(), *this, parser);
}


`
