package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . BrokerRepository

type BrokerRepository interface {
	FindAll() ([]*osbapi.Broker, error)
	Register(*osbapi.Broker) error
}

// TODO: godoc
type BrokersActor struct {
	Repository BrokerRepository
}

// TODO: godoc
func (a *BrokersActor) GetBrokers() ([]*osbapi.Broker, error) {
	return a.Repository.FindAll()
}

//TODO: Make names consistent
// TODO: godoc
func (a *BrokersActor) Register(broker *osbapi.Broker) error {
	return a.Repository.Register(broker)
}
