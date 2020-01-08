package main

import (
	"os"
	"strconv"
	"time"

	textapi "github.com/AYLIEN/aylien_textapi_go"
	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bpross/cc-hw/caption"
	"github.com/bpross/cc-hw/dao/combined"
	"github.com/bpross/cc-hw/datastore"
	"github.com/bpross/cc-hw/handler"
)

const (
	envApiKey       = "AYLIEN_API_KEY"
	envAppID        = "AYLIEN_APP_ID"
	envCaptionCount = "AYLIEN_CAPTION_COUNT"
)

func main() {
	// Load env variables panic if any of these are not set
	var (
		apiKey, appID, captionCountEnv string
		present                        bool
	)
	apiKey, present = os.LookupEnv(envApiKey)
	if !present {
		panic("AYLIEN_API_KEY must be set in env")
	}
	appID, present = os.LookupEnv(envAppID)
	if !present {
		panic("AYLIEN_APP_ID must be set in env")
	}
	captionCountEnv, present = os.LookupEnv(envCaptionCount)
	if !present {
		panic("AYLIEN_CAPTION_COUNT must be set in env")
	}
	captionCount, err := strconv.Atoi(captionCountEnv)
	if err != nil {
		panic(err.Error())
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	r := gin.New()        // don't use the Default(), since it comes with a logger
	r.Use(gin.Recovery()) // add Recovery middleware

	// Pulled from the gin-logrus docs
	useBanner := false
	useUTC := true
	logger := logrus.StandardLogger()
	r.Use(ginlogrus.WithTracing(
		logger,
		useBanner,
		time.RFC3339,
		useUTC,
		"",
		[]byte{}, // where jaeger might have put the trace id
		[]byte{}, // where the trace ID might already be populated in the headers
		ginlogrus.WithAggregateLogging(true)))

	// Setup datastores
	memDS := datastore.NewInMemoryDatastore(logger)
	cacheDS := datastore.NewNoOpCache(logger)

	// Setup DAO
	combinedPoster := combined.NewPoster(logger, cacheDS, memDS)

	// Setup generator
	auth := textapi.Auth{appID, apiKey}
	client, err := textapi.NewClient(auth, true)
	if err != nil {
		panic(err)
	}
	captionGenerator := caption.NewAylienGenerator(logger, client.Summarize)

	// Setup handler and routes
	baseHandler := handler.NewDefaultPoster(combinedPoster)
	generateHandler := handler.NewCaptionGeneratorPoster(baseHandler, combinedPoster, captionGenerator, captionCount)
	r.GET("/post/:id", generateHandler.Get)
	r.POST("/post", generateHandler.Post)
	r.PUT("/post/:id", generateHandler.Put)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
