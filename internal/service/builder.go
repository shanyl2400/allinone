package service

import (
	"errors"
	"gomssbuilder/internal/repository"
	"gomssbuilder/internal/repository/model"
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
	scriptsMessage []string
}

func (gb *GomssBuilder) ListGomssBranches() ([]string, error) {
	//ListGomssBranches与publish有冲突
	gb.Lock()
	defer gb.Unlock()

	return gb.gomss.GetBranches()
}

func (gb *GomssBuilder) ListZrtc() ([]*ZRTC, error) {
	ans, err := gb.zrtc.ListZrtcs()
	if err != nil {
		return nil, err
	}
	output := make([]*ZRTC, 0, len(ans))
	for _, z := range ans {
		output = append(output, &ZRTC{
			Name: z.Name,
			Path: z.Path,
		})
	}
	return output, nil
}

func (gb *GomssBuilder) Publish(gomssBranch, zrtcPath, version string, localZRTC bool) error {
	//publish必须是原子操作
	gb.Lock()
	defer gb.Unlock()

	gb.cleanPublishLog()

	if zrtcPath == "" || gomssBranch == "" {
		return ErrSelectBranchAndZrtc
	}

	// 1.checkout branch
	gb.pushPublishLog("正在切换分支...")
	msg, err := gb.gomss.Checkout(gomssBranch)
	gb.pushPublishLog(msg)
	if err != nil {
		return err
	}

	// 2.download zrtc
	if !localZRTC {
		gb.pushPublishLog("正在下载zrtc...")
		msg, err := gb.zrtc.DownloadZrtc(zrtcPath)
		gb.pushPublishLog(msg...)
		if err != nil {
			return err
		}
	}

	// 3.go get & make
	gb.pushPublishLog("正在编译gomss...")
	msgs, err := gb.gomss.Build(localZRTC)
	gb.pushPublishLog(msgs...)
	if err != nil {
		return err
	}

	// 4.publish & chmod
	gb.pushPublishLog("正在发布镜像...")
	msg, err = gb.gomss.Publish(version)
	gb.pushPublishLog(msg)
	if err != nil {
		return err
	}

	zrtcVersion := zrtcPath
	if localZRTC {
		zrtcVersion = "zrtc-outer"
	}
	err = repository.GetPublishRecordRepository().Put(&model.PublishRecord{
		Version:     version,
		GomssBranch: gomssBranch,
		ZrtcVersion: zrtcVersion,
	})
	if err != nil {
		return err
	}

	return nil
}

func (gb *GomssBuilder) PublishLogs() []string {
	if gb.scriptsMessage == nil {
		return make([]string, 0)
	}
	return gb.scriptsMessage
}

func (gb *GomssBuilder) RecentPublish() ([]*PublishRecord, error) {
	records, err := repository.GetPublishRecordRepository().Recent(10)
	if err != nil {
		return nil, err
	}
	output := make([]*PublishRecord, 0, len(records))
	for _, r := range records {
		output = append(output, &PublishRecord{
			ID:          r.ID,
			Version:     r.Version,
			GomssBranch: r.GomssBranch,
			ZrtcVersion: r.ZrtcVersion,
			RecordedAt:  r.RecordedAt,
		})
	}
	return output, nil
}

func (gb *GomssBuilder) pushPublishLog(logs ...string) {
	gb.scriptsMessage = append(gb.scriptsMessage, logs...)
}

func (gb *GomssBuilder) cleanPublishLog() {
	gb.scriptsMessage = make([]string, 0)
}

func NewBuilder() *GomssBuilder {
	return &GomssBuilder{
		zrtc:  new(scripts.ZrtcOp),
		gomss: new(scripts.GomssOp),
	}
}
