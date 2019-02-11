package acceptance

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CLI services command", func() {

	var (
		args    []string
		session *Session
	)

	BeforeEach(func() {
		args = []string{"services"}
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
			Eventually(session).Should(Say(`sm \[OPTIONS\] services <list>`))
			Eventually(session).Should(Say("\n"))
			Eventually(session).Should(Say("The services command group lets you list the available services in the"))
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
				Eventually(session).Should(Say(`sm \[OPTIONS\] services list`))
				Eventually(session).Should(Say("\n"))
				Eventually(session).Should(Say("List the services that are available in the marketplace\\."))
			})
		})

		When("0 brokers are registered", func() {
			It("displays 'No brokers found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No brokers found\\."))
			})
		})

		// 	When("valid args are passed", func() {
		// 		BeforeEach(func() {
		// 			args = append(args, "--name", "my-broker", "--url", "url", "--username", "username", "--password", "password")
		// 		})
		//
		// 		It("successfully registers the service broker", func() {
		// 			Eventually(session).Should(Exit(0))
		// 			Eventually(session).Should(Say("Broker 'my-broker' registered\\."))
	})
})
