package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	pathToSMCLI string
	kubeClient  client.Client
)

func TestAcceptance(t *testing.T) {
	BeforeSuite(func() {
		var err error
		pathToSMCLI, err = Build("github.com/pivotal-cf/ism/cmd/sm")
		Expect(err).NotTo(HaveOccurred())

		kubeClient, err = buildKubeClient()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func buildKubeClient() (client.Client, error) {
	home := os.Getenv("HOME")
	kubeconfigFilepath := fmt.Sprintf("%s/.kube/config", home)
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigFilepath)
	if err != nil {
		return nil, err
	}

	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	return client.New(clientConfig, client.Options{Scheme: scheme.Scheme})
}
