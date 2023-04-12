package scripts

import (
	"errors"
	"gomssbuilder/internal/config"
	"log"
	"os/exec"
	"strings"
)

var (
	ErrNoSuchBranch = errors.New("no such branch")
)

type GomssOp struct {
}

func (g *GomssOp) GetBranches() ([]string, error) {
	//fetch origin branches
	err := g.fetch()
	if err != nil {
		return nil, err
	}
	branches, err := g.branches()
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (g *GomssOp) Checkout(branchName string) (string, error) {
	//checkout & pull
	branches, err := g.GetBranches()
	if err != nil {
		return "", err
	}
	flag := false
	for _, branch := range branches {
		if branchName == branch {
			flag = true
			break
		}
	}
	if !flag {
		log.Printf("No such branch: %v", branchName)
		return "", ErrNoSuchBranch
	}

	return g.checkout(branchName)
}

func (g *GomssOp) Build(localZRTC bool) ([]string, error) {
	ans := make([]string, 0)
	output, err := execute("go", "mod", "tidy")
	ans = append(ans, output)
	if err != nil {
		return ans, err
	}

	if localZRTC {
		output2, err := execute("make", `tags=nolibopusfile zrtcoutside`)
		ans = append(ans, output2)
		if err != nil {
			return ans, err
		}
	} else {
		output2, err := execute("make", "tags=nolibopusfile")
		ans = append(ans, output2)
		if err != nil {
			return ans, err
		}
	}

	return ans, nil
}

func (g *GomssOp) Publish(version string) (string, error) {
	data, err := execute("sh", "publish_gomss.sh", "gomss", version)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (g *GomssOp) checkout(name string) (string, error) {
	resp, err := execute("git", "checkout", name)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (g *GomssOp) fetch() error {
	_, err := execute("git", "fetch", "--all")
	if err != nil {
		return err
	}
	return nil
}

func (g *GomssOp) branches() ([]string, error) {
	data, err := execute("git", "tag")
	if err != nil {
		return nil, err
	}

	branches := make([]string, 0)
	for _, raw := range strings.Split(string(data), "\n") {
		raw = strings.TrimSpace(raw)
		branch := strings.Split(raw, "->")[0]
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

func execute(name string, params ...string) (string, error) {
	log.Printf("run command: %v, %v", name, params)
	cmd := exec.Command(name, params...)
	cmd.Dir = config.GetConfig().GomssPath
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to call run command %v, params: %v, error: %v", name, params, err)
		return string(data), err
	}
	return string(data), nil
}
