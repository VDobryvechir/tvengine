@del %GOPATH%\bin\tvengine.exe
go install github.com/VDobryvechir/tvengine/pkg/main
@ren %GOPATH%\bin\main.exe tvengine.exe
copy %GOPATH%\bin\tvengine.exe c:\prg\tvengine\tvengine.exe
