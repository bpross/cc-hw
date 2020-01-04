package datastore

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"
)

// Record stores in the information about a url
type Record struct {
	ID       bson.ObjectId
	CustID   string
	URL      string
	Captions []string
}

// Datastore provides an interface for inserting, retrieving and updating information about records
type Datastore interface {
	Insert(string, *Record) (*Record, error)
	Get(string, bson.ObjectId) (*Record, error)
	Update(string, *Record) (*Record, error)
}

// InMemoryDatastore implements the Datastore interface for in memory storage
type InMemoryDatastore struct {
	logger log.Logger
	store  map[string]*Record
}

// NewInMemoryDatastore creates a new InMemoryDatastore with the provided options
func NewInMemoryDatastore(logger log.Logger) *InMemoryDatastore {
	m := make(map[string]*Record)
	return &InMemoryDatastore{
		logger: logger,
		store:  m,
	}
}

// Insert inserts a new record into the map, customerID is used to enforce tenancy
func (d *InMemoryDatastore) Insert(customerID string, record *Record) (*Record, error) {
	if record == nil {
		return nil, fmt.Errorf("must provide record")
	}

	if record.ID != "" {
		return nil, fmt.Errorf("cannot provide ID")
	}

	if customerID == "" {
		return nil, NewInvalidArugmentError("customerID")
	}

	logger := d.logger.WithFields(log.Fields{
		"customerID": customerID,
		"url":        record.URL,
	})

	logger.Debug("inserting")

	// Generate ID
	id := bson.NewObjectId()

	// Create new record
	r := &Record{
		ID:       id,
		CustID:   customerID,
		URL:      record.URL,
		Captions: record.Captions,
	}

	// Create composite ID to enforce tenancy
	storeID := createCompositeID(customerID, id)

	// Store record
	d.store[storeID] = r

	logger.WithFields(log.Fields{
		"record_id": id.Hex(),
	}).Info("successfully inserted record")

	return r, nil
}

// Get retrieves the recordID from the map, tenancy is enforced with the customerID
func (d *InMemoryDatastore) Get(customerID string, recordID bson.ObjectId) (*Record, error) {
	if recordID == "" {
		return nil, NewInvalidArugmentError("recordID")
	}

	if customerID == "" {
		return nil, NewInvalidArugmentError("customerID")
	}

	// Create composite id
	storeID := createCompositeID(customerID, recordID)

	logger := d.logger.WithFields(log.Fields{
		"customerID": customerID,
		"recordID":   recordID.Hex(),
	})

	logger.Debug("retrieving")

	// Find record in the datastore, if ok is false, the record DNE
	r, ok := d.store[storeID]
	if !ok {
		return nil, NewNotFoundError("record")
	}

	logger.Info("successfully retrieved")
	return r, nil
}

// Update stores the given record in the map
func (d *InMemoryDatastore) Update(customerID string, record *Record) (*Record, error) {
	if record == nil {
		return nil, fmt.Errorf("must provide record")
	}

	if record.ID == "" {
		return nil, NewInvalidArugmentError("recordID")
	}

	if customerID == "" {
		return nil, NewInvalidArugmentError("customerID")
	}

	logger := d.logger.WithFields(log.Fields{
		"customerID": customerID,
		"recordID":   record.ID.Hex(),
	})

	logger.Debug("updating")

	// Create composite id
	storeID := createCompositeID(customerID, record.ID)

	// Find record in the datastore, if ok is false, the record DNE
	_, ok := d.store[storeID]
	if !ok {
		return nil, NewNotFoundError("record")
	}

	// Store record
	d.store[storeID] = record

	logger.Info("successfully updated record")
	return record, nil
}

func createCompositeID(customerID string, recordID bson.ObjectId) string {
	return fmt.Sprintf("%s:%s", customerID, recordID.Hex())
}
