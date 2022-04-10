go build -o consoleclient.exe main.go
xcopy /y consoleclient.exe "..\..\..\..\..\bin\Release\golangclient\consoleclient\"
pause