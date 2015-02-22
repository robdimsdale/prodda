package timer_test

import (
	"time"

	"github.com/mfine30/prodda/timer"
	"github.com/mfine30/prodda/timer/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Timer", func() {
	Describe("#NewAlarm", func() {
		It("creates an alarm that finishes at the specified time", func() {
			dingAt := time.Now().Add(45 * time.Second)
			alarm, err := timer.NewAlarm(dingAt)
			Expect(err).NotTo(HaveOccurred())
			Expect(alarm.FinishesAt).To(Equal(dingAt))
		})

		It("rejects the time if it is in the past", func() {
			dingAt := time.Now().Add(-45 * time.Second)
			_, err := timer.NewAlarm(dingAt)
			Expect(err).To(HaveOccurred())
		})
	})

	It("Runs a task when the alarm expires", func() {
		fakeTask := new(fakes.FakeTask)
		fakeTask.RunReturns(nil)
		tick := time.NewTicker(time.Millisecond * 250)
		alarm := timer.Alarm{}
		alarm.Ticker = tick

		err := alarm.RunOnDing(fakeTask)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeTask.RunCallCount()).To(Equal(1))
	})

	It("Can cancel the alarm before it expires", func() {
		fakeTask := new(fakes.FakeTask)
		tick := time.NewTicker(time.Millisecond * 300)
		closeChannel := make(chan struct{})

		alarm := timer.Alarm{}
		alarm.Ticker = tick
		alarm.Alert = closeChannel

		go func() {
			time.Sleep(100 * time.Millisecond)
			close(alarm.Alert)
		}()
		err := alarm.RunOnDing(fakeTask)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeTask.RunCallCount()).To(Equal(0))
	})

	Describe("#UpdateAlarm", func() {
		var alarm *timer.Alarm
		var originalDingAt time.Time

		BeforeEach(func() {
			var err error
			originalDingAt = time.Now().Add(500 * time.Millisecond)
			alarm, err = timer.NewAlarm(originalDingAt)
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
				fakeTask := new(fakes.FakeTask)
				alarm.RunOnDing(fakeTask)
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
