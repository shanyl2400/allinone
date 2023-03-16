package repository

import (
	"gomssbuilder/internal/repository/boltdb"
	"gomssbuilder/internal/repository/model"
	"sync"
)

type PublishRecordRepository interface {
	Put(p *model.PublishRecord) error
	Recent(size int) ([]*model.PublishRecord, error)
}

var (
	_publishRecordRepo     PublishRecordRepository
	_publishRecordRepoOnce sync.Once
)

func GetPublishRecordRepository() PublishRecordRepository {
	var err error
	_publishRecordRepoOnce.Do(func() {
		_publishRecordRepo, err = boltdb.NewBoltdbPublishRecordRepository()
		if err != nil {
			panic(err)
		}
	})
	return _publishRecordRepo
}
