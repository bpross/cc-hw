package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/datastore"
)

const customerIdHeader = "x-customer-id"

// Poster defines the interface to handle post requests
type Poster interface {
	Get(*gin.Context)
	Post(*gin.Context)
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

func (p *DefaultPoster) Get(c *gin.Context) {
	urlID := c.Param("id")
	// Check if id is valid
	ok := bson.IsObjectIdHex(urlID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid post id"})
		return
	}

	id := bson.ObjectIdHex(urlID)

	// Get headers
	customerID := c.Request.Header.Get(customerIdHeader)

	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "must include customerID in headers"})
		return
	}

	post, err := p.ds.Get(customerID, id)
	if err != nil {
		switch err.(type) {
		case *datastore.InvalidArugment:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		case *datastore.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.PureJSON(http.StatusOK, post)
	return
}

func (p *DefaultPoster) Post(c *gin.Context) {
	// Get headers
	customerID := c.Request.Header.Get(customerIdHeader)

	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "must include customerID in headers"})
		return
	}

	// Hydrate post
	input := &dao.Post{}
	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	post, err := p.ds.Insert(customerID, input)
	if err != nil {
		switch err.(type) {
		case *datastore.InvalidArugment:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.PureJSON(http.StatusOK, post)
	return
}
