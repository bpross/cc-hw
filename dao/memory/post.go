package memory

import (
	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/datastore"
)

// MemoryPoster implements the Poster interface using an in-memory datastore
type MemoryPoster struct {
	logger *log.Logger
	ds     datastore.Datastore
}

// NewMemoryPoster creates a new MemoryPoster with the supplied options
func NewMemoryPoster(logger *log.Logger, ds datastore.Datastore) *MemoryPoster {
	return &MemoryPoster{
		logger: logger,
		ds:     ds,
	}
}

func (d *MemoryPoster) Insert(customerID string, post *dao.Post) (*dao.Post, error) {
	return d.ds.Insert(customerID, post)
}

func (d *MemoryPoster) Get(customerID string, postID bson.ObjectId) (*dao.Post, error) {
	return d.ds.Get(customerID, postID)
}

func (d *MemoryPoster) Update(customerID string, postID *dao.Post) (*dao.Post, error) {
	return d.ds.Update(customerID, postID)
}
