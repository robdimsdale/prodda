package domain_test

import (
	"github.com/mfine30/prodda/domain"
	domainfakes "github.com/mfine30/prodda/domain/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prod", func() {
	var fakeTask *domainfakes.FakeTask

	BeforeEach(func() {
		fakeTask = &domainfakes.FakeTask{}
	})

	Describe("#NewProd", func() {
		It("creates a prod", func() {
			schedule := "some schedule"
			prod, err := domain.NewProd(fakeTask, schedule)
			Expect(err).NotTo(HaveOccurred())
			Expect(prod).NotTo(BeNil())
			Expect(prod.Task).To(Equal(fakeTask))
			Expect(prod.Schedule).To(Equal(schedule))
		})
	})
})
