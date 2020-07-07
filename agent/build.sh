CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o agent.elf main.go get_sys_info.go send_flag.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o agent.exe main.go get_sys_info.go send_flag.go
