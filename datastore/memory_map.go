package datastore

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
)

// Datastore provides an interface for inserting, retrieving and updating information about posts
type Datastore interface {
	Insert(string, *dao.Post) (*dao.Post, error)
	Get(string, bson.ObjectId) (*dao.Post, error)
	Update(string, *dao.Post) (*dao.Post, error)
	Delete(string, bson.ObjectId) error
}

// InMemoryDatastore implements the Datastore interface for in memory storage
type InMemoryDatastore struct {
	logger *log.Logger
	store  map[string]*dao.Post
}

// NewInMemoryDatastore creates a new InMemoryDatastore with the provided options
func NewInMemoryDatastore(logger *log.Logger) *InMemoryDatastore {
	m := make(map[string]*dao.Post)
	return &InMemoryDatastore{
		logger: logger,
		store:  m,
	}
}

// Insert inserts a new post into the map, customerID is used to enforce tenancy
func (d *InMemoryDatastore) Insert(customerID string, post *dao.Post) (*dao.Post, error) {
	if post == nil {
		return nil, NewInvalidArugmentError("must provide post")
	}

	if post.ID != nil {
		return nil, NewInvalidArugmentError("cannot provide ID")
	}

	if customerID == "" {
		return nil, NewInvalidArugmentError("customerID")
	}

	logger := d.logger.WithFields(log.Fields{
		"customerID": customerID,
		"url":        post.URL,
	})

	logger.Info("inserting into memory map")

	// Generate ID
	id := bson.NewObjectId()

	// Create new post
	r := &dao.Post{
		ID:       &id,
		CustID:   customerID,
		URL:      post.URL,
		Captions: post.Captions,
	}

	// Create composite ID to enforce tenancy
	storeID := createCompositeID(customerID, id)

	// Store post
	d.store[storeID] = r

	logger.WithFields(log.Fields{
		"post_id": id.Hex(),
	}).Debug("successfully inserted post")

	return r, nil
}

// Get retrieves the postID from the map, tenancy is enforced with the customerID
func (d *InMemoryDatastore) Get(customerID string, postID bson.ObjectId) (*dao.Post, error) {
	if postID == "" {
		return nil, NewInvalidArugmentError("postID")
	}

	if customerID == "" {
		return nil, NewInvalidArugmentError("customerID")
	}

	// Create composite id
	storeID := createCompositeID(customerID, postID)

	logger := d.logger.WithFields(log.Fields{
		"customerID": customerID,
		"postID":     postID.Hex(),
	})

	logger.Info("retrieving from memory map")

	// Find post in the datastore, if ok is false, the post DNE
	r, ok := d.store[storeID]
	if !ok {
		return nil, NewNotFoundError("post")
	}

	logger.Debug("successfully retrieved")
	return r, nil
}

// Update stores the given post in the map
func (d *InMemoryDatastore) Update(customerID string, post *dao.Post) (*dao.Post, error) {
	if post == nil {
		return nil, fmt.Errorf("must provide post")
	}

	if post.ID == nil {
		return nil, NewInvalidArugmentError("postID")
	}

	if customerID == "" {
		return nil, NewInvalidArugmentError("customerID")
	}

	logger := d.logger.WithFields(log.Fields{
		"customerID": customerID,
		"postID":     post.ID.Hex(),
	})

	logger.Info("updating in memory map")

	// Create composite id
	storeID := createCompositeID(customerID, *post.ID)

	// Find post in the datastore, if ok is false, the post DNE
	prev, ok := d.store[storeID]
	if !ok {
		return nil, NewNotFoundError("post")
	}

	// Only copy over captions
	prev.Captions = post.Captions

	// Store post
	d.store[storeID] = prev

	logger.Debug("successfully updated post")
	return prev, nil
}

// Delete is not implemented for the in memory datastore
func (d *InMemoryDatastore) Delete(customerID string, postID bson.ObjectId) error {
	return fmt.Errorf("unimplemented")
}

func createCompositeID(customerID string, postID bson.ObjectId) string {
	return fmt.Sprintf("%s:%s", customerID, postID.Hex())
}
