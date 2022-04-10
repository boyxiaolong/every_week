package module

const ModuleHTemplate = `#pragma once
#include "gamelogic/player_module_interface.h"
#include "player_{{$.GetName}}_module_interface.h"


class Player{{$.GetClassName}}Module:public PlayerModuleInterface
					,public Player{{$.GetClassName}}ModuleInterface
{
public:
	explicit Player{{$.GetClassName}}Module(Player& player);
	virtual ~Player{{$.GetClassName}}Module();

private:
	virtual bool OnInit() override;
	virtual void OnRelease() override;
	virtual void OnRun() override;

	virtual async::Future<bool> OnLoadData() override;
	virtual void OnNewPlayer() override;
	virtual void OnDataReady() override;

	virtual void OnPlayerOnline() override;
	virtual void OnPlayerOffline() override;
	virtual async::Future<bool> OnBackup() override;

	virtual void OnMessage(uint16_t msg_type, const google::protobuf::Message& message) override;
	virtual void OnModuleEvent(const IModuleEvent& base_event) override;
	virtual int OnPMCommand(const StringParse& parser) override;
public:

private:

};

`
