@del %GOPATH%\bin\tvengine.exe
@Xcopy /E /I /Y C:\cloud\microcore\pkg %GOPATH%\pkg\mod\github.com\!dobryvechir\microcore@v1.0.0\pkg 
go install github.com/VDobryvechir/tvengine/pkg/main
@ren %GOPATH%\bin\main.exe tvengine.exe
