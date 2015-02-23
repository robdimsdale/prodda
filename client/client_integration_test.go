package client_test

import (
	"fmt"

	"github.com/mfine30/prodda/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientIntegration", func() {

	It("Can restart a specific build on travis", func() {
		travisClient := client.NewTravisClient(travisURL)
		resp, err := travisClient.TriggerBuild(travisToken, 50151622)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(resp.Flash)).To(BeNumerically(">", 0))

		if resp.Flash[0].Notice == "" {
			Expect(resp.Flash[0].Error).NotTo(Equal(""))
		} else if resp.Flash[0].Error == "" {
			Expect(resp.Flash[0].Notice).NotTo(Equal(""))
		} else {
			Fail(fmt.Sprintf("Unexpected response Flash message: %+v", resp.Flash[0]))
		}
	})
})
