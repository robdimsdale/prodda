package client_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/mfine30/prodda/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	It("Can get all builds for a specific repository", func() {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()

			Expect(r.URL.String()).To(Equal("/repos/fake_user/fake_repo/builds"))
			firstBuild := client.Build{
				Id: 1234,
			}
			secondBuild := client.Build{
				Id: 5678,
			}
			travisBuilds := []client.Build{
				secondBuild,
				firstBuild,
			} //build at index 0 is most recent

			jsonBuilds, err := json.Marshal(travisBuilds)
			Expect(err).NotTo(HaveOccurred())
			w.Write(jsonBuilds)
		}))

		travisClient := client.NewTravisClient(testServer.URL)
		travisBuilds, err := travisClient.GetBuilds("fake_user", "fake_repo")
		Expect(err).NotTo(HaveOccurred())

		builds := *travisBuilds
		Expect(builds[0]).To(Equal(client.Build{Id: 5678}))
	})

	It("Can restart a specific build for a given repo on travis", func() {
		restartNoticeText := "The build was successfully restarted."

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()

			Expect(r.URL.String()).To(Equal("/requests"))
			bodyBytes, err := ioutil.ReadAll(r.Body)
			body := string(bodyBytes)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal(`{"build_id": 1234}`))
			Expect(r.Header.Get("Authorization")).To(Equal("token fake-travis-token"))

			restartNotice := client.RestartNotice{
				Notice: restartNoticeText,
			}
			restartResponse := client.RestartResponse{
				Result: true,
				Flash:  []client.RestartNotice{restartNotice},
			}

			jsonResponse, err := json.Marshal(restartResponse)
			Expect(err).NotTo(HaveOccurred())
			w.Write(jsonResponse)
		}))

		travisClient := client.NewTravisClient(testServer.URL)
		resp, err := travisClient.TriggerBuild("fake_user", "fake_repo", "fake-travis-token", 1234)

		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(restartNoticeText))
	})
})
