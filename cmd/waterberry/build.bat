SET GOARM=6
SET GOARCH=arm
SET GOOS=linux
set CGO_ENABLED=0
go build -tags linux -o bin\waterberry  main.go

scp  bin\waterberry pi@192.168.4.245:~/waterberry