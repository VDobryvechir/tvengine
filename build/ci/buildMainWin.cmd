Xcopy /E /I /Y F:\cloud\go\root\src\github.com\Dobryvechir\microcore\pkg C:\Users\volod\go\pkg\mod\github.com\!dobryvechir\microcore@v1.0.0\pkg 
Xcopy /E /I /Y F:\cloud\go\root\src\github.com\VDobryvechir\dvserver\pkg ..\..\pkg 
copy ..\..\test\main.go ..\..\pkg\main\main.go
go install netcracker.com/oeconf/server/pkg/main
@del %GOPATH%\bin\oeconf.exe
@ren %GOPATH%\bin\main.exe oeconf.exe
copy %GOPATH%\bin\oeconf.exe c:\prg\Go\bin\oeconf.exe
