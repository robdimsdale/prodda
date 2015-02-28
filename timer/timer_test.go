package timer_test

import (
	"time"

	"github.com/mfine30/prodda/domain"
	domainfakes "github.com/mfine30/prodda/domain/fakes"
	"github.com/mfine30/prodda/timer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Timer", func() {
	var fakeTask *domainfakes.FakeTask
	frequency := domain.MiminumProdFrequency

	BeforeEach(func() {
		fakeTask = &domainfakes.FakeTask{}
	})

	// TODO: #88871252 add tests for propagation of error from domain.NewProd

	Describe("#Start", func() {
		It("Runs a task when the alarm expires", func() {
			dingAt := time.Now().Add(1 * time.Second)
			prod, err := domain.NewProd(dingAt, fakeTask, 0)
			Expect(err).NotTo(HaveOccurred())
			alarm, err := timer.NewAlarm(prod)
			Expect(err).NotTo(HaveOccurred())

			err = <-alarm.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTask.RunCallCount()).To(Equal(1))
		})
	})

	Describe("#Cancel", func() {
		It("cancels the alarm if it is running", func() {
			dingAt := time.Now().Add(300 * time.Millisecond)
			prod, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			alarm, err := timer.NewAlarm(prod)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				time.Sleep(100 * time.Millisecond)
				alarm.Cancel()
			}()

			err = <-alarm.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTask.RunCallCount()).To(Equal(0))
		})

		It("returns an error if the alarm has already finished", func() {
			dingAt := time.Now().Add(50 * time.Millisecond)
			prod, err := domain.NewProd(dingAt, fakeTask, 0)
			Expect(err).NotTo(HaveOccurred())
			alarm, err := timer.NewAlarm(prod)
			Expect(err).NotTo(HaveOccurred())

			err = <-alarm.Start()
			Expect(err).NotTo(HaveOccurred())
			err = alarm.Cancel()
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the alarm has not been started", func() {
			dingAt := time.Now().Add(300 * time.Millisecond)
			prod, err := domain.NewProd(dingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			alarm, err := timer.NewAlarm(prod)
			Expect(err).NotTo(HaveOccurred())

			err = alarm.Cancel()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("#UpdateAlarm", func() {
		var alarm *timer.Alarm
		var originalDingAt time.Time

		BeforeEach(func() {
			var err error
			originalDingAt = time.Now().Add(500 * time.Millisecond)
			prod, err := domain.NewProd(originalDingAt, fakeTask, frequency)
			Expect(err).NotTo(HaveOccurred())
			alarm, err = timer.NewAlarm(prod)
			Expect(err).NotTo(HaveOccurred())
			Expect(alarm.NextTime()).To(Equal(originalDingAt))
		})

		// TODO: #88871252 add tests for propagation of error from prod.Update

		Context("when alarm is running", func() {
			BeforeEach(func() {
				var err error
				prod, err := domain.NewProd(originalDingAt, fakeTask, frequency)
				Expect(err).NotTo(HaveOccurred())
				alarm, err = timer.NewAlarm(prod)
				Expect(err).NotTo(HaveOccurred())
				Expect(alarm.NextTime()).To(Equal(originalDingAt))
				alarm.Start()
			})

			It("successfully updates the alarm when time is in the future", func() {
				newDingAt := time.Now().Add(70 * time.Second)
				err := alarm.Update(newDingAt, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(alarm.NextTime()).To(Equal(newDingAt))
			})
		})

		Context("when alarm is not running", func() {
			It("successfully updates the alarm when time is in the future", func() {
				newDingAt := time.Now().Add(70 * time.Second)
				err := alarm.Update(newDingAt, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(alarm.NextTime()).To(Equal(newDingAt))
			})
		})
	})
})
