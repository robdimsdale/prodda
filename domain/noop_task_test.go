package domain_test

import (
	"time"

	"github.com/mfine30/prodda/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("No-op task", func() {
	var testLogger *lagertest.TestLogger
	schedule := ""

	BeforeEach(func() {
		testLogger = lagertest.NewTestLogger("no-op task test")
	})

	It("logs and sleeps", func() {
		sleepDuration := 50 * time.Millisecond
		task := domain.NewNoOpTask(schedule, sleepDuration, testLogger)

		startTime := time.Now()
		task.Run()
		duration := time.Now().Sub(startTime)
		Expect(duration).To(BeNumerically(">=", sleepDuration))
		Expect(testLogger.Buffer()).To(Say("started"))
		Expect(testLogger.Buffer()).To(Say("completed"))
	})
})
