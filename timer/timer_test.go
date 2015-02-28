package timer_test

import (
	"time"

	"github.com/mfine30/prodda/domain"
	domainfakes "github.com/mfine30/prodda/domain/fakes"
	"github.com/mfine30/prodda/timer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheduler", func() {
	var fakeTask *domainfakes.FakeTask
	frequency := domain.MiminumProdFrequency

	BeforeEach(func() {
		fakeTask = &domainfakes.FakeTask{}
	})

	// TODO: #88871252 add tests for propagation of error from domain.NewProd

	Describe("#Start", func() {
		It("Runs a task when the scheduler expires", func() {
			dingAt := time.Now().Add(1 * time.Second)
			prod, err := domain.NewProd(dingAt, fakeTask, 0)
			Expect(err).NotTo(HaveOccurred())
			scheduler, err := timer.NewScheduler(prod)
			Expect(err).NotTo(HaveOccurred())

			err = <-scheduler.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTask.RunCallCount()).To(Equal(1))
		})
	})

	Describe("#Cancel", func() {
		It("cancels the scheduler if it is running", func() {
			dingAt := time.Now().Add(300 * time.Millisecond)
			prod, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			scheduler, err := timer.NewScheduler(prod)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				time.Sleep(100 * time.Millisecond)
				scheduler.Cancel()
			}()

			err = <-scheduler.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTask.RunCallCount()).To(Equal(0))
		})

		It("returns an error if the scheduler has already finished", func() {
			dingAt := time.Now().Add(50 * time.Millisecond)
			prod, err := domain.NewProd(dingAt, fakeTask, 0)
			Expect(err).NotTo(HaveOccurred())
			scheduler, err := timer.NewScheduler(prod)
			Expect(err).NotTo(HaveOccurred())

			err = <-scheduler.Start()
			Expect(err).NotTo(HaveOccurred())
			err = scheduler.Cancel()
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the scheduler has not been started", func() {
			dingAt := time.Now().Add(300 * time.Millisecond)
			prod, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			scheduler, err := timer.NewScheduler(prod)
			Expect(err).NotTo(HaveOccurred())

			err = scheduler.Cancel()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("#UpdateScheduler", func() {
		var scheduler *timer.Scheduler
		var originalDingAt time.Time

		BeforeEach(func() {
			var err error
			originalDingAt = time.Now().Add(500 * time.Millisecond)
			prod, err := domain.NewProd(originalDingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			scheduler, err = timer.NewScheduler(prod)
			Expect(err).NotTo(HaveOccurred())
			Expect(scheduler.NextTime()).To(Equal(originalDingAt))
		})

		// TODO: #88871252 add tests for propagation of error from prod.Update

		Context("when scheduler is running", func() {
			BeforeEach(func() {
				var err error
				prod, err := domain.NewProd(originalDingAt, fakeTask, frequency)
				Expect(err).NotTo(HaveOccurred())
				scheduler, err = timer.NewScheduler(prod)
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduler.NextTime()).To(Equal(originalDingAt))
				scheduler.Start()
			})

			It("successfully updates the scheduler when time is in the future", func() {
				newDingAt := time.Now().Add(70 * time.Second)
				err := scheduler.Update(newDingAt, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduler.NextTime()).To(Equal(newDingAt))
			})
		})

		Context("when scheduler is not running", func() {
			It("successfully updates the scheduler when time is in the future", func() {
				newDingAt := time.Now().Add(70 * time.Second)
				err := scheduler.Update(newDingAt, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduler.NextTime()).To(Equal(newDingAt))
			})
		})
	})
})
