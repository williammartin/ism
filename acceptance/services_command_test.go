package acceptance

import (
	"os"
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
			It("displays 'No services found.' and exits 0", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say("No services found\\."))
			})
		})

		When("1 broker is registered", func() {
			BeforeEach(func() {
				// Step 0 - deploy a broker

				// Step 1 - get broker url, name and password from ENV
				brokerURL := os.Getenv("BROKER_URL")
				brokerUsername := os.Getenv("BROKER_USERNAME")
				brokerPassword := os.Getenv("BROKER_PASSWORD")

				// Step 2 - run sm broker register --name x -- etc.
				registerArgs := []string{"broker", "register",
					"--name", "test-broker",
					"--url", brokerURL,
					"--username", brokerUsername,
					"--password", brokerPassword}
				command := exec.Command(pathToSMCLI, registerArgs...)
				registerSession, err := Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(registerSession).Should(Exit(0))
			})

			XIt("displays services and plans for the broker", func() {
				Eventually(session).Should(Exit(0))
				Eventually(session).Should(Say(`^NAME\\s+PLANS\\s+BROKER\\s+DESCRIPTION$`))
				Eventually(session).Should(Say(`^overview-service\\s+simple\\s+test-broker\\s+lol whatevs$`))
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
