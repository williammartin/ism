package acceptance

import (
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI service command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"service"}
	})

	JustBeforeEach(func() {
		var err error

		command := exec.Command(pathToCLI, args...)
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
			Eventually(session).Should(Say(`sm \[OPTIONS\] service <list>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The service command group lets you list the available services in the"))
			Eventually(session).Should(Say("marketplace\\."))
		})
	})

	Describe("list sub command", func() {
		BeforeEach(func() {
			args = append(args, "list")
		})

		When("--help is passed", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] service list`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("List the services that are available in the marketplace\\."))
			})
		})

		//TODO: -randomizeAllSpecs causes this test to flake due to shared k8s state
		When("0 brokers are registered", func() {
			It("displays 'No services found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No services found\\."))
			})
		})

		//TODO: -randomizeAllSpecs causes this test to flake due to shared k8s state
		When("1 broker is registered", func() {
			BeforeEach(func() {
				brokerURL := os.Getenv("BROKER_URL")
				brokerUsername := os.Getenv("BROKER_USERNAME")
				brokerPassword := os.Getenv("BROKER_PASSWORD")

				Expect(brokerURL).NotTo(Equal(""))
				Expect(brokerUsername).NotTo(Equal(""))
				Expect(brokerPassword).NotTo(Equal(""))

				registerArgs := []string{"broker", "register",
					"--name", "test-broker",
					"--url", brokerURL,
					"--username", brokerUsername,
					"--password", brokerPassword}
				command := exec.Command(pathToCLI, registerArgs...)
				registerSession, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(registerSession).Should(Exit(0))

				//TODO: Temporarily sleep until #164240938 is done.
				time.Sleep(3 * time.Second)
			})

			AfterEach(func() {
				deleteBrokers("test-broker")
			})

			It("displays services and plans for the broker", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("SERVICE\\s+PLANS\\s+BROKER\\s+DESCRIPTION"))
				Eventually(session).Should(Say("overview-service\\s+simple, complex\\s+test-broker\\s+Provides an overview of any service instances and bindings that have been created by a platform"))
			})
		})
	})
})
