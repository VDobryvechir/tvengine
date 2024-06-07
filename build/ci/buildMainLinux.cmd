del %GOPATH%\bin\linux_amd64\main
SET CGO_ENABLED=0
SET GOOS=linux
go install github.com/VDobryvechir/tvengine/pkg/main
@copy %GOPATH%\bin\linux_amd64\main c:\prg\tvengine\tvengine 
@del %GOPATH%\bin\linux_amd64\main

