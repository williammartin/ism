package repositories_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/pivotal-cf/ism/pkg/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var cfg *rest.Config

func TestRepositories(t *testing.T) {
	var testEnv *envtest.Environment

	SetDefaultEventuallyTimeout(time.Second * 5)
	SetDefaultConsistentlyDuration(time.Second * 5)

	BeforeSuite(func() {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crds")},
		}
		apis.AddToScheme(scheme.Scheme)

		var err error
		cfg, err = testEnv.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		testEnv.Stop()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repositories Suite")
}
