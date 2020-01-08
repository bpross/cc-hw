package combined

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/datastore"
)

// Poster implements the Poster interface using a read/write thru cache backed by a persistent datastore
type Poster struct {
	logger     *log.Logger
	cache      datastore.Datastore
	persistent datastore.Datastore
}

// NewPoster creates a new Poster with the supplied options
func NewPoster(logger *log.Logger, cache, persistent datastore.Datastore) *Poster {
	return &Poster{
		logger:     logger,
		cache:      cache,
		persistent: persistent,
	}
}

// Insert calls both the persistent and cache datastores. It will only return error
// on persistent failure. cache failure just means a read to the persistent store later
func (d *Poster) Insert(customerID string, post *dao.Post) (*dao.Post, error) {
	d.logger.Info("inserting")
	dsPost, err := d.persistent.Insert(customerID, post)
	if err != nil {
		d.logger.Warn("failed to insert into persistent")
		return nil, err
	}

	logger := d.logger.WithFields(log.Fields{
		"post_id": dsPost.ID.Hex(),
	})

	_, err = d.cache.Insert(customerID, dsPost)
	if err != nil {
		logger.Warn("failed to insert into cache")
	}

	logger.Debug("successfully inserted")
	return dsPost, nil
}

// Get tries the cache first and then the persistent store, on any cache error
// the code will try to read from the persistent storage
func (d *Poster) Get(customerID string, postID bson.ObjectId) (*dao.Post, error) {
	logger := d.logger.WithFields(log.Fields{
		"post_id": postID.Hex(),
	})

	logger.Info("retrieving")
	post, err := d.cache.Get(customerID, postID)
	if err != nil || post == nil {
		if err == nil {
			err = errors.New("cache miss")
		}
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Info("failed to retrieve from cache")

		post, err = d.persistent.Get(customerID, postID)
		if err != nil {
			logger.WithFields(log.Fields{
				"error": err.Error(),
			}).Warn("failed to retrieve from persistent")
			return nil, err
		}
	}

	logger.Debug("successfully retrieved")
	return post, nil
}

// Update calls the persistent store first. On success, the cache is called. If the
// cache call fails, the value will be deleted from the cache
func (d *Poster) Update(customerID string, post *dao.Post) (*dao.Post, error) {
	logger := d.logger.WithFields(log.Fields{
		"post_id": post.ID.Hex(),
	})

	logger.Info("updating")
	post, err := d.persistent.Update(customerID, post)
	if err != nil {
		return nil, err
	}

	_, err = d.cache.Update(customerID, post)
	if err != nil {
		logger.Warn("failed to update into cache")

		err = d.cache.Delete(customerID, *post.ID)
		if err != nil {
			// This is bad, is there a better way to handle this?
			logger.Error("failed to delete into cache")
			return nil, err
		}
	}

	logger.Debug("successfully updated")
	return post, nil
}
