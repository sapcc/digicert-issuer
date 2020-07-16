/*
Copyright 2020 SAP SE.

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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DigicertIssuerSpec defines the desired state of DigicertIssuer
type DigicertIssuerSpec struct {
	// URL is the DigiCert cert-central API.
	URL string `json:"url,omitempty"`

	// Provisioner contains the DigiCert provisioner configuration.
	Provisioner DigicertProvisioner `json:"provisioner"`
}

// DigiCertProvisioner contains the DigiCert provisioner configuration.
type DigicertProvisioner struct {
	// APITokenReference references a secret in the same namespace containing the DigiCert API token.
	APITokenReference SecretKeySelector `json:"apiTokenReference"`

	// CACertID ...
	CACertID string `json:"caCertId,omitempty"`

	// OrganizationID is the ID of the organization in Digicert.
	OrganizationID *int `json:"organizationID,omitempty"`

	// OrganizationName is the name of the organization in Digicert.
	// If specified takes precedence over OrganizationID.
	OrganizationName string `json:"organizationName,omitempty"`

	// OrganizationUnits ...
	OrganizationUnits []string `json:"organizationUnits,omitempty"`

	// ValidityYears ...
	ValidityYears *int `json:"validityYears,omitempty"`

	// DisableRenewalNotifications ...
	DisableRenewalNotifications *bool `json:"disableRenewalNotifications,omitempty"`

	// PaymentMethod ...
	PaymentMethod string `json:"paymentMethod,omitempty"`

	// SkipApproval ...
	SkipApproval *bool `json:"skipApproval,omitempty"`

	// OrderType ...
	OrderType string `json:"orderType,omitempty"`
}

// SecretKeySelector references a secret in the same namespace containing sensitive configuration.
type SecretKeySelector struct {
	// The name of the secret.
	Name string `json:"name"`

	// The key in the secret.
	Key string `json:"key"`
}

// DigicertIssuerStatus defines the observed state of DigicertIssuer
type DigicertIssuerStatus struct {
	// Conditions ...
	// +optional
	Conditions []DigicertIssuerCondition `json:"conditions,omitempty"`
}

// DigicertIssuerCondition  ...
type DigicertIssuerCondition struct {
	// Type of the condition, currently ('Ready').
	Type ConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	// +optional
	Reason ConditionReason `json:"reason,omitempty"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	// +optional
	Message string `json:"message,omitempty"`
}

// ConditionType represents a DigicertIssuer condition type.
// +kubebuilder:validation:Enum=Ready
type ConditionType string

const (
	// ConditionReady indicates that a DigicertIssuer is ready for use.
	ConditionReady ConditionType = "Ready"

	// ConditionConfigurationError indicates any configuration error.
	// See the condition message for details.
	ConditionConfigurationError ConditionType = "ConfigurationError"
)

// ConditionStatus represents a condition's status.
// +kubebuilder:validation:Enum=True;False;Unknown
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in
// the condition; "ConditionFalse" means a resource is not in the condition;
// "ConditionUnknown" means kubernetes can't decide if a resource is in the
// condition or not. In the future, we could add other intermediate
// conditions, e.g. ConditionDegraded.
const (
	// ConditionTrue represents the fact that a given condition is true
	ConditionTrue ConditionStatus = "True"

	// ConditionFalse represents the fact that a given condition is false
	ConditionFalse ConditionStatus = "False"

	// ConditionUnknown represents the fact that a given condition is unknown
	ConditionUnknown ConditionStatus = "Unknown"
)

type ConditionReason string

const (
	ConditionReasonInvalidIssuerSpec     ConditionReason = "InvalidIssuerSpec"
	ConditionReasonSecretNotFoundOrEmpty ConditionReason = "SecretNotFoundOrEmpty"
)

// +kubebuilder:object:root=true

// DigicertIssuer is the Schema for the digicertissuers API
// +kubebuilder:subresource:status
type DigicertIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DigicertIssuerSpec    `json:"spec,omitempty"`
	Status *DigicertIssuerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DigicertIssuerList contains a list of DigicertIssuer
type DigicertIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DigicertIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DigicertIssuer{}, &DigicertIssuerList{})
}
