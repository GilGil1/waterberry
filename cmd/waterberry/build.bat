
Rem  Generate code
SET GOARCH=amd64
SET GOOS=windows
set GOTOOLDIR=c:\go\pkg\tool\windows_amd64
pushd ..\..\
go generate ./...
if %errorlevel% neq 0 exit /b %errorlevel%
popd

rem compile
SET GOARM=6
SET GOARCH=arm
SET GOOS=linux
set CGO_ENABLED=0
set HOST=pi@192.168.4.245
del  bin\waterberry

go build -tags linux -o bin\waterberry main.go
if %errorlevel% neq 0 exit /b %errorlevel%

REM Stop and upgrade the software
ssh  %HOST% sudo service waterd stop
scp  bin\waterberry %HOST%:~/waterberry
scp  waterd.service %HOST%:~/waterberry
scp  config.json %HOST%:~/waterberry

ssh  %HOST% sudo cp ~/waterberry/waterd.service /etc/systemd/system/waterd.service
ssh  %HOST% sudo systemctl enable waterd.service


ssh  %HOST% sudo systemctl daemon-reload
ssh  %HOST% sudo service waterd start
