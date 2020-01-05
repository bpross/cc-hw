package dao

import (
	"labix.org/v2/mgo/bson"
)

// Post stores in the information about a url
type Post struct {
	ID       bson.ObjectId
	CustID   string
	URL      string
	Captions []string
}

// Poster defines the interface for persisting posts
type Poster interface {
	Insert(string, *Post) (*Post, error)
	Get(string, bson.ObjectId) (*Post, error)
	Update(string, *Post) (*Post, error)
}
