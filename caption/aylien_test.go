package caption

import (
	"crypto/sha1"
	"errors"
	"io/ioutil"

	textapi "github.com/AYLIEN/aylien_textapi_go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("InMemoryDatastore", func() {
	var (
		logger        *log.Logger
		g             *AylienGenerator
		mockSummarize SummarizeFunc
		url           string
		numCaptions   int
		captions      []string
		err           error
	)

	BeforeEach(func() {
		logger = log.New()
		logger.Out = ioutil.Discard
		url = "https://test-url.com"
	})

	JustBeforeEach(func() {
		captions, err = g.Create(url, numCaptions)
	})

	Context("with cache hit", func() {
		var (
			cachedCaptions []string
			called         bool
		)
		BeforeEach(func() {
			// verify that we didnt call our summarize function
			mockSummarize = func(req *textapi.SummarizeParams) (*textapi.SummarizeResponse, error) {
				called = true
				return nil, nil
			}
			g = NewAylienGenerator(logger, mockSummarize)
			h := sha1.New()
			h.Write([]byte(url))
			id := string(h.Sum(nil))
			cachedCaptions = []string{
				"test4",
				"test5",
				"test6",
			}
			g.cache[id] = cachedCaptions
		})

		It("should NOT return an error", func() {
			Expect(err).To(BeNil())
		})

		It("should return captions", func() {
			Expect(captions).To(Equal(cachedCaptions))
		})

		It("should NOT call summarize", func() {
			Expect(called).To(BeFalse())
		})
	})

	Context("with summarize error", func() {
		BeforeEach(func() {
			mockSummarize = func(req *textapi.SummarizeParams) (*textapi.SummarizeResponse, error) {
				Expect(req.URL).To(Equal(url))
				Expect(req.NumberOfSentences).To(Equal(numCaptions))
				return nil, errors.New("summarize error")
			}
			g = NewAylienGenerator(logger, mockSummarize)
		})

		It("should return an error", func() {
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("summarize error"))
		})

		It("should NOT return captions", func() {
			Expect(captions).To(BeNil())
		})
	})

	Context("with summarize success", func() {
		var sentences []string
		BeforeEach(func() {
			sentences = []string{
				"one",
				"two",
				"three",
			}
			mockSummarize = func(req *textapi.SummarizeParams) (*textapi.SummarizeResponse, error) {
				Expect(req.URL).To(Equal(url))
				Expect(req.NumberOfSentences).To(Equal(numCaptions))
				resp := &textapi.SummarizeResponse{
					Sentences: sentences,
				}
				return resp, nil
			}
			g = NewAylienGenerator(logger, mockSummarize)
		})

		It("should NOT return an error", func() {
			Expect(err).To(BeNil())
		})

		It("should return captions", func() {
			Expect(captions).To(Equal(sentences))
		})

		It("should put captions in the cache", func() {
			h := sha1.New()
			h.Write([]byte(url))
			id := string(h.Sum(nil))
			Expect(g.cache).To(HaveKeyWithValue(id, sentences))
		})
	})
})
