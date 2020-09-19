SET GOARCH=amd64
SET GOOS=windows
set CGO_ENABLED=0
go run -tags win main_win.go

