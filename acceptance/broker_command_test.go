package acceptance

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
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

			It("successfully registers the service broker", func() {
				Eventually(session).Should(Exit(0))
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

		When("--name is not passed", func() {
			BeforeEach(func() {
				args = append(args, "--url", "url")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a Service Broker into the marketplace"))
			})
		})

		When("--url is not passed", func() {
			BeforeEach(func() {
				args = append(args, "--name", "name")
			})

			It("displays help and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("Usage:"))
				Eventually(session).Should(Say(`sm \[OPTIONS\] broker register \[register-OPTIONS\]`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("Register a Service Broker into the marketplace"))
			})
		})
	})
})
