@del %GOPATH%\bin\tvengine.exe
@Xcopy /E /I /Y C:\prg\libgo\microcore\pkg %GOPATH%\pkg\mod\github.com\!dobryvechir\microcore@v1.0.0\pkg 
go install github.com/VDobryvechir/tvengine/pkg/main
@ren %GOPATH%\bin\main.exe tvengine.exe
copy %GOPATH%\bin\tvengine.exe c:\prg\tvengine\tvengine.exe
