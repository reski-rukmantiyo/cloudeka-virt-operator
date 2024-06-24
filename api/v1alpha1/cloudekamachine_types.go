/*
Copyright 2024.

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

// CloudekaMachineSpec defines the desired state of CloudekaMachine
type CloudekaMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Processor int    `json:"processor"`
	Memory    string `json:"memory"`
	Disk      string `json:"disk"`
}

// CloudekaMachineStatus defines the observed state of CloudekaMachine
type CloudekaMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Running bool `json:"running"`
	Valid   bool `json:"valid"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CloudekaMachine is the Schema for the cloudekamachines API
type CloudekaMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudekaMachineSpec   `json:"spec,omitempty"`
	Status CloudekaMachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CloudekaMachineList contains a list of CloudekaMachine
type CloudekaMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudekaMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudekaMachine{}, &CloudekaMachineList{})
}
