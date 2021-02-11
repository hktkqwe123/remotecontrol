GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy.elf main.go sqlit3_handle.go
GO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o proxy.exe main.go sqlit3_handle.go
