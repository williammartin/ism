package broker

import (
	"context"

	osbapiv1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var ctx = context.TODO()

//TODO Shall we leave this here or move it to the test?

//go:generate counterfeiter . BrokerClient

type BrokerClient interface {
	osbapi.Client
}

type BrokerReconciler struct {
	kubeBrokerRepo     repositories.KubeBrokerRepo
	kubeServiceRepo    repositories.KubeServiceRepo
	kubePlanRepo       repositories.KubePlanRepo
	createBrokerClient osbapi.CreateFunc
}

func NewBrokerReconciler(
	createBrokerClient osbapi.CreateFunc,
	kubeBrokerRepo repositories.KubeBrokerRepo,
	kubeServiceRepo repositories.KubeServiceRepo,
	kubePlanRepo repositories.KubePlanRepo,
) *BrokerReconciler {
	return &BrokerReconciler{
		createBrokerClient: createBrokerClient,
		kubeBrokerRepo:     kubeBrokerRepo,
		kubeServiceRepo:    kubeServiceRepo,
		kubePlanRepo:       kubePlanRepo,
	}
}

func (r *BrokerReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	broker, err := r.kubeBrokerRepo.Get(request.NamespacedName)
	if err != nil {
		return reconcile.Result{}, err
	}

	if broker.Status.State == osbapiv1alpha1.BrokerStateRegistered {
		return reconcile.Result{}, nil
	}

	osbapiConfig := brokerClientConfig(broker)
	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	catalog, err := osbapiClient.GetCatalog()
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, catalogService := range catalog.Services {
		service, err := r.kubeServiceRepo.Create(broker, catalogService)
		if err != nil {
			return reconcile.Result{}, err
		}

		for _, catalogPlan := range catalogService.Plans {
			if err := r.kubePlanRepo.Create(service, catalogPlan); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	if err := r.kubeBrokerRepo.UpdateState(broker, osbapiv1alpha1.BrokerStateRegistered); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func brokerClientConfig(broker *osbapiv1alpha1.Broker) *osbapi.ClientConfiguration {
	osbapiConfig := osbapi.DefaultClientConfiguration()
	osbapiConfig.Name = broker.Spec.Name
	osbapiConfig.URL = broker.Spec.URL
	osbapiConfig.AuthConfig = &osbapi.AuthConfig{
		BasicAuthConfig: &osbapi.BasicAuthConfig{
			Username: broker.Spec.Username,
			Password: broker.Spec.Password,
		},
	}
	return osbapiConfig
}
