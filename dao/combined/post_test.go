package combined

import (
	"errors"
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	mock_datastore "github.com/bpross/cc-hw/mocks/github.com/bpross/cc-hw/datastore"
)

var _ = Describe("Poster", func() {
	var (
		logger         *log.Logger
		p              *Poster
		mockPersistent *mock_datastore.MockDatastore
		mockCache      *mock_datastore.MockDatastore
		mockCtrl       *gomock.Controller

		customerID string
		post       *dao.Post
		postID     bson.ObjectId
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger = log.New()
		logger.Out = ioutil.Discard
		mockPersistent = mock_datastore.NewMockDatastore(mockCtrl)
		mockCache = mock_datastore.NewMockDatastore(mockCtrl)
		p = NewPoster(logger, mockCache, mockPersistent)

		customerID = "test-customer"
		postID = bson.NewObjectId()
		post = &dao.Post{
			ID:     &postID,
			CustID: customerID,
			URL:    "test-url",
			Captions: []string{
				"caption1",
				"caption2",
				"caption3",
			},
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Insert", func() {
		var (
			retPost *dao.Post
			err     error
		)

		JustBeforeEach(func() {
			retPost, err = p.Insert(customerID, post)
		})

		Context("with persistent datastore error", func() {
			var (
				dsErr error
			)
			BeforeEach(func() {
				dsErr = errors.New("test-error")
				mockPersistent.EXPECT().Insert(customerID, post).Return(nil, dsErr)
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(dsErr))
			})

			It("should NOT return a post", func() {
				Expect(retPost).To(BeNil())
			})
		})

		Context("without datastore error", func() {
			BeforeEach(func() {
				mockPersistent.EXPECT().Insert(customerID, post).Return(post, nil)
			})

			Context("with cache error", func() {
				var dsErr error
				BeforeEach(func() {
					dsErr = errors.New("test-error")
					mockCache.EXPECT().Insert(customerID, post).Return(nil, dsErr)
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a post", func() {
					Expect(retPost).To(Equal(post))
				})
			})

			Context("without cache error", func() {
				BeforeEach(func() {
					mockCache.EXPECT().Insert(customerID, post).Return(post, nil)
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a post", func() {
					Expect(retPost).To(Equal(post))
				})
			})
		})
	})

	Describe("Get", func() {
		var (
			retPost *dao.Post
			err     error
		)

		JustBeforeEach(func() {
			retPost, err = p.Get(customerID, postID)
		})

		Context("with cache datastore success", func() {
			BeforeEach(func() {
				mockCache.EXPECT().Get(customerID, postID).Return(post, nil)
			})

			It("should NOT return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return a post", func() {
				Expect(retPost).To(Equal(post))
			})
		})

		// TODO turn the two contexts following this into a table driven test
		Context("with cache datastore error", func() {
			var cacheErr error
			BeforeEach(func() {
				cacheErr = errors.New("test-error")
				mockCache.EXPECT().Get(customerID, postID).Return(nil, cacheErr)
			})

			Context("with persistent datastore error", func() {
				var (
					dsErr error
				)
				BeforeEach(func() {
					dsErr = errors.New("test-error")
					mockPersistent.EXPECT().Get(customerID, postID).Return(nil, dsErr)
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(dsErr))
				})

				It("should NOT return a post", func() {
					Expect(retPost).To(BeNil())
				})
			})

			Context("without persistent datastore error", func() {
				BeforeEach(func() {
					mockPersistent.EXPECT().Get(customerID, postID).Return(post, nil)
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a post", func() {
					Expect(retPost).To(Equal(post))
				})
			})
		})

		Context("with cache datastore miss", func() {
			BeforeEach(func() {
				mockCache.EXPECT().Get(customerID, postID).Return(nil, nil)
			})

			Context("with persistent datastore error", func() {
				var (
					dsErr error
				)
				BeforeEach(func() {
					dsErr = errors.New("test-error")
					mockPersistent.EXPECT().Get(customerID, postID).Return(nil, dsErr)
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(dsErr))
				})

				It("should NOT return a post", func() {
					Expect(retPost).To(BeNil())
				})
			})

			Context("without persistent datastore error", func() {
				BeforeEach(func() {
					mockPersistent.EXPECT().Get(customerID, postID).Return(post, nil)
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a post", func() {
					Expect(retPost).To(Equal(post))
				})
			})
		})
	})

	Describe("Update", func() {
		var (
			retPost *dao.Post
			err     error
			dsErr   error
		)

		JustBeforeEach(func() {
			retPost, err = p.Update(customerID, post)
		})

		Context("with datastore error", func() {
			BeforeEach(func() {
				dsErr = errors.New("test-error")
				mockPersistent.EXPECT().Update(customerID, post).Return(nil, dsErr)
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(dsErr))
			})

			It("should NOT return a post", func() {
				Expect(retPost).To(BeNil())
			})
		})

		Context("without datastore error", func() {
			BeforeEach(func() {
				mockPersistent.EXPECT().Update(customerID, post).Return(post, nil)
			})

			Context("with cache update error", func() {
				var (
					cacheUpdateErr, cacheDeleteErr error
				)
				BeforeEach(func() {
					cacheUpdateErr = errors.New("update-error")
					mockCache.EXPECT().Update(customerID, post).Return(nil, cacheUpdateErr)
				})

				Context("with cache delete error", func() {
					BeforeEach(func() {
						cacheDeleteErr = errors.New("update-error")
						mockCache.EXPECT().Delete(customerID, *post.ID).Return(cacheDeleteErr)
					})
					It("should return an error", func() {
						Expect(err).NotTo(BeNil())
						Expect(err).To(Equal(cacheDeleteErr))
					})

					It("should NOT return a post", func() {
						Expect(retPost).To(BeNil())
					})
				})

				Context("With cache delete sucecss", func() {
					BeforeEach(func() {
						mockCache.EXPECT().Delete(customerID, *post.ID).Return(nil)
					})
					It("should NOT return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should return a post", func() {
						Expect(retPost).To(Equal(post))
					})
				})
			})

			Context("with cache update success", func() {
				BeforeEach(func() {
					mockCache.EXPECT().Update(customerID, post).Return(post, nil)
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a post", func() {
					Expect(retPost).To(Equal(post))
				})
			})
		})
	})
})
