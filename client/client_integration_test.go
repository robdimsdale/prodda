package client_test

import (
	"github.com/mfine30/prodda/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientIntegration", func() {

	It("Can restart a specific build on travis", func() {
		travisClient := client.NewTravisClient(travisURL)
		resp, err := travisClient.TriggerBuild("mfine30", "prodda", travisToken, 50151622)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(ContainSubstring("successfully restarted"))
	})
})
