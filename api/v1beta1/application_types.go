/*
Copyright 2025 cuisongliu.

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

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name"`

	// Optional: no more than one of the following may be specified.
	// Defaults to "".
	// +optional
	Value string `json:"value,omitempty"`
}

// HelmVar represents an environment variable present in a Container.
type HelmVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name"`

	// Optional: no more than one of the following may be specified.
	// Defaults to "".
	// +optional
	Value string `json:"value,omitempty"`
	// Source for the environment variable's value. Cannot be used if value is not empty.
	// +optional
	ValueFrom *HelmVarSource `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

// HelmVarSource represents a source for the value of an HelmVar.
type HelmVarSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *v1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *v1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// ApplicationSpec defines the desired state of Application.
type ApplicationSpec struct {
	// Container image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// This field is optional to allow higher level config management to default or override
	// container images in workload controllers like Deployments and StatefulSets.
	// +optional
	Image string `json:"image,omitempty"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	Env []EnvVar `json:"env,omitempty"`
	// List of helm variables to set in the container.
	// Cannot be updated.
	// +optional
	Helm []HelmVar `json:"helm,omitempty"`
}

type ApplicationPhase string

// These are the valid phases of node.
const (
	ApplicationPending   ApplicationPhase = "Pending"
	ApplicationError     ApplicationPhase = "Error"
	ApplicationReady     ApplicationPhase = "Ready"
	ApplicationInProcess ApplicationPhase = "InProcess"
)

// ApplicationStatus defines the observed state of Application.
type ApplicationStatus struct {
	// Phase represents the current phase of Application.
	// +kubebuilder:default:=Unknown
	Phase ApplicationPhase `json:"phase,omitempty"`
	// Conditions contains the different condition statuses for this Application.
	// +optional
	Conditions []metav1.Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Application is the Schema for the applications API.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
