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
	mock_dao "github.com/bpross/cc-hw/mocks/github.com/bpross/cc-hw/dao"
)

var _ = Describe("DefaulPoster", func() {
	var (
		mockCtrl   *gomock.Controller
		mockPoster *mock_dao.MockPoster
		handler    *DefaultPoster
		router     *gin.Engine
		customerID string
		recorder   *httptest.ResponseRecorder
		method     string
		url        string
		req        *http.Request
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockPoster = mock_dao.NewMockPoster(mockCtrl)
		handler = NewDefaultPoster(mockPoster)
		router = setupRouter(handler)
		customerID = "test-customer"
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Get", func() {
		var (
			postID bson.ObjectId
			err    error
		)

		BeforeEach(func() {
			postID = bson.NewObjectId()
			method = "GET"
		})

		JustBeforeEach(func() {
			router.ServeHTTP(recorder, req)
		})

		Context("with invalid id", func() {
			BeforeEach(func() {
				url = "/post/blah"
				req, err = http.NewRequest(method, url, nil)
				Expect(err).To(BeNil())
			})

			It("should return StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return a useful message", func() {
				expected := `{"message":"invalid post id"}`
				actual := strings.TrimSuffix(recorder.Body.String(), "\n")
				Expect(actual).To(Equal(expected))
			})
		})

		Context("with valid id", func() {
			BeforeEach(func() {
				url = "/post/" + postID.Hex()
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
				BeforeEach(func() {
					req, err = http.NewRequest(method, url, nil)
					Expect(err).To(BeNil())
					req.Header.Add(customerIDHeader, customerID)
				})

				Context("with datastore error", func() {
					Context("with InvalidArugment error", func() {
						BeforeEach(func() {
							daoErr := datastore.NewInvalidArugmentError("test-error")
							mockPoster.EXPECT().Get(customerID, postID).Return(nil, daoErr)
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

					Context("with NotFound error", func() {
						BeforeEach(func() {
							daoErr := datastore.NewNotFoundError("test-error")
							mockPoster.EXPECT().Get(customerID, postID).Return(nil, daoErr)
						})

						It("should return StatusNotFound", func() {
							Expect(recorder.Code).To(Equal(http.StatusNotFound))
						})

						It("should return the datastore message", func() {
							expected := `{"message":"test-error not found"}`
							actual := strings.TrimSuffix(recorder.Body.String(), "\n")
							Expect(actual).To(Equal(expected))
						})
					})

					Context("with unknown error", func() {
						BeforeEach(func() {
							daoErr := errors.New("test-error")
							mockPoster.EXPECT().Get(customerID, postID).Return(nil, daoErr)
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
					)
					BeforeEach(func() {
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
						mockPoster.EXPECT().Get(customerID, postID).Return(dsPost, nil)
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
						Captions: []string{
							"caption1",
							"caption2",
							"caption3",
						},
					}
					body, err = json.Marshal(post)
					Expect(err).To(BeNil())
					req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
					Expect(err).To(BeNil())
					req.Header.Add(customerIDHeader, customerID)
					req.Header.Add("Content-Type", "application/json")
				})

				Context("with datastore error", func() {
					Context("with InvalidArugment error", func() {
						BeforeEach(func() {
							daoErr := datastore.NewInvalidArugmentError("test-error")
							mockPoster.EXPECT().Insert(customerID, &post).Return(nil, daoErr)
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
							mockPoster.EXPECT().Insert(customerID, &post).Return(nil, daoErr)
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
						mockPoster.EXPECT().Insert(customerID, &post).Return(dsPost, nil)
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

	Describe("Put", func() {
		var (
			postID bson.ObjectId
			err    error
		)

		BeforeEach(func() {
			postID = bson.NewObjectId()
			method = "PUT"
		})

		JustBeforeEach(func() {
			router.ServeHTTP(recorder, req)
		})
		Context("with invalid id", func() {
			BeforeEach(func() {
				url = "/post/blah"
				req, err = http.NewRequest(method, url, nil)
				Expect(err).To(BeNil())
			})

			It("should return StatusBadRequest", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return a useful message", func() {
				expected := `{"message":"invalid post id"}`
				actual := strings.TrimSuffix(recorder.Body.String(), "\n")
				Expect(actual).To(Equal(expected))
			})
		})

		Context("with valid id", func() {
			BeforeEach(func() {
				url = "/post/" + postID.Hex()
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
							ID: &postID,
							Captions: []string{
								"caption1",
								"caption2",
								"caption3",
							},
						}
						body, err = json.Marshal(post)
						Expect(err).To(BeNil())
						req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
						Expect(err).To(BeNil())
						req.Header.Add(customerIDHeader, customerID)
						req.Header.Add("Content-Type", "application/json")
					})

					Context("with datastore error", func() {
						Context("with InvalidArugment error", func() {
							BeforeEach(func() {
								daoErr := datastore.NewInvalidArugmentError("test-error")
								mockPoster.EXPECT().Update(customerID, &post).Return(nil, daoErr)
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
								mockPoster.EXPECT().Update(customerID, &post).Return(nil, daoErr)
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
							mockPoster.EXPECT().Update(customerID, &post).Return(dsPost, nil)
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
