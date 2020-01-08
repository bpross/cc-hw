package dao

import (
	"labix.org/v2/mgo/bson"
)

// Post stores in the information about a url
type Post struct {
	ID       *bson.ObjectId `json:"id, omitempty"`
	CustID   string         `json:"-"` // do not return when we marshal to json
	URL      string         `json:"url"`
	Captions []string       `json:"captions", omitempty`
}

// Poster defines the interface for persisting posts
type Poster interface {
	Insert(string, *Post) (*Post, error)
	Get(string, bson.ObjectId) (*Post, error)
	Update(string, *Post) (*Post, error)
}
