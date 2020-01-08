package handler

import (
	"io/ioutil"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

func setupRouter(p Poster) *gin.Engine {
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	r.GET("/post/:id", p.Get)
	r.POST("/post", p.Post)
	r.PUT("/post/:id", p.Put)
	return r
}
