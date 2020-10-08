
Rem  Generate code
SET GOARCH=amd64
SET GOOS=windows
set GOTOOLDIR=c:\go\pkg\tool\windows_amd64
pushd ..\..\
go generate ./...
if %errorlevel% neq 0 exit /b %errorlevel%
popd

echo finished