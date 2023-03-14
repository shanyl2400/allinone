package service

import (
	"errors"
	"gomssbuilder/internal/scripts"
	"sync"
)

var (
	ErrNoSuchBranch        = errors.New("no such branch")
	ErrNoSuchZRTC          = errors.New("no such zrtc")
	ErrSelectBranchAndZrtc = errors.New("hasn't select gomss branch and zrtc")
)

type GomssBuilder struct {
	zrtc  *scripts.ZrtcOp
	gomss *scripts.GomssOp

	sync.Mutex
}

func (gb *GomssBuilder) ListGomssBranches() ([]string, error) {
	//ListGomssBranches与publish有冲突
	gb.Lock()
	defer gb.Unlock()

	return gb.gomss.GetBranches()
}

func (gb *GomssBuilder) ListZrtc() ([]*scripts.ZRTC, error) {
	return gb.zrtc.ListZrtcs()
}

func (gb *GomssBuilder) Publish(gomssBranch, zrtcPath, version string) error {
	//publish必须是原子操作
	gb.Lock()
	defer gb.Unlock()

	if zrtcPath == "" || gomssBranch == "" {
		return ErrSelectBranchAndZrtc
	}

	// 1.checkout branch
	err := gb.gomss.Checkout(gomssBranch)
	if err != nil {
		return err
	}

	// 2.download zrtc
	_, err = gb.zrtc.DownloadZrtc(zrtcPath)
	if err != nil {
		return err
	}

	// 3.go get & make
	_, err = gb.gomss.Build()
	if err != nil {
		return err
	}

	// 4.publish & chmod
	_, err = gb.gomss.Publish(version)
	if err != nil {
		return err
	}

	return nil
}

func NewBuilder() *GomssBuilder {
	return &GomssBuilder{
		zrtc:  new(scripts.ZrtcOp),
		gomss: new(scripts.GomssOp),
	}
}
