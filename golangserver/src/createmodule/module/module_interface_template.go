package module

const ModuleInterfaceTemplate = `#pragma once
class Player{{$.GetClassName}}ModuleInterface
{
public:
	Player{{$.GetClassName}}ModuleInterface() {};
	virtual ~Player{{$.GetClassName}}ModuleInterface() = default;
public:

};
`