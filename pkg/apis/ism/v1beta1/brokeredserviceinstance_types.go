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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BrokeredServiceInstanceSpec defines the desired state of BrokeredServiceInstance
type BrokeredServiceInstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ServiceID string `json:"serviceId,omitempty"`
	PlanID    string `json:"planId,omitempty"`
	Name      string `json:"name,omitempty"`
	GUID      string `json:"guid,omitempty"`
	Migrated  bool   `json:"migrated,omitempty"`
}

// BrokeredServiceInstanceStatus defines the observed state of BrokeredServiceInstance
type BrokeredServiceInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Success bool `json:"success,omitempty"`
	Async   bool `json:"async,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokeredServiceInstance is the Schema for the brokeredserviceinstances API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type BrokeredServiceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokeredServiceInstanceSpec   `json:"spec,omitempty"`
	Status BrokeredServiceInstanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokeredServiceInstanceList contains a list of BrokeredServiceInstance
type BrokeredServiceInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrokeredServiceInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrokeredServiceInstance{}, &BrokeredServiceInstanceList{})
}
