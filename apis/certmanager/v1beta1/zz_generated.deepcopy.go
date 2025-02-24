// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and sapcc contributors
// SPDX-License-Identifier: Apache-2.0

//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DigicertIssuer) DeepCopyInto(out *DigicertIssuer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(DigicertIssuerStatus)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DigicertIssuer.
func (in *DigicertIssuer) DeepCopy() *DigicertIssuer {
	if in == nil {
		return nil
	}
	out := new(DigicertIssuer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DigicertIssuer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DigicertIssuerCondition) DeepCopyInto(out *DigicertIssuerCondition) {
	*out = *in
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DigicertIssuerCondition.
func (in *DigicertIssuerCondition) DeepCopy() *DigicertIssuerCondition {
	if in == nil {
		return nil
	}
	out := new(DigicertIssuerCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DigicertIssuerList) DeepCopyInto(out *DigicertIssuerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DigicertIssuer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DigicertIssuerList.
func (in *DigicertIssuerList) DeepCopy() *DigicertIssuerList {
	if in == nil {
		return nil
	}
	out := new(DigicertIssuerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DigicertIssuerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DigicertIssuerSpec) DeepCopyInto(out *DigicertIssuerSpec) {
	*out = *in
	in.Provisioner.DeepCopyInto(&out.Provisioner)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DigicertIssuerSpec.
func (in *DigicertIssuerSpec) DeepCopy() *DigicertIssuerSpec {
	if in == nil {
		return nil
	}
	out := new(DigicertIssuerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DigicertIssuerStatus) DeepCopyInto(out *DigicertIssuerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]DigicertIssuerCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DigicertIssuerStatus.
func (in *DigicertIssuerStatus) DeepCopy() *DigicertIssuerStatus {
	if in == nil {
		return nil
	}
	out := new(DigicertIssuerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DigicertProvisioner) DeepCopyInto(out *DigicertProvisioner) {
	*out = *in
	out.APITokenReference = in.APITokenReference
	if in.OrganizationID != nil {
		in, out := &in.OrganizationID, &out.OrganizationID
		*out = new(int)
		**out = **in
	}
	if in.OrganizationUnits != nil {
		in, out := &in.OrganizationUnits, &out.OrganizationUnits
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ValidityDays != nil {
		in, out := &in.ValidityDays, &out.ValidityDays
		*out = new(int)
		**out = **in
	}
	if in.ValidityYears != nil {
		in, out := &in.ValidityYears, &out.ValidityYears
		*out = new(int)
		**out = **in
	}
	if in.DisableRenewalNotifications != nil {
		in, out := &in.DisableRenewalNotifications, &out.DisableRenewalNotifications
		*out = new(bool)
		**out = **in
	}
	if in.SkipApproval != nil {
		in, out := &in.SkipApproval, &out.SkipApproval
		*out = new(bool)
		**out = **in
	}
	if in.ContainerID != nil {
		in, out := &in.ContainerID, &out.ContainerID
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DigicertProvisioner.
func (in *DigicertProvisioner) DeepCopy() *DigicertProvisioner {
	if in == nil {
		return nil
	}
	out := new(DigicertProvisioner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretKeySelector) DeepCopyInto(out *SecretKeySelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretKeySelector.
func (in *SecretKeySelector) DeepCopy() *SecretKeySelector {
	if in == nil {
		return nil
	}
	out := new(SecretKeySelector)
	in.DeepCopyInto(out)
	return out
}
