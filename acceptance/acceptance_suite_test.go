package acceptance

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var (
	pathToSMCLI       string
	kubeClient        client.Client
	controllerSession *Session
	testEnv           *envtest.Environment
)

func TestAcceptance(t *testing.T) {
	SetDefaultEventuallyTimeout(time.Second * 5)
	SetDefaultConsistentlyDuration(time.Second * 5)

	BeforeSuite(func() {
		var err error
		pathToSMCLI, err = Build("github.com/pivotal-cf/ism/cmd/sm")
		Expect(err).NotTo(HaveOccurred())

		testEnv = &envtest.Environment{
			UseExistingCluster: true,
		}

		testEnvConfig, err := testEnv.Start()
		Expect(err).NotTo(HaveOccurred())

		Expect(v1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())

		kubeClient, err = client.New(testEnvConfig, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		installCRDs()
		startController()
	})

	AfterSuite(func() {
		Expect(testEnv.Stop()).To(Succeed())
		terminateController()
		CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func installCRDs() {
	command := exec.Command("make", "install")
	command.Dir = filepath.Join("..")
	command.Stdout = GinkgoWriter
	command.Stderr = GinkgoWriter
	Expect(command.Run()).To(Succeed())
}

func startController() {
	pathToController, err := Build("github.com/pivotal-cf/ism/cmd/manager")
	Expect(err).NotTo(HaveOccurred())

	command := exec.Command(pathToController)
	controllerSession, err = Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
}

func terminateController() {
	controllerSession.Terminate()
}

func ensureBrokerExists(brokerName string) {
	key := types.NamespacedName{
		Name:      brokerName,
		Namespace: "default",
	}

	fetched := &v1alpha1.Broker{}
	Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())
}

func deleteBrokers(brokerNames ...string) {
	for _, b := range brokerNames {
		bToDelete := &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      b,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), bToDelete)).To(Succeed())
	}
}
