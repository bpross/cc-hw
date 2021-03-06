package memory

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
		logger   *log.Logger
		p        *Poster
		mockDs   *mock_datastore.MockDatastore
		mockCtrl *gomock.Controller

		customerID string
		post       *dao.Post
		postID     bson.ObjectId
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger = log.New()
		logger.Out = ioutil.Discard
		mockDs = mock_datastore.NewMockDatastore(mockCtrl)
		p = NewPoster(logger, mockDs)

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

		Context("with datastore error", func() {
			var (
				dsErr error
			)
			BeforeEach(func() {
				dsErr = errors.New("test-error")
				mockDs.EXPECT().Insert(customerID, post).Return(nil, dsErr)
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
				mockDs.EXPECT().Insert(customerID, post).Return(post, nil)
			})

			It("should NOT return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return a post", func() {
				Expect(retPost).To(Equal(post))
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

		Context("with datastore error", func() {
			var (
				dsErr error
			)
			BeforeEach(func() {
				dsErr = errors.New("test-error")
				mockDs.EXPECT().Get(customerID, postID).Return(nil, dsErr)
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
				mockDs.EXPECT().Get(customerID, postID).Return(post, nil)
			})

			It("should NOT return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should return a post", func() {
				Expect(retPost).To(Equal(post))
			})
		})
	})

	Describe("Update", func() {
		var (
			retPost *dao.Post
			err     error
		)

		JustBeforeEach(func() {
			retPost, err = p.Update(customerID, post)
		})

		Context("with datastore error", func() {
			var (
				dsErr error
			)
			BeforeEach(func() {
				dsErr = errors.New("test-error")
				mockDs.EXPECT().Update(customerID, post).Return(nil, dsErr)
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
				mockDs.EXPECT().Update(customerID, post).Return(post, nil)
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
