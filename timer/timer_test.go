package timer_test

import (
	"time"

	"github.com/mfine30/prodda/timer"
	"github.com/mfine30/prodda/timer/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Timer", func() {
	var fakeTask *fakes.FakeTask

	BeforeEach(func() {
		fakeTask = &fakes.FakeTask{}
	})

	Describe("#NewAlarm", func() {
		It("creates an alarm that finishes at the specified time", func() {
			dingAt := time.Now().Add(45 * time.Second)
			alarm, err := timer.NewAlarm(dingAt, fakeTask)
			Expect(err).NotTo(HaveOccurred())
			Expect(alarm.FinishesAt).To(Equal(dingAt))
		})

		It("rejects the time if it is in the past", func() {
			dingAt := time.Now().Add(-45 * time.Second)
			_, err := timer.NewAlarm(dingAt, fakeTask)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Time must not be in the past"))
		})

		It("rejects a nil task", func() {
			dingAt := time.Now().Add(45 * time.Second)
			_, err := timer.NewAlarm(dingAt, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task must not be nil"))
		})
	})

	Describe("#Start", func() {
		It("Runs a task when the alarm expires", func() {
			dingAt := time.Now().Add(1 * time.Second)
			alarm, err := timer.NewAlarm(dingAt, fakeTask)

			err = alarm.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTask.RunCallCount()).To(Equal(1))
		})
	})

	Describe("#Cancel", func() {
		It("cancels the alarm if it is running", func() {
			dingAt := time.Now().Add(300 * time.Millisecond)
			alarm, err := timer.NewAlarm(dingAt, fakeTask)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				time.Sleep(100 * time.Millisecond)
				alarm.Cancel()
			}()
			err = alarm.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTask.RunCallCount()).To(Equal(0))
		})

		It("returns an error if the alarm has already finished", func() {
			dingAt := time.Now().Add(50 * time.Millisecond)
			alarm, err := timer.NewAlarm(dingAt, fakeTask)
			Expect(err).NotTo(HaveOccurred())

			err = alarm.Start()
			Expect(err).NotTo(HaveOccurred())
			err = alarm.Cancel()
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the alarm has not been started", func() {
			dingAt := time.Now().Add(300 * time.Millisecond)
			alarm, err := timer.NewAlarm(dingAt, fakeTask)
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
			alarm, err = timer.NewAlarm(originalDingAt, fakeTask)
			Expect(err).NotTo(HaveOccurred())
			Expect(alarm.FinishesAt).To(Equal(originalDingAt))
		})

		It("returns an error when the time is in the past", func() {
			newDingAt := time.Now().Add(-70 * time.Second)
			err := alarm.UpdateAlarm(newDingAt)
			Expect(err).To(HaveOccurred())
		})

		Context("when alarm is running", func() {
			BeforeEach(func() {
				var err error
				alarm, err = timer.NewAlarm(originalDingAt, fakeTask)
				Expect(err).NotTo(HaveOccurred())
				Expect(alarm.FinishesAt).To(Equal(originalDingAt))
				alarm.Start()
			})

			It("successfully updates the alarm when time is in the future", func() {
				newDingAt := time.Now().Add(70 * time.Second)
				err := alarm.UpdateAlarm(newDingAt)
				Expect(err).NotTo(HaveOccurred())
				Expect(alarm.FinishesAt).To(Equal(newDingAt))
			})
		})

		Context("when alarm is not running", func() {
			It("successfully updates the alarm when time is in the future", func() {
				newDingAt := time.Now().Add(70 * time.Second)
				err := alarm.UpdateAlarm(newDingAt)
				Expect(err).NotTo(HaveOccurred())
				Expect(alarm.FinishesAt).To(Equal(newDingAt))
			})
		})
	})
})
