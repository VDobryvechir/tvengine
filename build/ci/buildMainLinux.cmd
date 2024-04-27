Xcopy /E /I /Y F:\cloud\go\root\src\github.com\Dobryvechir\microcore\pkg C:\Users\volod\go\pkg\mod\github.com\!dobryvechir\microcore@v1.0.0\pkg 
Xcopy /E /I /Y F:\cloud\go\root\src\github.com\VDobryvechir\dvserver\pkg ..\..\pkg 
copy ..\..\test\main.go ..\..\pkg\main\main.go
del %GOPATH%\bin\linux_amd64\main
call copyAllConfigs.cmd
SET CGO_ENABLED=0
SET GOOS=linux
go install netcracker.com/oeconf/server/pkg/main
@copy %GOPATH%\bin\linux_amd64\main F:\a\nec\order-capture-config\configs\oeconf\oeconf 
@del %GOPATH%\bin\linux_amd64\main

