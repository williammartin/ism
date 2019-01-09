/*
Copyright 2018 The ISM Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package broker

import (
	"context"
	"fmt"

	ismv1beta1 "github.com/pivotal-cf/ism/pkg/apis/ism/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Broker Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this ism.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBroker{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("broker-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Broker
	err = c.Watch(&source.Kind{Type: &ismv1beta1.Broker{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &ismv1beta1.BrokeredService{}}, &handler.EnqueueRequestForOwner{
		IsController: false,
		OwnerType:    &ismv1beta1.Broker{},
	})
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileBroker{}

// ReconcileBroker reconciles a Broker object
type ReconcileBroker struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Broker object and makes changes based on the state read
// and what is in the Broker.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=ism.ism.pivotal.io,resources=brokers,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileBroker) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Broker instance
	broker := &ismv1beta1.Broker{}
	err := r.Get(context.TODO(), request.NamespacedName, broker)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// err := r.List(context.TODO(),

	fmt.Printf("RECONCILE: %#v\n", broker)

	cat, err := getBrokerCatalog(broker)
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, service := range cat.Services {
		brokeredService := &ismv1beta1.BrokeredService{}
		err := r.Get(context.TODO(), types.NamespacedName{Namespace: broker.Namespace, Name: service.ID}, brokeredService)

		if errors.IsNotFound(err) {
			brokeredServiceSpec := ismv1beta1.BrokeredServiceSpec{
				Name:        service.Name,
				ID:          service.ID,
				Bindable:    service.Bindable,
				Description: service.Description,
			}
			brokeredService = &ismv1beta1.BrokeredService{Spec: brokeredServiceSpec}
			brokeredService.Name = service.ID
			brokeredService.Namespace = request.Namespace
			brokeredService.Labels = map[string]string{
				"ServiceName": service.Name,
			}

			if err := controllerutil.SetControllerReference(broker, brokeredService, r.scheme); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.Create(context.TODO(), brokeredService); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			fmt.Printf("service %s already exists", service.Name)
		}

		for _, plan := range service.Plans {
			brokeredPlan := &ismv1beta1.BrokeredServicePlan{}
			err := r.Get(context.TODO(), types.NamespacedName{Namespace: broker.Namespace, Name: plan.ID}, brokeredPlan)
			if !errors.IsNotFound(err) {
				fmt.Printf("plan %s already exists", plan.Name)
				break
			}

			brokeredPlansSpec := ismv1beta1.BrokeredServicePlanSpec{
				Name:        plan.Name,
				ID:          plan.ID,
				Description: plan.Description,
			}
			brokeredPlan = &ismv1beta1.BrokeredServicePlan{Spec: brokeredPlansSpec}
			brokeredPlan.Name = plan.ID
			brokeredPlan.Namespace = request.Namespace
			brokeredPlan.Labels = map[string]string{
				"PlanName": brokeredPlan.Name,
			}

			if err := controllerutil.SetControllerReference(brokeredService, brokeredPlan, r.scheme); err != nil {
				return reconcile.Result{}, err
			}

			if err := r.Create(context.TODO(), brokeredPlan); err != nil {
				return reconcile.Result{}, err
			}
		}

	}

	fmt.Printf("CATALOG: %#v\n", cat)

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	// if !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Printf("Updating Deployment %s/%s\n", deploy.Namespace, deploy.Name)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }

	return reconcile.Result{}, nil
}

func getBrokerCatalog(instance *ismv1beta1.Broker) (*osb.CatalogResponse, error) {
	config := osb.DefaultClientConfiguration()
	config.URL = instance.Spec.URL
	basicAuthConfig := osb.AuthConfig{
		BasicAuthConfig: &osb.BasicAuthConfig{
			Username: instance.Spec.Username,
			Password: instance.Spec.Password,
		},
	}
	config.AuthConfig = &basicAuthConfig

	client, err := osb.NewClient(config)
	if err != nil {
		return nil, err
	}

	cat, err := client.GetCatalog()
	fmt.Printf("CAT: %+v\n", cat)
	if err != nil {
		return nil, err
	}

	return cat, nil
}
