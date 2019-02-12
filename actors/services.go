package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . ServiceRepository

// TODO: godoc
type ServiceRepository interface {
	FindByBroker(brokerID string) ([]*osbapi.Service, error)
}

// TODO: godoc
type ServicesActor struct {
	Repository ServiceRepository
}

// TODO: godoc
func (a *ServicesActor) GetServices(brokerID string) ([]*osbapi.Service, error) {
	return a.Repository.FindByBroker(brokerID)
}
