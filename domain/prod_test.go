package domain_test

import (
	"time"

	"github.com/mfine30/prodda/domain"
	domainfakes "github.com/mfine30/prodda/domain/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prod", func() {
	var fakeTask *domainfakes.FakeTask
	frequency := domain.MiminumProdFrequency + 500*time.Millisecond

	BeforeEach(func() {
		fakeTask = &domainfakes.FakeTask{}
	})

	Describe("#NewProd", func() {
		It("creates a prod that finishes at the specified time", func() {
			dingAt := time.Now().Add(45 * time.Second)
			prod, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			Expect(prod.NextTime).To(Equal(dingAt))
		})

		It("rejects the time if it is in the past", func() {
			dingAt := time.Now().Add(-45 * time.Second)
			_, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Time must not be in the past"))
		})

		It("rejects a nil task", func() {
			dingAt := time.Now().Add(45 * time.Second)
			_, err := domain.NewProd(dingAt, nil, frequency)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task must not be nil"))
		})

		It("accepts frequency of 0", func() {
			dingAt := time.Now().Add(45 * time.Second)
			_, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
		})

		It("accepts frequency of MinimumFrequency", func() {
			dingAt := time.Now().Add(45 * time.Second)
			_, err := domain.NewProd(dingAt, fakeTask, domain.MiminumProdFrequency)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if the frequency is between 0 and MinimumFrequency", func() {
			frequency := 10 * time.Millisecond
			dingAt := time.Now().Add(45 * time.Second)
			_, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Frequency must either be 0 or greater than"))
		})
	})

	Describe("#Update", func() {
		var prod *domain.Prod
		var originalDingAt time.Time

		BeforeEach(func() {
			var err error
			originalDingAt = time.Now().Add(500 * time.Millisecond)
			prod, err = domain.NewProd(originalDingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			Expect(prod.NextTime).To(Equal(originalDingAt))
		})

		It("returns an error if the time is in the past", func() {
			newDingAt := time.Now().Add(-70 * time.Second)
			err := prod.Update(newDingAt, 0)
			Expect(err).To(HaveOccurred())
		})

		It("accepts frequency of 0", func() {
			newFrequency := time.Duration(0)
			err := prod.Update(originalDingAt, newFrequency)
			Expect(err).NotTo(HaveOccurred())
			Expect(prod.Frequency).To(Equal(time.Duration(0)))
		})

		It("accepts frequency of MinimumFrequency", func() {
			newFrequency := time.Duration(domain.MiminumProdFrequency)
			err := prod.Update(originalDingAt, newFrequency)
			Expect(err).NotTo(HaveOccurred())
			Expect(prod.Frequency).To(Equal(time.Duration(domain.MiminumProdFrequency)))
		})

		It("returns an error if the frequency is between 0 and MinimumFrequency", func() {
			newFrequency := 10 * time.Millisecond
			err := prod.Update(originalDingAt, newFrequency)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Frequency must either be 0 or greater than"))
		})
	})
})
