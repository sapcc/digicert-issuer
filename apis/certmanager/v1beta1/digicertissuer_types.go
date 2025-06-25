// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DigicertIssuerSpec defines the desired state of DigicertIssuer
type DigicertIssuerSpec struct {
	// Optional URL is the DigiCert cert-central API.
	// +optional
	URL string `json:"url,omitempty"`

	// Provisioner contains the DigiCert provisioner configuration.
	Provisioner DigicertProvisioner `json:"provisioner"`
}

// +kubebuilder:validation:XValidation:message="only one of validityDays and validityYears can be set.",rule="has(self.validityDays) && !has(self.validityYears) || !has(self.validityDays) && has(self.validityYears)"

// DigiCertProvisioner contains the DigiCert provisioner configuration.
type DigicertProvisioner struct {
	// APITokenReference references a secret in the same namespace containing the DigiCert API token.
	APITokenReference SecretKeySelector `json:"apiTokenReference"`

	// CACertID is the ID of the CA if multiple CA certificates are configured in the (sub-)account.
	CACertID string `json:"caCertID,omitempty"`

	// OrganizationID is the ID of the organization in Digicert.
	OrganizationID *int `json:"organizationID,omitempty"`

	// OrganizationName is the name of the organization in Digicert.
	// If specified takes precedence over OrganizationID.
	OrganizationName string `json:"organizationName,omitempty"`

	// OrganizationUnits is the list of organizational units.
	OrganizationUnits []string `json:"organizationUnits,omitempty"`

	// ValidityDays is the validity of the order and certificate in days. Overrides ValidityYears if set.
	ValidityDays *int `json:"validityDays,omitempty"`

	// ValidityYears is the validity of the order and certificate in years. Defaults to 1 year if not set.
	// Can be overridden by ValidityDays.
	ValidityYears *int `json:"validityYears,omitempty"`

	// DisableRenewalNotifications disables email renewal notifications for expiring certificates.
	DisableRenewalNotifications *bool `json:"disableRenewalNotifications,omitempty"`

	// PaymentMethod is the configured payment method in the Digicert account.
	PaymentMethod string `json:"paymentMethod,omitempty"`

	// SkipApproval skips the approval of the certificate.
	SkipApproval *bool `json:"skipApproval,omitempty"`

	// OrderType is the certificate order type.
	OrderType string `json:"orderType,omitempty"`

	// ContainerID is the ID of the division
	ContainerID *int `json:"containerID,omitempty"`
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
	// Conditions is a list of DigicertIssuerConditions describing the current status.
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
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DigicertIssuer is the Schema for the digicertissuers API
type DigicertIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   DigicertIssuerSpec    `json:"spec"`
	Status *DigicertIssuerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DigicertIssuerList contains a list of DigicertIssuer
type DigicertIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DigicertIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DigicertIssuer{}, &DigicertIssuerList{})
}
