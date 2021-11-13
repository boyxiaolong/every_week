::功能: 编译服务器代码
@echo off
SETLOCAL

set src_dir=%1/server.sln
if "%src_dir%"=="" goto usage

set is_rebuild=%2
if "%is_rebuild%"=="" goto usage

set build_type=%3
if "%build_type%"=="" goto usage

CALL "C:\Program Files (x86)\Microsoft Visual Studio\2019\Community\Common7\Tools\VsDevCmd.bat"

echo src_dir,%src_dir%,is_rebuild,%is_rebuild%,build_type,%build_type%

set compile_project_list=libzeroproto libcommon libregion mapserver gameserver
(for %%a in (%compile_project_list%) do ( 
   echo begin compile %%a
   devenv %src_dir% /%is_rebuild%  "%build_type%|x64" /project %%a
   if %errorlevel% NEQ 0 (
		echo compile %%a error
		exit
   )
))

echo build finish
exit /B 0

:usage
@echo Usage: %0 ^<src_dir_sln^> ^<Build/Rebuild^> ^<Release/Debug^>
exit /B 0

ENDLOCAL