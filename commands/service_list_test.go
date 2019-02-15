package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
	"github.com/pivotal-cf/ism/usecases"
)

var _ = Describe("Service List Command", func() {

	var (
		fakeUsecase *commandsfakes.FakeServiceListUsecase
		fakeUI      *commandsfakes.FakeUI

		listCommand ServiceListCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeUsecase = &commandsfakes.FakeServiceListUsecase{}
		fakeUI = &commandsfakes.FakeUI{}

		listCommand = ServiceListCommand{
			ServiceListUsecase: fakeUsecase,
			UI:                 fakeUI,
		}
	})

	JustBeforeEach(func() {
		executeErr = listCommand.Execute(nil)
	})

	When("there are no services", func() {
		BeforeEach(func() {
			fakeUsecase.GetServicesReturns([]*usecases.Service{}, nil)
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

	When("there is 1 or more services", func() {
		BeforeEach(func() {
			fakeUsecase.GetServicesReturns([]*usecases.Service{
				{
					Name:        "redis",
					Description: "redis service description",
					PlanNames:   []string{"small", "large"},
					BrokerName:  "redis-broker",
				},
				{
					Name:        "mysql",
					Description: "mysql service description",
					PlanNames:   []string{"medium"},
					BrokerName:  "mysql-broker",
				},
			}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays the services in a table", func() {
			Expect(fakeUI.DisplayTableCallCount()).NotTo(BeZero())
			data := fakeUI.DisplayTableArgsForCall(0)

			Expect(data[0]).To(Equal([]string{"SERVICE", "PLANS", "BROKER", "DESCRIPTION"}))
			Expect(data[1]).To(Equal([]string{"redis", "small, large", "redis-broker", "redis service description"}))
			Expect(data[2]).To(Equal([]string{"mysql", "medium", "mysql-broker", "mysql service description"}))
		})
	})

	When("getting services returns an error", func() {
		BeforeEach(func() {
			fakeUsecase.GetServicesReturns([]*usecases.Service{}, errors.New("error-getting-services"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-services"))
		})
	})
})
