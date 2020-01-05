package datastore

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
)

var _ = Describe("InMemoryDatastore", func() {
	var (
		logger     *log.Logger
		ds         *InMemoryDatastore
		post       *dao.Post
		customerID string
	)

	BeforeEach(func() {
		// TODO pass in no-op logger, so we dont log during tests
		logger = log.New()
		ds = NewInMemoryDatastore(logger)
		customerID = "test-customer"
	})

	Describe("Insert", func() {
		var (
			retPost *dao.Post
			err     error
		)

		JustBeforeEach(func() {
			retPost, err = ds.Insert(customerID, post)
		})

		Context("without post", func() {
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("invalid must provide post"))
			})

			It("should NOT return a post", func() {
				Expect(retPost).To(BeNil())
			})
		})

		Context("with post.ID", func() {
			BeforeEach(func() {
				id := bson.NewObjectId()
				post = &dao.Post{
					ID: &id,
				}
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("invalid cannot provide ID"))
			})

			It("should NOT return a post", func() {
				Expect(retPost).To(BeNil())
			})
		})

		Context("without post.ID", func() {
			BeforeEach(func() {
				post = &dao.Post{
					URL: "test-url",
					Captions: []string{
						"caption1",
						"caption2",
						"caption3",
					},
				}
			})

			Context("without customerID", func() {
				BeforeEach(func() {
					customerID = ""
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(Equal("invalid customerID"))
				})

				It("should NOT return a post", func() {
					Expect(retPost).To(BeNil())
				})
			})

			Context("with customerID", func() {
				BeforeEach(func() {
					customerID = "test-customer"
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a post", func() {
					Expect(retPost.ID).NotTo(Equal(""))
					Expect(retPost.CustID).To(Equal(customerID))
					Expect(retPost.URL).To(Equal(post.URL))
					Expect(retPost.Captions).To(Equal(post.Captions))
				})

				It("should insert a post", func() {
					storeID := createCompositeID(customerID, *retPost.ID)
					Expect(ds.store).To(HaveKeyWithValue(storeID, retPost))
				})
			})
		})
	})

	Describe("Get", func() {
		var (
			retPost *dao.Post
			err     error
			postID  bson.ObjectId
		)

		JustBeforeEach(func() {
			retPost, err = ds.Get(customerID, postID)
		})

		Context("without postID", func() {
			BeforeEach(func() {
				postID = ""
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("invalid postID"))
			})

			It("should NOT return a post", func() {
				Expect(retPost).To(BeNil())
			})
		})

		Context("with postID", func() {
			BeforeEach(func() {
				postID = bson.NewObjectId()
			})

			Context("without customerID", func() {
				BeforeEach(func() {
					customerID = ""
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(Equal("invalid customerID"))
				})

				It("should NOT return a post", func() {
					Expect(retPost).To(BeNil())
				})
			})

			Context("with customerID", func() {
				BeforeEach(func() {
					customerID = "test-customer"
				})

				Context("with post not found", func() {
					BeforeEach(func() {
						ds.store = make(map[string]*dao.Post)
					})

					It("should return an error", func() {
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("post not found"))
					})

					It("should NOT return a post", func() {
						Expect(retPost).To(BeNil())
					})
				})

				Context("with post found", func() {
					var post *dao.Post
					BeforeEach(func() {
						ds.store = make(map[string]*dao.Post)
						storeID := createCompositeID(customerID, postID)
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
						ds.store[storeID] = post
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

	Describe("Update", func() {
		var (
			retPost *dao.Post
			err     error
			post    *dao.Post
		)

		JustBeforeEach(func() {
			retPost, err = ds.Update(customerID, post)
		})

		Context("without post", func() {
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("must provide post"))
			})

			It("should NOT return a post", func() {
				Expect(retPost).To(BeNil())
			})
		})

		Context("with post", func() {
			BeforeEach(func() {
				post = &dao.Post{
					CustID: customerID,
					URL:    "test-url1",
					Captions: []string{
						"caption4",
						"caption5",
						"caption6",
					},
				}
			})

			Context("without post.ID", func() {
				BeforeEach(func() {
					post.ID = nil
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(Equal("invalid postID"))
				})

				It("should NOT return a post", func() {
					Expect(retPost).To(BeNil())
				})
			})

			Context("with post.ID", func() {
				var postID bson.ObjectId
				BeforeEach(func() {
					postID = bson.NewObjectId()
					post.ID = &postID
				})

				Context("without customerID", func() {
					BeforeEach(func() {
						customerID = ""
					})

					It("should return an error", func() {
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("invalid customerID"))
					})

					It("should NOT return a post", func() {
						Expect(retPost).To(BeNil())
					})
				})

				Context("with customerID", func() {
					BeforeEach(func() {
						ds.store = make(map[string]*dao.Post)
						storeID := createCompositeID(customerID, postID)
						storedPost := &dao.Post{
							ID:     &postID,
							CustID: customerID,
							URL:    "test-url",
							Captions: []string{
								"caption1",
								"caption2",
								"caption3",
							},
						}
						ds.store[storeID] = storedPost
					})

					It("should NOT return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should return a post", func() {
						Expect(retPost).To(Equal(post))
					})

					It("should update the post", func() {
						storeID := createCompositeID(customerID, *retPost.ID)
						Expect(ds.store).To(HaveKeyWithValue(storeID, retPost))
					})
				})
			})
		})
	})
})
