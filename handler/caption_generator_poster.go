package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bpross/cc-hw/caption"
	"github.com/bpross/cc-hw/dao"
)

type generatePostRequest struct {
	URL string `json:"url"`
}

// CaptionGeneratorPoster implements the Poster interface and generates captions
// it takes a base Poster, because we only need to implement the Post method
type CaptionGeneratorPoster struct {
	Poster
	ds               dao.Poster
	captionGenerator caption.Generator
	numCaptions      int
}

// NewCaptionGeneratorPoster returns a CaptionGeneratorPoster with the provided options
func NewCaptionGeneratorPoster(base Poster, ds dao.Poster, g caption.Generator, numCaptions int) *CaptionGeneratorPoster {
	return &CaptionGeneratorPoster{
		base,
		ds,
		g,
		numCaptions,
	}
}

// Post defines the handler for POST requests. This generates captions and then saves
func (p *CaptionGeneratorPoster) Post(c *gin.Context) {
	// Get headers
	customerID := getAndValidateHeaders(c)
	if customerID == "" {
		return
	}

	// Hydrate post
	req := &generatePostRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Generate captions
	captions, err := p.captionGenerator.Create(req.URL, p.numCaptions)
	if err != nil {
		// log here
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to generate captions"})
		return
	}

	// Save post
	input := generatePostRequestToPost(*req, captions)
	post, err := p.ds.Insert(customerID, input)
	if err != nil {
		setReturnError(err, c)
		return
	}
	c.PureJSON(http.StatusOK, post)
	return
}

func generatePostRequestToPost(req generatePostRequest, captions []string) *dao.Post {
	return &dao.Post{
		URL:      req.URL,
		Captions: captions,
	}
}
