/*
Copyright 2019 Microsoft Corporation.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TODO: just use KFService and kill this
// Copied from https://github.com/kubeflow/kfserving/blob/master/pkg/apis/serving/v1alpha1/kfservice_types.go

// ModelServiceSpec defines the desired state of ModelService
type ModelServiceSpec struct {
	MinReplicas int32 `json:"minReplicas"`
	MaxReplicas int32 `json:"maxReplicas"`

	Default ModelSpec `json:"default"`
	// Optional Canary definition
	Canary *CanarySpec `json:"canary,omitempty"`
}

type ModelSpec struct {
	// The following fields follow a "1-of" semantic. Users must specify exactly one spec.
	Custom *CustomSpec `json:"custom,omitempty"`
}

// CanarySpec defines an alternate configuration to route a percentage of traffic.
type CanarySpec struct {
	ModelSpec      `json:",inline"`
	TrafficPercent int32 `json:"trafficPercent"`
}

type CustomSpec struct {
	Container v1.Container `json:"container"`
}

// ModelServiceStatus defines the observed state of ModelService
type ModelServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ModelService is the Schema for the modelservices API
// +k8s:openapi-gen=true
type ModelService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelServiceSpec   `json:"spec,omitempty"`
	Status ModelServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ModelServiceList contains a list of ModelService
type ModelServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelService{}, &ModelServiceList{})
}
