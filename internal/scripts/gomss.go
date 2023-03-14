package scripts

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var (
	gomssDir = "/home/shanyonglong/gomss"

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

func (g *GomssOp) Checkout(branchName string) error {
	//checkout & pull
	branches, err := g.GetBranches()
	if err != nil {
		return err
	}
	flag := false
	for _, branch := range branches {
		if branchName == branch {
			flag = true
			break
		}
	}
	if !flag {
		log.Fatalf("No such branch: %v", branchName)
		return ErrNoSuchBranch
	}

	return g.checkout(branchName)
}

func (g *GomssOp) Build() (map[string]string, error) {
	ans := make(map[string]string)
	output, err := execute("go", "mod", "tidy")
	ans["get"] = output
	if err != nil {
		return ans, err
	}

	output2, err := execute("make")
	ans["make"] = output2
	if err != nil {
		return ans, err
	}
	return ans, nil
}

func (g *GomssOp) Publish(version string) (string, error) {
	data, err := execute("sh", "publish2.sh", "gomss", version)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (g *GomssOp) checkout(name string) error {
	_, err := execute("git", "checkout", name)
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
		return err
	}
	return nil
}

func (g *GomssOp) fetch() error {
	_, err := execute("git", "fetch")
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
		return err
	}
	return nil
}

func (g *GomssOp) branches() ([]string, error) {
	data, err := execute("git", "branch", "-r")
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
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
	cmd := exec.Command(name, params...)
	cmd.Dir = gomssDir
	data, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(data))
		log.Fatalf("failed to call run command %v, error: %v", name, err)
		return string(data), err
	}
	return string(data), nil
}
