package caption_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCaption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Caption Suite")
}
