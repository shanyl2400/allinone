package boltdb

import (
	"gomssbuilder/internal/config"
	"log"
	"sync"

	"go.etcd.io/bbolt"
)

type Client struct {
	db   *bbolt.DB
	path string
}

func (c *Client) Open() error {
	db, err := bbolt.Open(c.path, 0666, nil)
	if err != nil {
		log.Printf("open boltdb failed, err: %v", err)
		return err
	}
	c.db = db
	return nil
}

func (c *Client) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

var (
	_client     *Client
	_clientOnce sync.Once
)

func GetClient() *Client {
	_clientOnce.Do(func() {
		_client = &Client{
			path: config.GetConfig().BoltDBPath,
		}
	})

	return _client
}
