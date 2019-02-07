package acceptance

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI broker command", func() {
	When("--help is passed", func() {
		It("displays help and exits 0", func() {
			command := exec.Command(pathToSMCLI, "broker", "--help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))
			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`sm \[OPTIONS\] broker <register>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The broker command group lets you register, update and deregister Service"))
			Eventually(session).Should(Say("Brokers from the marketplace"))
		})
	})

	Describe("register sub command", func() {
		When("valid args are passed", func() {
			It("successfully registers the service broker", func() {
				command := exec.Command(pathToSMCLI, "broker", "register", "--name", "my-broker", "--url", "url", "--username", "username", "--password", "password")
				session, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Broker 'my-broker' registered\\."))
			})
		})

		When("--help is passed", func() {
			It("displays help and exits 0", func() {
				command := exec.Command(pathToSMCLI, "broker", "register", "--help")
				session, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a Service Broker into the marketplace"))
			})
		})
	})
})
