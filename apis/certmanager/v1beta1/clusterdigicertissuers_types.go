// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and sapcc contributors
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const ClusterDigicertIssuerKind = "ClusterDigicertIssuer"

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster

// ClusterDigicertIssuer is the Schema for the clusterdigicertissuers API
type ClusterDigicertIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DigicertIssuerSpec    `json:"spec"`
	Status *DigicertIssuerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterDigicertIssuerList contains a list of ClusterDigicertIssuer
type ClusterDigicertIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterDigicertIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDigicertIssuer{}, &ClusterDigicertIssuerList{})
}
