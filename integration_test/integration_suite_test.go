package integration_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const (
	travisURL = "https://api.travis-ci.org"
)

var (
	travisToken string
)

var _ = BeforeSuite(func() {
	travisToken = os.Getenv("TRAVIS_TOKEN")
	if travisToken == "" {
		Fail("TRAVIS_TOKEN must be provided")
	}
})

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}
