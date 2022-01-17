package git

import (
	"fmt"
	"git.dreame.tech/robot_commercial/dkit/dcomm"
	"git.dreame.tech/robot_commercial/dkit/dlog"
	"net/http"
	"os/exec"
	"strings"
)

func NewHelper(address, token string) *Helper {
	return &Helper{address, token}
}

type Helper struct {
	Address string
	Token   string
}

func (helper *Helper) ListAllGroupProjects(group string) []Repo {
	result := make([]Repo, 0)
	if group == "" {
		dlog.Error("input group name is empty")
		return result
	}
	currentPage := 1
	for {
		req := fmt.Sprintf("%s/api/v4/projects?private_token=%s&simple=true&include_subgroups=true&per_page=100&page=%d", helper.Address, helper.Token, currentPage)
		//可能由于gitlab版本问题
		//req := fmt.Sprintf("%s/api/v4/projects?private_token=%s&simple=true&owned=true&per_page=100&page=%d", helper.Address, helper.Token, currentPage)
		res, err := dcomm.HttpGet(req)
		if err != nil {
			dlog.Error("Get all group project failed:", err.Error())
			return []Repo{}
		}
		if repos := makeRepoList(res.GetBody()); len(repos) > 0 {
			for _, repo := range repos {
				if strings.Contains(repo.HttpUrl, group) {
					result = append(result, repo)
				}
			}
			currentPage++
			continue
		}
		dlog.Infof("Get %d repos", len(result))
		return result
	}
}

func (helper *Helper) GroupClone(group, dir, branch string) {
	if group == "" {
		dlog.Error("can not create new branch,group is empty")
		return
	}
	if branch == "" {
		branch = "master"
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			if out, err := exec.Command("/usr/bin/git", "clone", r.SshUrl, fmt.Sprintf("%s/%s", dir, r.Name), "-b", branch).CombinedOutput(); err != nil {
				dlog.Errorf("git clone [%s] failed:%s", r.SshUrl, string(out))
				return
			}
			dlog.Infof("git clone [%s] to [%s] success", r.Name, fmt.Sprintf("%s/%s/%s", dir, group, r.Name))
		}
		f(repo)
	}
}

func (helper *Helper) GroupMigrateClone(group, dir string) {
	if group == "" {
		dlog.Error("can not create new branch,group is empty")
		return
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			if out, err := exec.Command("/usr/bin/git", "clone", r.SshUrl, fmt.Sprintf("%s/%s", dir, r.Name)).CombinedOutput(); err != nil {
				dlog.Errorf("git clone [%s] failed:%s", r.SshUrl, string(out))
				return
			}
			if out, err := exec.Command("rm", "-r", fmt.Sprintf("%s/%s/.git", dir, r.Name)).CombinedOutput(); err != nil {
				dlog.Errorf("remove .git file [%s] failed:%s", r.Name, string(out))
				return
			}
			if out, err := exec.Command("/usr/bin/git", "clone", r.SshUrl, fmt.Sprintf("%s/%s/.git", dir, r.Name), "--bare").CombinedOutput(); err != nil {
				dlog.Errorf("git branch [%s] failed:%s", r.SshUrl, string(out))
				return
			}
			dlog.Infof("git migrate clone [%s] to [%s] success", r.Name, fmt.Sprintf("%s/%s/%s", dir, group, r.Name))
		}
		f(repo)
	}
}

func (helper *Helper) GroupNew(group, branchFrom, branchNew string) {
	dlog.Infof("git new branch,group:%s ,branch from:%s ,branch new:%s", group, branchFrom, branchNew)
	if group == "" {
		dlog.Error("can not create new branch,group is empty")
		return
	}
	if branchFrom == "" || branchNew == "" {
		dlog.Error("can not create new branch,branch name is empty")
		return
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			//info := map[string]string{"ref": branchFrom, "branch": branchNew}
			req := fmt.Sprintf("%s/api/v4/projects/%d/repository/branches?private_token=%s&ref=%s&branch=%s", helper.Address, r.ID, helper.Token, branchFrom, branchNew)
			rep, err := dcomm.HttpPost(req, "")
			if err != nil {
				dlog.Error("[%s] create new branch failed %s", r.Name, err.Error())
				return
			}
			if rep != nil && rep.GetCode() != http.StatusCreated {
				dlog.Errorf("[%s] create new branch failed:[%d]", r.Name, rep.GetCode())
				dlog.Error(rep.GetBody())
				return
			}
			dlog.Infof("[%s] create  branch [%s] from [%s] success", r.Name, branchNew, branchFrom)
		}
		f(repo)
	}
}

