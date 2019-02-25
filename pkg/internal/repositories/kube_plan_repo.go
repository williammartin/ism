package repositories

import (
	"context"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type kubePlanRepo struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewKubePlanRepo(client client.Client) *kubePlanRepo {
	return &kubePlanRepo{
		client: client,
		scheme: scheme.Scheme,
	}
}

func (repo *kubePlanRepo) Create(service *v1alpha1.BrokerService, catalogPlan osbapi.Plan) error {
	plan := &v1alpha1.BrokerServicePlan{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.ObjectMeta.Name + "." + catalogPlan.ID,
			Namespace: service.Namespace,
		},
		Spec: v1alpha1.BrokerServicePlanSpec{
			Name: catalogPlan.Name,
		},
	}

	if err := controllerutil.SetControllerReference(service, plan, repo.scheme); err != nil {
		return err
	}

	return repo.client.Create(context.TODO(), plan)
}
