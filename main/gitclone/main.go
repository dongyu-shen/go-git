package main

import (
	"flag"
	"git.dreame.tech/robot_commercial/dkit/dlog"
	"git.dreame.tech/robot_commercial/dkit/dtool"
	"git.dreame.tech/robot_commercial/platform/common/example/git/git"
)

func main() {
	dlog.OpenLog("git")
	address := flag.String("address", "http://192.168.10.10", "gitlab address,eg:http://example.com")
	group := flag.String("g", "robot_commercial/platform_android", "group name")
	dir := flag.String("d", ".", "clone to dir")
	branch := flag.String("b", "master", "branch name")
	token := flag.String("t", "", "token")
	flag.Parse()
	helper := git.NewHelper(*address, *token)
	helper.GroupClone(*group, *dir, *branch)
	dtool.Ok()
}
