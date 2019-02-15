package usecases_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/ism/osbapi"
	. "github.com/pivotal-cf/ism/usecases"
	"github.com/pivotal-cf/ism/usecases/usecasesfakes"
)

var _ = Describe("Service List Usecase", func() {

	var (
		fakeBrokersActor  *usecasesfakes.FakeBrokersActor
		fakeServicesActor *usecasesfakes.FakeServicesActor
		fakePlansActor    *usecasesfakes.FakePlansActor

		serviceListUsecase ServiceListUsecase

		services   []*Service
		executeErr error
	)

	BeforeEach(func() {
		fakeBrokersActor = &usecasesfakes.FakeBrokersActor{}
		fakeServicesActor = &usecasesfakes.FakeServicesActor{}
		fakePlansActor = &usecasesfakes.FakePlansActor{}

		serviceListUsecase = ServiceListUsecase{
			BrokersActor:  fakeBrokersActor,
			ServicesActor: fakeServicesActor,
			PlansActor:    fakePlansActor,
		}
	})

	JustBeforeEach(func() {
		services, executeErr = serviceListUsecase.GetServices()
	})

	It("fetches all brokers", func() {
		Expect(fakeBrokersActor.GetBrokersCallCount()).NotTo(BeZero())
	})

	When("fetching brokers errors", func() {
		BeforeEach(func() {
			fakeBrokersActor.GetBrokersReturns([]*osbapi.Broker{}, errors.New("error-getting-brokers"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-brokers"))
		})
	})

	When("there are no brokers", func() {
		BeforeEach(func() {
			fakeBrokersActor.GetBrokersReturns([]*osbapi.Broker{}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("returns an empty list of services", func() {
			Expect(services).To(HaveLen(0))
		})
	})

	When("there are one or more brokers", func() {
		BeforeEach(func() {
			fakeBrokersActor.GetBrokersReturns([]*osbapi.Broker{
				{ID: "broker1-id", Name: "broker1"},
				{ID: "broker2-id", Name: "broker2"}}, nil)
		})

		It("fetches services for each broker", func() {
			Expect(fakeServicesActor.GetServicesCallCount()).To(Equal(2))
			Expect(fakeServicesActor.GetServicesArgsForCall(0)).To(Equal("broker1-id"))
			Expect(fakeServicesActor.GetServicesArgsForCall(1)).To(Equal("broker2-id"))
		})

		When("fetching services errors", func() {
			BeforeEach(func() {
				fakeServicesActor.GetServicesReturns([]*osbapi.Service{}, errors.New("error-getting-services"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-getting-services"))
			})
		})

		When("all the brokers have services", func() {
			BeforeEach(func() {
				fakeServicesActor.GetServicesReturnsOnCall(0, []*osbapi.Service{
					{ID: "service1-id", Name: "service1", Description: "service1 description"},
					{ID: "service2-id", Name: "service2", Description: "service2 description"}}, nil)

				fakeServicesActor.GetServicesReturnsOnCall(1, []*osbapi.Service{
					{ID: "service3-id", Name: "service3", Description: "service3 description"},
					{ID: "service4-id", Name: "service4", Description: "service4 description"}}, nil)
			})

			It("fetches plans for each service", func() {
				Expect(fakePlansActor.GetPlansCallCount()).To(Equal(4))
				Expect(fakePlansActor.GetPlansArgsForCall(0)).To(Equal("service1-id"))
				Expect(fakePlansActor.GetPlansArgsForCall(1)).To(Equal("service2-id"))
				Expect(fakePlansActor.GetPlansArgsForCall(2)).To(Equal("service3-id"))
				Expect(fakePlansActor.GetPlansArgsForCall(3)).To(Equal("service4-id"))
			})

			When("fetching plans errors", func() {
				BeforeEach(func() {
					fakePlansActor.GetPlansReturns([]*osbapi.Plan{}, errors.New("error-getting-plans"))
				})

				It("propagates the error", func() {
					Expect(executeErr).To(MatchError("error-getting-plans"))
				})
			})

			When("all the services have plans", func() {
				BeforeEach(func() {
					fakePlansActor.GetPlansReturnsOnCall(0, []*osbapi.Plan{{Name: "plan1"}, {Name: "extra-plan"}}, nil)
					fakePlansActor.GetPlansReturnsOnCall(1, []*osbapi.Plan{{Name: "plan2"}}, nil)
					fakePlansActor.GetPlansReturnsOnCall(2, []*osbapi.Plan{{Name: "plan3"}}, nil)
					fakePlansActor.GetPlansReturnsOnCall(3, []*osbapi.Plan{{Name: "plan4"}}, nil)
				})

				It("doesn't error", func() {
					Expect(executeErr).NotTo(HaveOccurred())
				})

				It("returns a list of services", func() {
					Expect(services[0]).To(Equal(&Service{Name: "service1", Description: "service1 description", PlanNames: []string{"plan1", "extra-plan"}, BrokerName: "broker1"}))
					Expect(services[1]).To(Equal(&Service{Name: "service2", Description: "service2 description", PlanNames: []string{"plan2"}, BrokerName: "broker1"}))
					Expect(services[2]).To(Equal(&Service{Name: "service3", Description: "service3 description", PlanNames: []string{"plan3"}, BrokerName: "broker2"}))
					Expect(services[3]).To(Equal(&Service{Name: "service4", Description: "service4 description", PlanNames: []string{"plan4"}, BrokerName: "broker2"}))
				})
			})
		})
	})
})
