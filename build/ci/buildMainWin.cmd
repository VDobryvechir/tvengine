Xcopy /E /I /Y C:\cloud\go\root\src\github.com\Dobryvechir\microcore\pkg C:\Users\volod\go\pkg\mod\github.com\!dobryvechir\microcore@v1.0.0\pkg 
Xcopy /E /I /Y C:\cloud\go\root\src\github.com\VDobryvechir\dvserver\pkg ..\..\pkg 
copy ..\..\test\main.go ..\..\pkg\main\main.go
go install github.com/VDobryvechir/tvengine/pkg/main
@del %GOPATH%\bin\tvengine.exe
@ren %GOPATH%\bin\main.exe tvengine.exe
copy %GOPATH%\bin\tvengine.exe c:\prg\Go\bin\tvengine.exe
