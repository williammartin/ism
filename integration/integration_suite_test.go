package integration

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var pathToSMCLI string

func TestIntegration(t *testing.T) {
	BeforeSuite(func() {
		var err error
		pathToSMCLI, err = Build("github.com/pivotal-cf/ism/cmd/sm")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
