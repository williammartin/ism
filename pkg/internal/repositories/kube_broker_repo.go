package repositories

import (
	"context"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ctx = context.TODO()

type KubeBrokerRepo struct {
	client client.Client
}

func NewKubeBrokerRepo(client client.Client) *KubeBrokerRepo {
	return &KubeBrokerRepo{
		client: client,
	}
}

func (repo *KubeBrokerRepo) Get(resource types.NamespacedName) (*v1alpha1.Broker, error) {
	broker := &v1alpha1.Broker{}

	err := repo.client.Get(ctx, resource, broker)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

func (repo *KubeBrokerRepo) UpdateState(broker *v1alpha1.Broker, newState v1alpha1.BrokerState) error {
	broker.Status.State = newState

	return repo.client.Status().Update(ctx, broker)
}
