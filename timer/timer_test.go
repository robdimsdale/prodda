package timer_test

import (
	"time"

	"github.com/mfine30/prodda/fakes"
	"github.com/mfine30/prodda/timer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Timer", func() {
	It("creates a timer that finishes at the specified time", func() {
		dingAt := time.Now().Add(45 * time.Second)
		alarm := timer.MakeTimer(dingAt.Year(), dingAt.Month(), dingAt.Day(), dingAt.Hour(), dingAt.Minute(), dingAt.Second())
		Expect(alarm.FinishesAt).To(Equal(dingAt.Truncate(time.Second)))
	})

	It("Runs a task when the timer expires", func() {
		fakeTask := new(fakes.FakeTask)
		fakeTask.RunReturns(nil)
		tick := time.NewTicker(time.Millisecond * 250)
		alarm := timer.Alarm{}
		alarm.Ticker = tick

		err := alarm.RunOnDing(fakeTask)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeTask.RunCallCount()).To(Equal(1))
	})

	It("Can cancel the timer before it expires", func() {
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
})
