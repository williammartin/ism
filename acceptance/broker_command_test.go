package acceptance

import (
	"context"
	"os/exec"

	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("CLI broker command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"broker"}
	})

	JustBeforeEach(func() {
		var err error

		command := exec.Command(pathToSMCLI, args...)
		session, err = Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	When("--help is passed", func() {
		BeforeEach(func() {
			args = append(args, "--help")
		})

		It("displays help and exits 0", func() {
			Eventually(session).Should(Exit(0))
			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`sm \[OPTIONS\] broker <register>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The broker command group lets you register, update and deregister Service"))
			Eventually(session).Should(Say("Brokers from the marketplace"))
		})
	})

	Describe("register sub command", func() {
		BeforeEach(func() {
			args = append(args, "register")
		})

		When("valid args are passed", func() {
			BeforeEach(func() {
				args = append(args, "--name", "my-broker", "--url", "url", "--username", "username", "--password", "password")
			})

			AfterEach(func() {
				deleteBrokers("my-broker")
			})

			It("successfully registers the broker, and displays a message", func() {
				Eventually(session).Should(Exit(0))

				ensureBrokerExists("my-broker")

				Eventually(session).Should(Say("Broker 'my-broker' registered\\."))
			})
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a Service Broker into the marketplace"))
			})
		})

		When("required arguments are not passed", func() {
			It("displays an informative message and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("the required flags `--name', `--password', `--url' and `--username' were not specified"))
			})
		})
	})
})

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
