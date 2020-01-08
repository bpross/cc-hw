package caption

import (
	"crypto/sha1"

	textapi "github.com/AYLIEN/aylien_textapi_go"
	log "github.com/sirupsen/logrus"
)

// SummarizeFunc defines the function used by AylienGenerator to request captions
type SummarizeFunc func(*textapi.SummarizeParams) (*textapi.SummarizeResponse, error)

// AylienGenerator implements the generator interface and uses Aylien API to do so
type AylienGenerator struct {
	logger        *log.Logger
	summarizeFunc SummarizeFunc
	cache         map[string][]string // TODO make this access a datastore
}

// NewAylienGenerator creates an AylienGenerator with the provided options
func NewAylienGenerator(logger *log.Logger, summarizeFunc SummarizeFunc) *AylienGenerator {
	cache := make(map[string][]string)
	return &AylienGenerator{
		logger:        logger,
		summarizeFunc: summarizeFunc,
		cache:         cache,
	}
}

// Create sends a request to Aylien to Summarize the provided url
func (g *AylienGenerator) Create(url string, numCaptions int) ([]string, error) {
	logger := g.logger.WithFields(log.Fields{
		"url":         url,
		"numCaptions": string(numCaptions),
	})
	logger.Info("generating captions")
	// First check if we have already summarized this url
	h := sha1.New()
	h.Write([]byte(url))
	id := string(h.Sum(nil))
	if captions, ok := g.cache[id]; ok {
		logger.Debug("cache hit")
		return captions, nil
	}
	logger.Debug("cache miss")
	// Create request
	req := &textapi.SummarizeParams{
		URL:               url,
		NumberOfSentences: numCaptions,
	}

	// Send request
	logger.Debug("sending request")
	resp, err := g.summarizeFunc(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("request successfull, adding captions to cache")
	// insert into cache
	g.cache[id] = resp.Sentences
	return resp.Sentences, nil
}
