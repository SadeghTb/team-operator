/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TeamSpec defines the desired state of Team
type TeamSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ManagementState indicates whether and how the operator should manage the component.
	// Indicator if the resource is 'Managed' or 'Unmanaged' by the operator.
	ManagementState ManagementState `json:"managementState"`

	// Foo is an example field of Team. Edit team_types.go to remove/update
	TeamAdmin string `json:"teamAdmin,omitempty"`
	// Foo is an example field of Team. Edit team_types.go to remove/update
	Namespaces []string `json:"namespaces,omitempty"`
}

// TeamStatus defines the observed state of Team
type TeamStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:resource:scope=Cluster
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Team is the Schema for the teams API
type Team struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamSpec   `json:"spec,omitempty"`
	Status TeamStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TeamList contains a list of Team
type TeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Team `json:"items"`
}

// Managed means that the operator is actively managing its resources and trying to keep the component active.
// It will only upgrade the component if it is safe to do so
// Unmanaged means that the operator will not take any action related to the component
//
// +kubebuilder:validation:Enum:=Managed;Unmanaged

type ManagementState string

const (
	ManagementStateManaged   ManagementState = "Managed"
	ManagementStateUnmanaged ManagementState = "Unmanaged"
)

func init() {
	SchemeBuilder.Register(&Team{}, &TeamList{})
}