func (helper *Helper) GroupLock(group, branch string) {
	if group == "" {
		dlog.Error("can not lock branch,group is empty")
		return
	}
	if branch == "" {
		dlog.Error("can not lock branch,branch name is empty")
		return
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			//info := map[string]string{"ref": branchFrom, "branch": branchNew}
			req := fmt.Sprintf("%s/api/v4/projects/%d/protected_branches?private_token=%s&name=%s&push_access_level=0&merge_access_level=40", helper.Address, r.ID, helper.Token, branch)
			rep, err := dcomm.HttpPost(req, "")
			if err != nil {
				dlog.Error(" [%s] lock branch failed:%s", r.Name, err.Error())
				return
			}
			if rep != nil && rep.GetCode() != http.StatusCreated {
				dlog.Errorf("[%s] lock branch failed:[%d]", r.Name)
				dlog.Error(rep.GetBody())
				return
			}
			dlog.Infof("[%s] lock branch [%s] success", r.Name, branch)
		}
		f(repo)
	}
}

func (helper *Helper) GroupUnLock(group, branch string) {
	if group == "" {
		dlog.Error("can not lock branch,group is empty")
		return
	}
	if branch == "" {
		dlog.Error("can not lock branch,branch name is empty")
		return
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			//info := map[string]string{"ref": branchFrom, "branch": branchNew}
			req := fmt.Sprintf("%s/api/v4/projects/%d/protected_branches/%s?private_token=%s", helper.Address, r.ID, branch, helper.Token)
			rep, err := dcomm.HttpDelete(req)
			if err != nil {
				dlog.Error(" [%s] unlock branch failed:%s", r.Name, err.Error())
				return
			}
			if rep != nil && rep.GetCode() != http.StatusNoContent {
				dlog.Errorf("[%s] unlock branch failed:[%d]", r.Name, rep.GetCode())
				dlog.Error(rep.GetBody())
				return
			}
			dlog.Infof("[%s] unlock branch [%s] success", r.Name, branch)
		}
		f(repo)
	}
}

func (helper *Helper) GroupDeleteBranch(group, branch string) {
	if group == "" {
		dlog.Error("can not delete branch,group is empty")
		return
	}
	if branch == "" {
		dlog.Error("can not delete branch,branch name is empty")
		return
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			req := fmt.Sprintf("%s/api/v4/projects/%d/repository/branches/%s?private_token=%s", helper.Address, r.ID, branch, helper.Token)
			rep, err := dcomm.HttpDelete(req)
			if err != nil {
				dlog.Errorf("[%s] delete branch [%s] failed %s", r.Name, branch, err.Error())
				return
			}
			if rep != nil && rep.GetCode() != http.StatusNoContent {
				dlog.Errorf("[%s] delete branch [%s] failed:[%d]", r.Name, branch, rep.GetCode())
				dlog.Error(rep.GetBody())
				return
			}
			dlog.Infof("[%s] delete  branch [%s] success", r.Name, branch)
		}
		f(repo)
	}
}

func (helper *Helper) GroupMergeBranch(group, branchFrom, branchTo string) {
	if group == "" {
		dlog.Error("can not merge branch,group is empty")
		return
	}
	if branchFrom == "" || branchTo == "" {
		dlog.Error("can not merge branch,branch name is empty")
		return
	}
	//	wg := sync.WaitGroup{}
	//wg.Add()
	for _, repo := range helper.ListAllGroupProjects(group) {
		f := func(r Repo) {
			title := fmt.Sprintf("merge %s branch %s to %s", r.Name, branchFrom, branchTo)
			userID := 226 //shendongyu
			req := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests?private_token=%s&title=%s&source_branch=%s&target_branch=%s&assignee_id=%d", helper.Address, r.ID, helper.Token, title, branchFrom, branchTo, userID)
			rep, err := dcomm.HttpPost(req, "")
			if err != nil {
				dlog.Errorf("[%s] merge branch [%s->%s] failed: %s", r.Name, branchFrom, branchTo, err.Error())
				return
			}
			if rep != nil && rep.GetCode() != http.StatusNoContent {
				dlog.Errorf("[%s] merge branch [%s->%s] failed:[%d]", r.Name, branchFrom, branchTo, rep.GetCode())
				dlog.Error(rep.GetBody())
				return
			}
			dlog.Infof("[%s] merge branch [%s->%s] failed success", r.Name, branchFrom, branchTo)
		}
		f(repo)
	}
}
