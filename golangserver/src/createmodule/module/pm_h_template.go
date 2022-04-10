package module

const PmHTemplate = `#pragma once
#include "gamelogic/pm/pm_interface.h"
#include "player_{{$.GetName}}_module.h"
#include "toolkit/singletontemplate.h"


class Player{{$.GetClassName}}PMCommand:public BasePMCommandHandler<Player{{$.GetClassName}}Module>
{
public:
	Player{{$.GetClassName}}PMCommand();
	virtual ~Player{{$.GetClassName}}PMCommand();
	virtual bool Init() override;
	virtual void Release() override;
private:

};

using Player{{$.GetClassName}}PMCommandInst = ChaSingleton<Player{{$.GetClassName}}PMCommand>;

`
