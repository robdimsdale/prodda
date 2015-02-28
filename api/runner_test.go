package api_test

import (
	"fmt"
	"net"
	"os"

	"github.com/mfine30/prodda/api"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
)

var _ = Describe("APIRunner", func() {
	It("shuts down gracefully when signalled", func() {
		apiPort := 10000 + GinkgoParallelNode()
		logger := lagertest.NewTestLogger("APIRunner Test")
		username := "username"
		password := "password"
		handler := api.NewHandler(logger, username, password, nil)
		apiRunner := api.NewRunner(uint(apiPort), handler, logger)
		apiProcess := ifrit.Invoke(apiRunner)
		apiProcess.Signal(os.Kill)
		Eventually(apiProcess.Wait()).Should(Receive())

		_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", apiPort))
		Expect(err).To(HaveOccurred())
	})
})
