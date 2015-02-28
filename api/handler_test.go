package api_test

import (
	"net/http"

	"github.com/mfine30/prodda/api"
	apifakes "github.com/mfine30/prodda/api/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Handler", func() {
	username := "username"
	password := "password"
	var handler http.Handler

	Context("when a request panics", func() {

		JustBeforeEach(func() {
			logger := lagertest.NewTestLogger("Handler Test")
			handler = api.NewHandler(logger, username, password, nil)
		})

		var (
			realHomeHandleFunc func(rw http.ResponseWriter, r *http.Request)
			responseWriter     *apifakes.FakeResponseWriter
			request            *http.Request
		)

		BeforeEach(func() {
			realHomeHandleFunc = api.HomeHandleFunc
			api.HomeHandleFunc = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				panic("fake request panic")
			})

			responseWriter = &apifakes.FakeResponseWriter{}
			var err error
			request, err = http.NewRequest("GET", "/api/v0", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			api.HomeHandleFunc = realHomeHandleFunc
		})

		It("recovers from panics and responds with an internal server error", func() {
			handler.ServeHTTP(responseWriter, request) // should not panic

			Expect(responseWriter.WriteHeaderCallCount()).To(Equal(1))
			Expect(responseWriter.WriteHeaderArgsForCall(0)).To(Equal(http.StatusInternalServerError))
		})
	})
})
