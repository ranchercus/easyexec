GOOS="linux" CGO_ENABLED="0" GOARCH=amd64 go build -o bin/easyexec_linux_amd64 -ldflags "-s -w" main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/easyexec_win_x64.exe -ldflags "-s -w" main.go
go build -o bin/easyexec_darwin -ldflags "-s -w" main.go