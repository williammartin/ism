package broker

import (
	"context"

	osbapiv1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var ctx = context.TODO()

//go:generate counterfeiter . KubeClient

type KubeClient interface {
	client.Client
}

//TODO Shall we leave this here or move it to the test?

//go:generate counterfeiter . KubeStatusWriter

type KubeStatusWriter interface {
	client.StatusWriter
}

//TODO Shall we leave this here or move it to the test?

//go:generate counterfeiter . BrokerClient

type BrokerClient interface {
	osbapi.Client
}

type BrokerReconciler struct {
	kubeClient         KubeClient
	createBrokerClient osbapi.CreateFunc
}

func NewBrokerReconciler(kubeClient KubeClient, createBrokerClient osbapi.CreateFunc) *BrokerReconciler {
	return &BrokerReconciler{
		kubeClient:         kubeClient,
		createBrokerClient: createBrokerClient,
	}
}

func (r *BrokerReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	//1. Fetch the broker resource
	broker := &osbapiv1alpha1.Broker{}
	if err := r.kubeClient.Get(ctx, request.NamespacedName, broker); err != nil {
		return reconcile.Result{}, err
	}

	//1.5 Check that the broker has not been "registered" already"
	if broker.Status.State == osbapiv1alpha1.BrokerStateRegistered {
		return reconcile.Result{}, nil
	}

	//2. Parse spec for broker details
	osbapiConfig := brokerClientConfig(broker)
	osbapiClient, err := r.createBrokerClient(osbapiConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	//3. Call the broker /v2/catalog
	catalog, err := osbapiClient.GetCatalog()
	if err != nil {
		return reconcile.Result{}, err
	}

	//4. For each service
	//4.1 Create a new service resource
	for _, catalogService := range catalog.Services {
		service := &osbapiv1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: broker.ObjectMeta.Namespace,
				Name:      broker.ObjectMeta.Name + "." + catalogService.ID,
			},
			Spec: osbapiv1alpha1.BrokerServiceSpec{
				Name:        catalogService.Name,
				Description: catalogService.Description,
			},
		}

		if err := r.kubeClient.Create(ctx, service); err != nil {
			return reconcile.Result{}, err
		}

		//4.2 For each plan
		//4.2.1 Create a new plan resource
		for _, catalogPlan := range catalogService.Plans {
			plan := &osbapiv1alpha1.BrokerServicePlan{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: broker.ObjectMeta.Namespace,
					Name:      broker.ObjectMeta.Name + "." + catalogService.ID + "." + catalogPlan.ID,
				},
				Spec: osbapiv1alpha1.BrokerServicePlanSpec{
					Name: catalogPlan.Name,
				},
			}
			if err := r.kubeClient.Create(ctx, plan); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	//5. Done. Report success in broker resource
	broker.Status.State = osbapiv1alpha1.BrokerStateRegistered

	if err := r.kubeClient.Status().Update(ctx, broker); err != nil {
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
