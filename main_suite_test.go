package main_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

const (
	username = "admin"
	password = "password"
)

var (
	pathToExecutable string
	appPort          int
)

var _ = BeforeSuite(func() {
	var err error
	pathToExecutable, err = gexec.Build("github.com/mfine30/prodda")
	appPort = 40000 + config.GinkgoConfig.ParallelNode
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestTimer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}
