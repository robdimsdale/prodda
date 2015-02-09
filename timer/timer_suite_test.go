package prodda_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTimer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timer Suite")
}
