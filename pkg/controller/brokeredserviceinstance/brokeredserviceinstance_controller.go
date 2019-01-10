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

package brokeredserviceinstance

import (
	"context"
	"fmt"
	"reflect"

	ismv1beta1 "github.com/pivotal-cf/ism/pkg/apis/ism/v1beta1"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
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
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new BrokeredServiceInstance Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this ism.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBrokeredServiceInstance{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("brokeredserviceinstance-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to BrokeredServiceInstance
	err = c.Watch(&source.Kind{Type: &ismv1beta1.BrokeredServiceInstance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileBrokeredServiceInstance{}

// ReconcileBrokeredServiceInstance reconciles a BrokeredServiceInstance object
type ReconcileBrokeredServiceInstance struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a BrokeredServiceInstance object and makes changes based on the state read
// and what is in the BrokeredServiceInstance.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ism.ism.pivotal.io,resources=brokeredserviceinstances,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileBrokeredServiceInstance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("Reconcile %s\n", request.NamespacedName)
	// Fetch the BrokeredServiceInstance instance
	instance := &ismv1beta1.BrokeredServiceInstance{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.

			fmt.Println("instance not found")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.Status.Success {
		fmt.Println("instance already created")
		return reconcile.Result{}, err
	}

	//service
	service := &ismv1beta1.BrokeredService{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: instance.Spec.ServiceID}, service)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.

			//TODO: Set status to failed
			fmt.Println("service not found")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//broker
	broker := &ismv1beta1.Broker{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: service.OwnerReferences[0].Name}, broker)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.

			//TODO: Set status to failed
			fmt.Println("broker not found")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	beforeInstance := instance.DeepCopy()

	provisionResp, err := provision(broker.Spec.URL, broker.Spec.Username, broker.Spec.Password, instance.Spec.PlanID, service.Spec.ID, instance.Spec.GUID)
	if err != nil {
		return reconcile.Result{}, err
	}

	instance.Status.Async = provisionResp.Async
	instance.Status.Success = true

	if err := controllerutil.SetControllerReference(broker, instance, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(beforeInstance, instance) {
		if err := r.Status().Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func provision(url, username, password, planID, serviceID, instanceID string) (*osb.ProvisionResponse, error) {
	config := osb.DefaultClientConfiguration()
	config.URL = url
	basicAuthConfig := osb.AuthConfig{
		BasicAuthConfig: &osb.BasicAuthConfig{
			Username: username,
			Password: password,
		},
	}
	config.AuthConfig = &basicAuthConfig

	client, err := osb.NewClient(config)
	if err != nil {
		return nil, err
	}

	req := osb.ProvisionRequest{
		InstanceID:        instanceID,
		PlanID:            planID,
		ServiceID:         serviceID,
		OrganizationGUID:  "o",
		SpaceGUID:         "s",
		AcceptsIncomplete: false,
	}

	resp, err := client.ProvisionInstance(&req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
