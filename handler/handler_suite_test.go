package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

func setupRouter(p *DefaultPoster) *gin.Engine {
	r := gin.Default()
	r.GET("/post/:id", p.Get)
	r.POST("/post", p.Post)
	r.PUT("/post/:id", p.Put)
	return r
}
