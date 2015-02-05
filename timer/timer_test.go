package timer_test

import (
	"time"

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
})
