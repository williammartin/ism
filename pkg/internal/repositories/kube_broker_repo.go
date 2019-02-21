package repositories

import (
	"context"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate counterfeiter . KubeBrokerRepo

type KubeBrokerRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.Broker, error)
}

type kubeBrokerRepo struct {
	client client.Client
}

func NewKubeBrokerRepo(client client.Client) KubeBrokerRepo {
	return &kubeBrokerRepo{
		client: client,
	}
}

func (repo *kubeBrokerRepo) Get(resource types.NamespacedName) (*v1alpha1.Broker, error) {
	broker := &v1alpha1.Broker{}

	err := repo.client.Get(context.TODO(), resource, broker)
	if err != nil {
		return nil, err
	}

	return broker, nil
}
