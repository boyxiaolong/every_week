package module

const MsgHTemplate = `#pragma once
#include "protocol/base_message_handler.h"
#include "toolkit/singletontemplate.h"


class Player{{$.GetClassName}}MsgHandler:public BaseMessageHandler
{
public:
  virtual bool Init(BaseServer* server) override;
};

using Player{{$.GetClassName}}MsgHandlerInst = ChaSingleton<Player{{$.GetClassName}}MsgHandler>;
`
