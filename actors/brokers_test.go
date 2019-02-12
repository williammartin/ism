package actors_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/actors/actorsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Brokers Actor", func() {

	var (
		fakeBrokerRepository *actorsfakes.FakeBrokerRepository

		brokersActor *BrokersActor
	)

	BeforeEach(func() {
		fakeBrokerRepository = &actorsfakes.FakeBrokerRepository{}

		brokersActor = &BrokersActor{
			Repository: fakeBrokerRepository,
		}
	})

	Describe("GetBrokers", func() {
		var (
			brokers []*osbapi.Broker
			err     error
		)

		BeforeEach(func() {
			fakeBrokerRepository.FindAllReturns([]*osbapi.Broker{
				{Name: "broker-1"},
				{Name: "broker-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			brokers, err = brokersActor.GetBrokers()
		})

		It("finds all brokers from the repository", func() {
			Expect(fakeBrokerRepository.FindAllCallCount()).NotTo(BeZero())
			Expect(brokers).To(Equal([]*osbapi.Broker{
				{Name: "broker-1"},
				{Name: "broker-2"},
			}))
		})

		When("finding all brokers returns an error", func() {
			BeforeEach(func() {
				fakeBrokerRepository.FindAllReturns([]*osbapi.Broker{}, errors.New("error-finding-brokers"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-brokers"))
			})
		})
	})

	Describe("RegisterBroker", func() {
		var err error

		JustBeforeEach(func() {
			err = brokersActor.Register(&osbapi.Broker{
				Name: "broker-1",
			})
		})

		It("registers the broker", func() {
			Expect(fakeBrokerRepository.RegisterArgsForCall(0)).To(Equal(&osbapi.Broker{
				Name: "broker-1",
			}))
		})

		When("registering the broker fails", func() {
			BeforeEach(func() {
				fakeBrokerRepository.RegisterReturns(errors.New("error-registering-broker"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-registering-broker"))
			})
		})
	})
})
