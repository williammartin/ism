package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . PlanRepository

// TODO: godoc
type PlanRepository interface {
	FindByService(serviceID string) ([]*osbapi.Plan, error)
}

// TODO: godoc
type PlansActor struct {
	Repository PlanRepository
}

// TODO: godoc
func (a *PlansActor) GetPlans(serviceID string) ([]*osbapi.Plan, error) {
	return a.Repository.FindByService(serviceID)
}
