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

// CloudekaDataVolumeSpec defines the desired state of CloudekaDataVolume
type CloudekaDataVolumeSpec struct {
	Type         string `json:"type"`
	Version      string `json:"version"`
	Size         string `json:"size"`
	StorageClass string `json:"storageClass"`
}

// CloudekaDataVolumeStatus defines the observed state of CloudekaDataVolume
type CloudekaDataVolumeStatus struct {
	Phase    string `json:"phase"`
	Progress string `json:"progress"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=cdv
//+kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
//+kubebuilder:printcolumn:name="Size",type="string",JSONPath=".spec.size"
//+kubebuilder:printcolumn:name="StorageClass",type="string",JSONPath=".spec.storageClass"
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// CloudekaDataVolume is the Schema for the cloudekadatavolumes API
type CloudekaDataVolume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudekaDataVolumeSpec   `json:"spec,omitempty"`
	Status CloudekaDataVolumeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CloudekaDataVolumeList contains a list of CloudekaDataVolume
type CloudekaDataVolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudekaDataVolume `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudekaDataVolume{}, &CloudekaDataVolumeList{})
}
