package combined_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCombined(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dao Combined Suite")
}
