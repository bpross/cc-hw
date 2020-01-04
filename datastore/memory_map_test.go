package datastore

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"labix.org/v2/mgo/bson"
)

var _ = Describe("InMemoryDatastore", func() {
	var (
		logger     *log.Logger
		ds         *InMemoryDatastore
		record     *Record
		customerID string
	)

	BeforeEach(func() {
		// TODO pass in no-op logger, so we dont log during tests
		logger = log.New()
		ds = NewInMemoryDatastore(*logger)
		customerID = "test-customer"
	})

	Describe("Insert", func() {
		var (
			retRecord *Record
			err       error
		)

		JustBeforeEach(func() {
			retRecord, err = ds.Insert(customerID, record)
		})

		Context("without record", func() {
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("must provide record"))
			})

			It("should NOT return a record", func() {
				Expect(retRecord).To(BeNil())
			})
		})

		Context("with record.ID", func() {
			BeforeEach(func() {
				record = &Record{
					ID: bson.NewObjectId(),
				}
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("cannot provide ID"))
			})

			It("should NOT return a record", func() {
				Expect(retRecord).To(BeNil())
			})
		})

		Context("without record.ID", func() {
			BeforeEach(func() {
				record = &Record{
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

				It("should NOT return a record", func() {
					Expect(retRecord).To(BeNil())
				})
			})

			Context("with customerID", func() {
				BeforeEach(func() {
					customerID = "test-customer"
				})

				It("should NOT return an error", func() {
					Expect(err).To(BeNil())
				})

				It("should return a record", func() {
					Expect(retRecord.ID).NotTo(Equal(""))
					Expect(retRecord.CustID).To(Equal(customerID))
					Expect(retRecord.URL).To(Equal(record.URL))
					Expect(retRecord.Captions).To(Equal(record.Captions))
				})

				It("should insert a record", func() {
					storeID := createCompositeID(customerID, retRecord.ID)
					Expect(ds.store).To(HaveKeyWithValue(storeID, retRecord))
				})
			})
		})
	})

	Describe("Get", func() {
		var (
			retRecord *Record
			err       error
			recordID  bson.ObjectId
		)

		JustBeforeEach(func() {
			retRecord, err = ds.Get(customerID, recordID)
		})

		Context("without recordID", func() {
			BeforeEach(func() {
				recordID = ""
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("invalid recordID"))
			})

			It("should NOT return a record", func() {
				Expect(retRecord).To(BeNil())
			})
		})

		Context("with recordID", func() {
			BeforeEach(func() {
				recordID = bson.NewObjectId()
			})

			Context("without customerID", func() {
				BeforeEach(func() {
					customerID = ""
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(Equal("invalid customerID"))
				})

				It("should NOT return a record", func() {
					Expect(retRecord).To(BeNil())
				})
			})

			Context("with customerID", func() {
				BeforeEach(func() {
					customerID = "test-customer"
				})

				Context("with record not found", func() {
					BeforeEach(func() {
						ds.store = make(map[string]*Record)
					})

					It("should return an error", func() {
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("record not found"))
					})

					It("should NOT return a record", func() {
						Expect(retRecord).To(BeNil())
					})
				})

				Context("with record found", func() {
					var record *Record
					BeforeEach(func() {
						ds.store = make(map[string]*Record)
						storeID := createCompositeID(customerID, recordID)
						record = &Record{
							ID:     recordID,
							CustID: customerID,
							URL:    "test-url",
							Captions: []string{
								"caption1",
								"caption2",
								"caption3",
							},
						}
						ds.store[storeID] = record
					})

					It("should NOT return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should return a record", func() {
						Expect(retRecord).To(Equal(record))
					})
				})
			})
		})
	})

	Describe("Update", func() {
		var (
			retRecord *Record
			err       error
			record    *Record
		)

		JustBeforeEach(func() {
			retRecord, err = ds.Update(customerID, record)
		})

		Context("without record", func() {
			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("must provide record"))
			})

			It("should NOT return a record", func() {
				Expect(retRecord).To(BeNil())
			})
		})

		Context("with record", func() {
			BeforeEach(func() {
				record = &Record{
					CustID: customerID,
					URL:    "test-url1",
					Captions: []string{
						"caption4",
						"caption5",
						"caption6",
					},
				}
			})

			Context("without record.ID", func() {
				BeforeEach(func() {
					record.ID = ""
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(Equal("invalid recordID"))
				})

				It("should NOT return a record", func() {
					Expect(retRecord).To(BeNil())
				})
			})

			Context("with record.ID", func() {
				var recordID bson.ObjectId
				BeforeEach(func() {
					recordID = bson.NewObjectId()
					record.ID = recordID
				})

				Context("without customerID", func() {
					BeforeEach(func() {
						customerID = ""
					})

					It("should return an error", func() {
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("invalid customerID"))
					})

					It("should NOT return a record", func() {
						Expect(retRecord).To(BeNil())
					})
				})

				Context("with customerID", func() {
					BeforeEach(func() {
						ds.store = make(map[string]*Record)
						storeID := createCompositeID(customerID, recordID)
						storedRecord := &Record{
							ID:     recordID,
							CustID: customerID,
							URL:    "test-url",
							Captions: []string{
								"caption1",
								"caption2",
								"caption3",
							},
						}
						ds.store[storeID] = storedRecord
					})

					It("should NOT return an error", func() {
						Expect(err).To(BeNil())
					})

					It("should return a record", func() {
						Expect(retRecord).To(Equal(record))
					})

					It("should update the record", func() {
						storeID := createCompositeID(customerID, retRecord.ID)
						Expect(ds.store).To(HaveKeyWithValue(storeID, retRecord))
					})
				})
			})
		})
	})
})
