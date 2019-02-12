package actors_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/actors/actorsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Plans Actor", func() {

	var (
		fakePlanRepository *actorsfakes.FakePlanRepository

		plansActor *PlansActor
	)

	BeforeEach(func() {
		fakePlanRepository = &actorsfakes.FakePlanRepository{}

		plansActor = &PlansActor{
			Repository: fakePlanRepository,
		}
	})

	Describe("GetPlans", func() {
		var (
			plans []*osbapi.Plan
			err   error
		)

		BeforeEach(func() {
			fakePlanRepository.FindByServiceReturns([]*osbapi.Plan{
				{Name: "plan-1"},
				{Name: "plan-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			plans, err = plansActor.GetPlans("service-1")
		})

		It("finds plans by service id", func() {
			Expect(fakePlanRepository.FindByServiceArgsForCall(0)).To(Equal("service-1"))

			Expect(plans).To(Equal([]*osbapi.Plan{
				{Name: "plan-1"},
				{Name: "plan-2"},
			}))
		})

		When("finding plans returns an error", func() {
			BeforeEach(func() {
				fakePlanRepository.FindByServiceReturns([]*osbapi.Plan{}, errors.New("error-finding-plans"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-plans"))
			})
		})
	})
})
