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

func (g *GomssOp) Build(gomssBranch string, localZRTC bool) ([]string, error) {
	ans := make([]string, 0)
	output, err := execute("go", "mod", "tidy")
	ans = append(ans, output)
	if err != nil {
		return ans, err
	}

	var output2 string
	switch gomssBranch[:2] {
	case "v2":
		output2, err = g.buildV2()
	case "v3":
		output2, err = g.buildV3(localZRTC)
	case "v4":
		output2, err = g.buildV4(localZRTC)
	default:
		output2, err = g.buildV2()
	}
	if output2 != "" {
		ans = append(ans, output2)
	}
	if err != nil {
		return ans, err
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

func (g *GomssOp) buildV3(localZRTC bool) (string, error) {
	var output string
	var err error
	if localZRTC {
		output, err = execute("make", `tags=nolibopusfile zrtcoutside`)
		if err != nil {
			return output, err
		}
	} else {
		output, err = execute("make", `tags=nolibopusfile`)
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

func (g *GomssOp) buildV4(localZRTC bool) (string, error) {
	var output string
	var err error
	if localZRTC {
		output, err = execute("make", `tags=noffmpeg,nozrtc`)
		if err != nil {
			return output, err
		}
	} else {
		output, err = execute("make", "tags=noffmpeg")
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

func (g *GomssOp) buildV2() (string, error) {
	output, err := execute("make")
	if err != nil {
		return output, err
	}
	return output, nil
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
