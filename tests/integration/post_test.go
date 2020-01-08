// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"labix.org/v2/mgo/bson"

	"github.com/bpross/cc-hw/dao"
	"github.com/bpross/cc-hw/handler"
)

// This describes the tests enumerated in the design doc
var _ = Describe("Post", func() {
	var (
		reqUrl   string
		url      string
		resp     *http.Response
		postBody *dao.Post
		getBody  *dao.Post
	)

	Describe("Test Case 1", func() {
		BeforeEach(func() {
			// Make POST request
			reqUrl = "http://api:8080/post"
			url = "https://blog.cloudcampaign.io/2018/03/11/7-social-media-stats-you-can-leverage-to-land-more-clients/"
			bodyStruct := handler.GeneratePostRequest{
				URL: url,
			}
			body, err := json.Marshal(bodyStruct)
			Expect(err).To(BeNil())
			req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			postBody = &dao.Post{}
			err = json.Unmarshal(respBody, postBody)
			Expect(err).To(BeNil())

			// Make GET request
			reqUrl = reqUrl + "/" + postBody.ID.Hex()
			req, err = http.NewRequest("GET", reqUrl, nil)
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			respBody, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			getBody = &dao.Post{}
			err = json.Unmarshal(respBody, getBody)
			Expect(err).To(BeNil())
		})

		It("should return the same information on get that was returned on post", func() {
			Expect(getBody).To(Equal(postBody))
		})
	})

	Describe("Test Case 2", func() {
		BeforeEach(func() {
			client := &http.Client{}
			dneID := bson.NewObjectId()
			reqUrl = fmt.Sprintf("http://api:8080/post/%s", dneID.Hex())
			req, err := http.NewRequest("GET", reqUrl, nil)
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
		})

		It("should return 404", func() {
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Describe("Test Case 3", func() {
		BeforeEach(func() {
			// Make POST request
			reqUrl = "http://api:8080/post"
			url = "http://google.com"
			bodyStruct := handler.GeneratePostRequest{
				URL: url,
			}
			body, err := json.Marshal(bodyStruct)
			Expect(err).To(BeNil())
			req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			postBody = &dao.Post{}
			err = json.Unmarshal(respBody, postBody)
			Expect(err).To(BeNil())
		})

		It("should not return captions", func() {
			Expect(postBody.Captions).To(BeEmpty())
		})
	})

	Describe("Test Case 4", func() {
		var newCaptions []string
		BeforeEach(func() {
			// Make POST request
			reqUrl = "http://api:8080/post"
			url = "https://blog.cloudcampaign.io/2018/03/11/7-social-media-stats-you-can-leverage-to-land-more-clients/"
			bodyStruct := handler.GeneratePostRequest{
				URL: url,
			}
			body, err := json.Marshal(bodyStruct)
			Expect(err).To(BeNil())
			req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			postBody = &dao.Post{}
			err = json.Unmarshal(respBody, postBody)
			Expect(err).To(BeNil())

			// Make PUT request
			reqUrl = reqUrl + "/" + postBody.ID.Hex()
			newCaptions = []string{
				"test1",
			}
			updateStruct := dao.Post{
				Captions: newCaptions,
			}
			body, err = json.Marshal(updateStruct)
			Expect(err).To(BeNil())
			req, err = http.NewRequest("PUT", reqUrl, bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			firstID := postBody.ID.Hex()

			// Make second POST request
			reqUrl = "http://api:8080/post"
			url = "https://blog.cloudcampaign.io/2018/03/11/7-social-media-stats-you-can-leverage-to-land-more-clients/"
			bodyStruct = handler.GeneratePostRequest{
				URL: url,
			}
			body, err = json.Marshal(bodyStruct)
			Expect(err).To(BeNil())
			req, err = http.NewRequest("POST", reqUrl, bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Make GET request
			reqUrl = reqUrl + "/" + firstID
			req, err = http.NewRequest("GET", reqUrl, nil)
			req.Header.Add("x-customer-id", "1")
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			Expect(err).To(BeNil())
			respBody, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			getBody = &dao.Post{}
			err = json.Unmarshal(respBody, getBody)
			Expect(err).To(BeNil())
		})

		It("should have the captions from PUT request", func() {
			Expect(getBody.Captions).To(Equal(newCaptions))
		})
	})
})
