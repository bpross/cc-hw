package main

import (
	"time"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bpross/cc-hw/dao/combined"
	"github.com/bpross/cc-hw/datastore"
	"github.com/bpross/cc-hw/handler"
)

func main() {
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

	// Setup handler and routes
	handler := handler.NewDefaultPoster(combinedPoster)
	r.GET("/post/:id", handler.Get)
	r.POST("/post", handler.Post)
	r.PUT("/post/:id", handler.Put)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
