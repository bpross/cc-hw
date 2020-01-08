package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/datastore"
	mock_caption "github.com/bpross/cc-hw/mocks/github.com/bpross/cc-hw/caption"
	mock_dao "github.com/bpross/cc-hw/mocks/github.com/bpross/cc-hw/dao"
)

var _ = Describe("CaptionGeneratorPoster", func() {
	var (
		mockCtrl      *gomock.Controller
		mockPoster    *mock_dao.MockPoster
		baseHandler   *DefaultPoster
		handler       *CaptionGeneratorPoster
		mockGenerator *mock_caption.MockGenerator
		router        *gin.Engine
		customerID    string
		recorder      *httptest.ResponseRecorder
		method        string
		url           string
		req           *http.Request
		numCaptions   int
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockPoster = mock_dao.NewMockPoster(mockCtrl)
		mockGenerator = mock_caption.NewMockGenerator(mockCtrl)
		baseHandler = NewDefaultPoster(mockPoster)
		numCaptions = 3
		handler = NewCaptionGeneratorPoster(baseHandler, mockPoster, mockGenerator, numCaptions)
		router = setupRouter(handler)
		customerID = "test-customer"
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Post", func() {
		var err error
		BeforeEach(func() {
			method = "POST"
			url = "/post"
		})

		JustBeforeEach(func() {
			router.ServeHTTP(recorder, req)
		})

		Context("without customerID in header", func() {
			BeforeEach(func() {
				req, err = http.NewRequest(method, url, nil)
				Expect(err).To(BeNil())
			})

			It("should return StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return a useful message", func() {
				expected := `{"message":"must include customerID in headers"}`
				actual := strings.TrimSuffix(recorder.Body.String(), "\n")
				Expect(actual).To(Equal(expected))
			})
		})

		Context("with customerID in header", func() {
			Context("with json error", func() {
				BeforeEach(func() {
					req, err = http.NewRequest(method, url, nil)
					Expect(err).To(BeNil())
					req.Header.Add(customerIDHeader, customerID)
				})

				It("should return StatusBadRequest", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("should return a useful message", func() {
					expected := `{"message":"invalid request"}`
					actual := strings.TrimSuffix(recorder.Body.String(), "\n")
					Expect(actual).To(Equal(expected))
				})
			})

			Context("with valid json", func() {
				var (
					post dao.Post
					body []byte
				)
				BeforeEach(func() {
					post = dao.Post{
						URL: "test-url",
					}
					body, err = json.Marshal(post)
					Expect(err).To(BeNil())
					req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
					Expect(err).To(BeNil())
					req.Header.Add(customerIDHeader, customerID)
					req.Header.Add("Content-Type", "application/json")
				})

				Context("with generator error", func() {
					var genErr error
					BeforeEach(func() {
						genErr = errors.New("generator error")
						mockGenerator.EXPECT().Create(post.URL, numCaptions).Return(nil, genErr)
					})

					It("should return StatusInternalServerError", func() {
						Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
					})

					It("should return a useful message", func() {
						expected := `{"message":"unable to generate captions"}`
						actual := strings.TrimSuffix(recorder.Body.String(), "\n")
						Expect(actual).To(Equal(expected))
					})
				})

				Context("with generator success", func() {
					var (
						captions     []string
						generatePost dao.Post
					)
					BeforeEach(func() {
						captions = []string{
							"caption1",
							"caption2",
							"caption3",
						}
						generatePost = dao.Post{
							URL:      "test-url",
							Captions: captions,
						}
						mockGenerator.EXPECT().Create(post.URL, numCaptions).Return(captions, nil)
					})
					Context("with datastore error", func() {
						Context("with InvalidArugment error", func() {
							BeforeEach(func() {
								daoErr := datastore.NewInvalidArugmentError("test-error")
								mockPoster.EXPECT().Insert(customerID, &generatePost).Return(nil, daoErr)
							})

							It("should return StatusBadRequest", func() {
								Expect(recorder.Code).To(Equal(http.StatusBadRequest))
							})

							It("should return the datastore message", func() {
								expected := `{"message":"invalid test-error"}`
								actual := strings.TrimSuffix(recorder.Body.String(), "\n")
								Expect(actual).To(Equal(expected))
							})
						})

						Context("with unknown error", func() {
							BeforeEach(func() {
								daoErr := errors.New("test-error")
								mockPoster.EXPECT().Insert(customerID, &generatePost).Return(nil, daoErr)
							})

							It("should return StatusInternalServerError", func() {
								Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
							})

							It("should return the datastore message", func() {
								expected := `{"message":"test-error"}`
								actual := strings.TrimSuffix(recorder.Body.String(), "\n")
								Expect(actual).To(Equal(expected))
							})
						})
					})

					Context("with datastore success", func() {
						var (
							dsPost *dao.Post
							postID bson.ObjectId
						)
						BeforeEach(func() {
							postID = bson.NewObjectId()
							dsPost = &dao.Post{
								ID:     &postID,
								CustID: customerID,
								URL:    "test-url",
								Captions: []string{
									"caption1",
									"caption2",
									"caption3",
								},
							}
							mockPoster.EXPECT().Insert(customerID, &generatePost).Return(dsPost, nil)
						})

						It("should return StatusOK", func() {
							Expect(recorder.Code).To(Equal(http.StatusOK))
						})

						It("should return a post", func() {
							expected := fmt.Sprintf(`{"id":"%s","url":"test-url","captions":["caption1","caption2","caption3"]}`, postID.Hex())
							actual := strings.TrimSuffix(recorder.Body.String(), "\n")
							Expect(actual).To(Equal(expected))
						})
					})
				})
			})
		})
	})
})
