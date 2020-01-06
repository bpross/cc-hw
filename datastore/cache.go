package datastore

import (
	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
)

// NoOpCache implements the Datastore interface and is a placeholder for a read/write thru cache
type NoOpCache struct {
	logger *log.Logger
}

// NewNoOpCache creates a new NoOpCache with the provided options
func NewNoOpCache(logger *log.Logger) *NoOpCache {
	return &NoOpCache{
		logger: logger,
	}
}

// Insert just logs that insert was called
func (c *NoOpCache) Insert(customerID string, post *dao.Post) (*dao.Post, error) {
	c.logger.Info("calling cache insert")
	return nil, nil
}

// Get just logs that get was called
func (c *NoOpCache) Get(customerID string, postID bson.ObjectId) (*dao.Post, error) {
	c.logger.Info("calling cache get")
	return nil, nil
}

// Update just logs that update was called
func (c *NoOpCache) Update(customerID string, post *dao.Post) (*dao.Post, error) {
	c.logger.Info("calling cache update")
	return nil, nil
}

// Delete just logs that delete was called
func (c *NoOpCache) Delete(customerID string, postID bson.ObjectId) error {
	c.logger.Info("calling cache delete")
	return nil
}
