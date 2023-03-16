package boltdb

import (
	"encoding/json"
	"gomssbuilder/internal/repository/model"
	"time"

	"go.etcd.io/bbolt"
	"gopkg.in/mgo.v2/bson"
)

const (
	publishRecordBucket = "publish_record"
)

type BoltdbPublishRecordRepository struct {
	client *Client
	bucket *bbolt.Bucket
}

func (b *BoltdbPublishRecordRepository) Put(p *model.PublishRecord) error {
	return b.client.db.Update(func(tx *bbolt.Tx) error {
		p.ID = bson.NewObjectId().Hex()
		p.RecordedAt = time.Now()
		data, _ := json.Marshal(p)
		b := tx.Bucket([]byte(publishRecordBucket))
		return b.Put([]byte(p.ID), data)
	})
}
func (b *BoltdbPublishRecordRepository) Recent(size int) ([]*model.PublishRecord, error) {
	output := make([]*model.PublishRecord, 0)
	b.client.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(publishRecordBucket))

		c := b.Cursor()
		i := 0
		for k, v := c.Last(); k != nil && i < size; k, v = c.Prev() {
			record := new(model.PublishRecord)
			err := json.Unmarshal(v, record)
			if err != nil {
				return err
			}
			output = append(output, record)
			i++
		}
		return nil
	})

	return output, nil
}

func NewBoltdbPublishRecordRepository() (*BoltdbPublishRecordRepository, error) {
	var bucket *bbolt.Bucket
	err := GetClient().db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(publishRecordBucket))
		if err != nil {
			return err
		}

		bucket = b
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BoltdbPublishRecordRepository{
		bucket: bucket,
		client: GetClient(),
	}, nil
}
