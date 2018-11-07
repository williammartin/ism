/*
Copyright 2018 The Kubernetes Authors.

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

package serviceattachment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	ismv1beta1 "github.com/pivotal-cf/ism-controller-fun/pkg/apis/ism/v1beta1"
	servicecatalogv1beta1 "github.com/pivotal-cf/ism-controller-fun/pkg/apis/servicecatalog/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ServiceAttachment Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this ism.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileServiceAttachment{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("serviceattachment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to ServiceAttachment
	err = c.Watch(&source.Kind{Type: &ismv1beta1.ServiceAttachment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by ServiceAttachment - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ismv1beta1.ServiceAttachment{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileServiceAttachment{}

// ReconcileServiceAttachment reconciles a ServiceAttachment object
type ReconcileServiceAttachment struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ServiceAttachment object and makes changes based on the state read
// and what is in the ServiceAttachment.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ism.k8s.io,resources=serviceattachments,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileServiceAttachment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ServiceAttachment instance
	attachment := &ismv1beta1.ServiceAttachment{}
	err := r.Get(context.TODO(), request.NamespacedName, attachment)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	binding := &servicecatalogv1beta1.ServiceBinding{}

	err = r.Get(context.TODO(), types.NamespacedName{Name: attachment.Spec.BindingName, Namespace: attachment.Namespace}, binding)
	if err != nil {
		fmt.Printf("OMG: %s", err.Error())
		return reconcile.Result{}, err
	}

	fmt.Printf("WORKED: %+v", binding)
	if binding.Status.Conditions[0].Status == "True" {
		secret := &corev1.Secret{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: binding.Spec.SecretName, Namespace: binding.Namespace}, secret)

		fmt.Printf("ZOMG: %+v", secret)

		if err := sendToCF(attachment.Spec.PlatformURL, attachment.Spec.Username, attachment.Spec.Password, attachment.Spec.AttachmentContext, secret.Data); err != nil {
			fmt.Printf("OOPS: %+v", err)
			return reconcile.Result{}, err
		}
		fmt.Print("WOOHOO")

	}
	return reconcile.Result{}, nil
}

func sendToCF(url, username, password, attachContext string, credentials map[string][]byte) error {
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

	bindingBody, err := createBindingRequest("moomins", attachContext, credentials)

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
