package main_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Routing", func() {
	var url string
	var session *gexec.Session

	BeforeEach(func() {
		os.Setenv("PORT", strconv.Itoa(appPort))
		os.Setenv("USERNAME", username)
		os.Setenv("PASSWORD", password)

		command := exec.Command(pathToExecutable)
		var err error
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say("Prodda started"))
		url = fmt.Sprintf("http://localhost:%d", appPort)
	})

	AfterEach(func() {
		session.Terminate().Wait()
	})

	Describe("/", func() {
		It("returns 200 without authentication", func() {
			resp, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("/api/v0/prods", func() {
		BeforeEach(func() {
			url = fmt.Sprintf("%s/api/v0/prods", url)
		})

		It("returns 401 when no credentials provided", func() {
			resp, err := http.Post(url, "", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("returns 401 when bad credentials provided", func() {
			req, err := http.NewRequest("GET", url, nil)
			Expect(err).NotTo(HaveOccurred())
			req.SetBasicAuth("baduser", "somepassword")
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		Context("when authenticated and authorized", func() {
			var req *http.Request

			BeforeEach(func() {
				var err error
				req, err = http.NewRequest("POST", url, nil)
				Expect(err).NotTo(HaveOccurred())
				req.SetBasicAuth(username, password)
			})

			It("returns 404 for GET", func() {
				req.Method = "GET"
				client := &http.Client{}
				resp, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("returns 200 for POST", func() {
				client := &http.Client{}
				resp, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
