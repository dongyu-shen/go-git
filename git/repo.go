package git

import (
	"git.dreame.tech/robot_commercial/dkit/dencode/djson"
	"git.dreame.tech/robot_commercial/dkit/dlog"
)

type Repo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	SshUrl      string `json:"ssh_url_to_repo"`
	HttpUrl     string `json:"http_url_to_repo"`
}

func makeRepoList(data string) []Repo {
	var repos = make([]Repo, 0)
	if err := djson.BindFromString(data, &repos); err != nil {
		dlog.Error("make repo failed:", data)
	}
	return repos
}
