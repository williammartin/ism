package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/ism/actors"
	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
)

var _ = Describe("Services List Command", func() {

	var (
		fakeServicesActor *commandsfakes.FakeServicesActor
		fakeUI            *commandsfakes.FakeUI

		listCommand ListCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeServicesActor = &commandsfakes.FakeServicesActor{}
		fakeUI = &commandsfakes.FakeUI{}

		listCommand = ListCommand{
			ServicesActor: fakeServicesActor,
			UI:            fakeUI,
		}
	})

	JustBeforeEach(func() {
		executeErr = listCommand.Execute(nil)
	})

	When("there are no services", func() {
		BeforeEach(func() {
			fakeServicesActor.GetServicesReturns([]actors.Service{}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays that no services were found", func() {
			Expect(fakeUI.DisplayTextCallCount()).NotTo(BeZero())
			text, _ := fakeUI.DisplayTextArgsForCall(0)

			Expect(text).To(Equal("No services found."))
		})
	})
})
