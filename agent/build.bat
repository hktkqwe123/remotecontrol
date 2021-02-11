SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-H windowsgui" -o agent.elf main.go

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARH=amd64
go build -ldflags "-H windowsgui" -o agent.exe main.go



