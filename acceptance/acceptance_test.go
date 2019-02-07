package acceptance

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI", func() {
	When("no command or flag is passed", func() {
		It("displays help for the cli and exits 0", func() {
			command := exec.Command(pathToSMCLI)
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`sm \[OPTIONS\] <broker>`))
		})
	})

	When("--help is passed", func() {
		It("displays help for the cli and exits 0", func() {
			command := exec.Command(pathToSMCLI, "--help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(0))

			Eventually(session).Should(Say("Usage:"))
			Eventually(session).Should(Say(`sm \[OPTIONS\] <broker>`))
		})
	})

	When("broker is passed", func() {
		When("register is passed with args", func() {
			It("successfully registers", func() {
				command := exec.Command(pathToSMCLI, "broker", "register", "--name", "my-broker", "--url", "url", "--username", "username", "--password", "password")
				session, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Broker 'my-broker' registered\\."))
			})
		})

		When("--help is passed", func() {
			It("displays help for the broker command and exits 0", func() {
				command := exec.Command(pathToSMCLI, "broker", "--help")
				session, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] broker <register>`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("The broker command group lets you register, update and deregister service"))
				Eventually(session).Should(Say("brokers from the marketplace"))
			})
		})
	})
})
