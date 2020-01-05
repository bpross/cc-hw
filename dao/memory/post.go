package memory

import (
	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/datastore"
)

// Poster implements the Poster interface using an in-memory datastore
type Poster struct {
	logger *log.Logger
	ds     datastore.Datastore
}

// NewPoster creates a new Poster with the supplied options
func NewPoster(logger *log.Logger, ds datastore.Datastore) *Poster {
	return &Poster{
		logger: logger,
		ds:     ds,
	}
}

// Insert handles post insert requests using the underlying in memory datastore
func (d *Poster) Insert(customerID string, post *dao.Post) (*dao.Post, error) {
	return d.ds.Insert(customerID, post)
}

// Get handles post get requests using the underlying in memory datastore
func (d *Poster) Get(customerID string, postID bson.ObjectId) (*dao.Post, error) {
	return d.ds.Get(customerID, postID)
}

// Update handles post update requests using the underlying in memory datastore
func (d *Poster) Update(customerID string, postID *dao.Post) (*dao.Post, error) {
	return d.ds.Update(customerID, postID)
}
