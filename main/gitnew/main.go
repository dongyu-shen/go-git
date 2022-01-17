package main

import (
	"flag"

	"git.dreame.tech/robot_commercial/dkit/dlog"
	"git.dreame.tech/robot_commercial/dkit/dtool"
	"git.dreame.tech/robot_commercial/platform/common/example/git/git"
)

//go:generate GOARCH=amd64 GOOS=linux go build -o gitnew main.go
func main() {
	dlog.OpenLog("git")
	address := flag.String("address", "http://192.168.10.10", "gitlab address,eg:http://example.com")
	group := flag.String("g", "robot_commercial/platform_android", "group name")
	branchFrom := flag.String("f", "v3.0.0", "branch name to create branch from")
	branchNew := flag.String("n", "v3.0.1", "new branch name")
	token := flag.String("t", "", "token")
	flag.Parse()
	helper := git.NewHelper(*address, *token)
	helper.GroupNew(*group, *branchFrom, *branchNew)
	dtool.Ok()
}
