package module

const PmCppTemplate = `#include "stdafx.h"
#include "player_{{$.GetName}}_pm_command.h"
#include "error_code.pb.h"
#include "gamelogic/pm/pm_command_manager.h"
#include "gamelogic/player_module_type.h"
#include "toolkit/string_parse.h"


Player{{$.GetClassName}}PMCommand::Player{{$.GetClassName}}PMCommand()
{

}

Player{{$.GetClassName}}PMCommand::~Player{{$.GetClassName}}PMCommand()
{

}

bool Player{{$.GetClassName}}PMCommand::Init()
{
	//REG_PM("rmhero", Player{{$.GetClassName}}PMCommand::RemoveHero, "DelHero",
		//"rmhero <hero id>", PlayerModuleType::kPlayerModule{{$.GetClassName}});
	return true;
}

void Player{{$.GetClassName}}PMCommand::Release()
{

}


`
