package actors_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/actors/actorsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Services Actor", func() {

	var (
		fakeServiceRepository *actorsfakes.FakeServiceRepository

		servicesActor *ServicesActor
	)

	BeforeEach(func() {
		fakeServiceRepository = &actorsfakes.FakeServiceRepository{}

		servicesActor = &ServicesActor{
			Repository: fakeServiceRepository,
		}
	})

	Describe("GetServices", func() {
		var (
			services []*osbapi.Service
			err      error
		)

		BeforeEach(func() {
			fakeServiceRepository.FindByBrokerReturns([]*osbapi.Service{
				{Name: "service-1"},
				{Name: "service-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			services, err = servicesActor.GetServices("broker-1")
		})

		It("finds services by broker id", func() {
			Expect(fakeServiceRepository.FindByBrokerArgsForCall(0)).To(Equal("broker-1"))

			Expect(services).To(Equal([]*osbapi.Service{
				{Name: "service-1"},
				{Name: "service-2"},
			}))
		})

		When("finding services returns an error", func() {
			BeforeEach(func() {
				fakeServiceRepository.FindByBrokerReturns([]*osbapi.Service{}, errors.New("error-finding-services"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-services"))
			})
		})
	})

})
