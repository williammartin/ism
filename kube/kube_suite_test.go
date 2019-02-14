package kube_test

import (
	"path/filepath"
	"testing"

	"github.com/pivotal-cf/ism/pkg/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var cfg *rest.Config

func TestKube(t *testing.T) {
	var testEnv *envtest.Environment

	BeforeSuite(func() {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "config", "crds")},
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
	RunSpecs(t, "Kube Suite")
}
