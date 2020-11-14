for /f %%i in ('git rev-list -1 head') do set GIT_COMMIT=%%i
for /f %%i in ('date /t') do set BUILD_DATE=%%i
for /f %%i in ('time /t') do set BUILD_TIME=%%i
echo %GIT_COMMIT%
echo %BUILD_TIME%

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
set LD_FLAGS="-X main.GitCommit='%GIT_COMMIT%' -X main.BuildTime='%BUILD_DATE%-%BUILD_TIME%'"
go build -ldflags=%LD_FLAGS% -tags linux -o bin\waterberry main.go 
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
