@SET BINDRY=..\..\..\bin
xcopy /y "..\..\..\..\..\..\src\libzeroproto\proto\database\mall_config.proto"  "%~dp0"
xcopy /y "..\..\..\..\..\..\src\libzeroproto\proto\database\mall_lv_config.proto"  "%~dp0"
%BINDRY%\protoc -I="." --proto_path="../../public/protomsg" --go_out . mall_config.proto mall_lv_config.proto flash_reward_config.proto
