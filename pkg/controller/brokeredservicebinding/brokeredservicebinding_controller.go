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

package brokeredservicebinding

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"reflect"

	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
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

// Add creates a new BrokeredServiceBinding Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this ism.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBrokeredServiceBinding{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("brokeredservicebinding-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to BrokeredServiceBinding
	err = c.Watch(&source.Kind{Type: &ismv1beta1.BrokeredServiceBinding{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileBrokeredServiceBinding{}

// ReconcileBrokeredServiceBinding reconciles a BrokeredServiceBinding object
type ReconcileBrokeredServiceBinding struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a BrokeredServiceBinding object and makes changes based on the state read
// and what is in the BrokeredServiceBinding.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ism.ism.pivotal.io,resources=brokeredservicebindings,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileBrokeredServiceBinding) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("Reconile %s\n", request.NamespacedName)
	// Fetch the BrokeredServiceBinding instance
	binding := &ismv1beta1.BrokeredServiceBinding{}
	err := r.Get(context.TODO(), request.NamespacedName, binding)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if binding.Status.Success {
		fmt.Println("binding already created & attached")
		return reconcile.Result{}, nil
	}

	instance := &ismv1beta1.BrokeredServiceInstance{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: binding.Namespace, Name: binding.Spec.ServiceInstanceGUID}, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.

			//TODO: Set status to failed
			fmt.Println("instance not found")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if binding.Spec.Migrated {
		fmt.Println("binding migrated, skipping creation")
		binding.Status.Success = true
		binding.Status.Credentials = binding.Spec.MigratedCredentials

		if err := controllerutil.SetControllerReference(instance, binding, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		if err := r.Update(context.TODO(), binding); err != nil {
			return reconcile.Result{}, err
		}

		if err := r.Status().Update(context.TODO(), binding); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, err
	}

	//broker
	broker := &ismv1beta1.Broker{}
	err = r.Get(context.TODO(), types.NamespacedName{Namespace: binding.Namespace, Name: instance.OwnerReferences[0].Name}, broker)
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

	resp, err := bind(broker.Spec.URL, broker.Spec.Username, broker.Spec.Password, instance.Spec.ServiceID, instance.Spec.PlanID, instance.Spec.GUID, string(binding.GetUID()))
	if err != nil {
		fmt.Println("binding failed")
		return reconcile.Result{}, err
	}

	beforeBinding := binding.DeepCopy()

	binding.Status.Success = true

	s, err := json.Marshal(resp.Credentials)
	if err != nil {
		return reconcile.Result{}, err
	}
	binding.Status.Credentials = string(s)

	if err := controllerutil.SetControllerReference(broker, binding, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(beforeBinding, binding) {
		if err := r.Update(context.TODO(), binding); err != nil {
			return reconcile.Result{}, err
		}

		if err := r.Status().Update(context.TODO(), binding); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func bind(url, username, password, serviceID, planID, instanceID, bindingID string) (*osb.BindResponse, error) {
	config := osb.DefaultClientConfiguration()
	config.URL = url
	config.Insecure = true
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

	req := osb.BindRequest{
		InstanceID: instanceID,
		BindingID:  bindingID,
		ServiceID:  serviceID,
		PlanID:     planID,
	}

	resp, err := client.Bind(&req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func sendToCF(url, username, password, attachContext, bindingID string, credentials map[string]interface{}) error {
	fmt.Println(url)
	c := &cfclient.Config{
		ApiAddress:        url,
		Username:          username,
		Password:          password,
		SkipSslValidation: true,
	}

	client, err := cfclient.NewClient(c)

	if err != nil {
		return err
	}

	bindingBody, err := createBindingRequest(bindingID, attachContext, credentials)

	if err != nil {
		return err
	}

	req := client.NewRequestWithBody("POST", "/v3/external_service_bindings", bytes.NewBuffer(bindingBody))

	resp, err := client.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%+v", body)
		return err
	}
	return nil
}

func createBindingRequest(bindingID string, appGUID string, credentials interface{}) ([]byte, error) {
	b, err := json.Marshal(credentials)

	if err != nil {
		return nil, err
	}

	return []byte(`{
    "type": "app",
		"name": "` + bindingID + `",
    "relationships": {
      "app": {
        "data": {
          "guid": "` + appGUID + `"
         }
       }
    },
    "data": {
      "credentials": ` + string(b) + `
    }
  }`), nil
}
