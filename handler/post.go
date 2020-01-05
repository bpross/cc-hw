package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/datastore"
)

const customerIDHeader = "x-customer-id"

type PostRequest struct {
	URL      string   `json:"url"`
	Captions []string `json:"captions", omitempty`
}

type PutRequest struct {
	Captions []string `json:"captions", omitempty`
}

// Poster defines the interface to handle post requests
type Poster interface {
	Get(*gin.Context)
	Post(*gin.Context)
	Put(*gin.Context)
}

// DefaultPoster implements the Poster interface
type DefaultPoster struct {
	ds dao.Poster
}

// NewDefaultPoster returns a DefaultPoster with the provided options
func NewDefaultPoster(ds dao.Poster) *DefaultPoster {
	return &DefaultPoster{
		ds: ds,
	}
}

// Get defines the handler for post GET requests
func (p *DefaultPoster) Get(c *gin.Context) {
	urlID := c.Param("id")
	// Check if id is valid
	ok := validateID(c, urlID)
	if !ok {
		return
	}

	id := bson.ObjectIdHex(urlID)

	// Get headers
	customerID := getAndValidateHeaders(c)
	if customerID == "" {
		return
	}

	post, err := p.ds.Get(customerID, id)
	if err != nil {
		setReturnError(err, c)
		return
	}
	c.PureJSON(http.StatusOK, post)
	return
}

// Post defines the handler for post POST requests
func (p *DefaultPoster) Post(c *gin.Context) {
	// Get headers
	customerID := getAndValidateHeaders(c)
	if customerID == "" {
		return
	}
	// Hydrate post
	req := &PostRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input := postRequestToPost(*req)
	post, err := p.ds.Insert(customerID, input)
	if err != nil {
		setReturnError(err, c)
		return
	}
	c.PureJSON(http.StatusOK, post)
	return
}

// Put defines the handler for handling post PUT requests
func (p *DefaultPoster) Put(c *gin.Context) {
	urlID := c.Param("id")
	// Check if id is valid
	ok := validateID(c, urlID)
	if !ok {
		return
	}

	id := bson.ObjectIdHex(urlID)

	// Get headers
	customerID := getAndValidateHeaders(c)
	if customerID == "" {
		return
	}

	// Hydrate post
	req := &PutRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	input := putRequestToPost(*req, id)
	post, err := p.ds.Update(customerID, input)
	if err != nil {
		setReturnError(err, c)
		return
	}
	c.PureJSON(http.StatusOK, post)
	return
}

func validateID(c *gin.Context, urlID string) bool {
	ok := bson.IsObjectIdHex(urlID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid post id"})
		return false
	}
	return true
}

func getAndValidateHeaders(c *gin.Context) string {
	// Get headers
	customerID := c.Request.Header.Get(customerIDHeader)

	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "must include customerID in headers"})
	}
	return customerID
}

func setReturnError(dsErr error, c *gin.Context) {
	switch dsErr.(type) {
	case *datastore.InvalidArugment:
		c.JSON(http.StatusBadRequest, gin.H{"message": dsErr.Error()})
		return
	case *datastore.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": dsErr.Error()})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": dsErr.Error()})
		return
	}
}

func postRequestToPost(req PostRequest) *dao.Post {
	return &dao.Post{
		URL:      req.URL,
		Captions: req.Captions,
	}
}

func putRequestToPost(req PutRequest, id bson.ObjectId) *dao.Post {
	return &dao.Post{
		ID:       &id,
		Captions: req.Captions,
	}
}
