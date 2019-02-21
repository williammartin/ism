package repositories

import (
	"context"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ctx = context.TODO()

//go:generate counterfeiter . KubeBrokerRepo

//TODO: move to internal reconciler
type KubeBrokerRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.Broker, error)
	UpdateState(broker *v1alpha1.Broker, newState v1alpha1.BrokerState) error
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

	err := repo.client.Get(ctx, resource, broker)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

func (repo *kubeBrokerRepo) UpdateState(broker *v1alpha1.Broker, newState v1alpha1.BrokerState) error {
	broker.Status.State = newState

	return repo.client.Status().Update(ctx, broker)
}
