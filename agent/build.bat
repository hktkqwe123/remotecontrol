SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o agent.efl main.go get_sys_info.go send_flag.go

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARH=amd64
go build -o agent.exe main.go get_sys_info.go send_flag.go



