::删除指定文件夹下的服务器日志文件
echo off
set DIR="%1"
echo DIR=%DIR%

for /R %DIR% %%f in (log) do ( 
::echo %%f
del /F /Q %%f
)

for /R %DIR% %%f in (behaviorlog) do ( 
::echo %%f
del /F /Q %%f
)