package usecases

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . BrokersActor

//TODO: Rename this
type BrokersActor interface {
	GetBrokers() ([]*osbapi.Broker, error)
}

//go:generate counterfeiter . ServicesActor

//TODO: Rename this
type ServicesActor interface {
	GetServices(brokerID string) ([]*osbapi.Service, error)
}

//go:generate counterfeiter . PlansActor

//TODO: Rename this
type PlansActor interface {
	GetPlans(serviceID string) ([]*osbapi.Plan, error)
}

type ServiceListUsecase struct {
	BrokersActor  BrokersActor
	ServicesActor ServicesActor
	PlansActor    PlansActor
}

func (u *ServiceListUsecase) GetServices() ([]*Service, error) {
	brokers, err := u.BrokersActor.GetBrokers()
	if err != nil {
		return []*Service{}, err
	}

	var servicesToDisplay []*Service
	for _, b := range brokers {
		services, err := u.ServicesActor.GetServices(b.ID)
		if err != nil {
			return []*Service{}, err
		}

		for _, s := range services {
			plans, err := u.PlansActor.GetPlans(s.ID)
			if err != nil {
				return []*Service{}, err
			}

			serviceToDisplay := &Service{
				Name:        s.Name,
				Description: s.Description,
				PlanNames:   plansToNames(plans),
				BrokerName:  b.Name,
			}
			servicesToDisplay = append(servicesToDisplay, serviceToDisplay)
		}
	}

	return servicesToDisplay, nil
}

type Service struct {
	Name        string
	Description string
	PlanNames   []string
	BrokerName  string
}

func plansToNames(plans []*osbapi.Plan) []string {
	names := []string{}
	for _, p := range plans {
		names = append(names, p.Name)
	}

	return names
}
