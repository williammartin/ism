package kube

import (
	"context"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Service struct {
	KubeClient client.Client
}

func (s *Service) FindByBroker(brokerID string) ([]*osbapi.Service, error) {
	list := &v1alpha1.BrokerServiceList{}
	if err := s.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Service{}, err
	}

	services := []*osbapi.Service{}
	for _, s := range list.Items {
		if s.Spec.BrokerID == brokerID {
			services = append(services, &osbapi.Service{
				ID:          string(s.UID),
				Name:        s.Spec.Name,
				Description: s.Spec.Description,
				BrokerID:    s.Spec.BrokerID,
			})
		}
	}

	return services, nil
}
