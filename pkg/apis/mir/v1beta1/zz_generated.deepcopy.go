// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CanarySpec) DeepCopyInto(out *CanarySpec) {
	*out = *in
	in.ModelSpec.DeepCopyInto(&out.ModelSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CanarySpec.
func (in *CanarySpec) DeepCopy() *CanarySpec {
	if in == nil {
		return nil
	}
	out := new(CanarySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomSpec) DeepCopyInto(out *CustomSpec) {
	*out = *in
	in.Container.DeepCopyInto(&out.Container)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomSpec.
func (in *CustomSpec) DeepCopy() *CustomSpec {
	if in == nil {
		return nil
	}
	out := new(CustomSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentSource) DeepCopyInto(out *ModelDeploymentSource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentSource.
func (in *ModelDeploymentSource) DeepCopy() *ModelDeploymentSource {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelDeploymentSource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentSourceList) DeepCopyInto(out *ModelDeploymentSourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelDeploymentSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentSourceList.
func (in *ModelDeploymentSourceList) DeepCopy() *ModelDeploymentSourceList {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentSourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelDeploymentSourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentSourceSpec) DeepCopyInto(out *ModelDeploymentSourceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentSourceSpec.
func (in *ModelDeploymentSourceSpec) DeepCopy() *ModelDeploymentSourceSpec {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentSourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelDeploymentSourceStatus) DeepCopyInto(out *ModelDeploymentSourceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelDeploymentSourceStatus.
func (in *ModelDeploymentSourceStatus) DeepCopy() *ModelDeploymentSourceStatus {
	if in == nil {
		return nil
	}
	out := new(ModelDeploymentSourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelInferenceResource) DeepCopyInto(out *ModelInferenceResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelInferenceResource.
func (in *ModelInferenceResource) DeepCopy() *ModelInferenceResource {
	if in == nil {
		return nil
	}
	out := new(ModelInferenceResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelInferenceResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelInferenceResourceList) DeepCopyInto(out *ModelInferenceResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelInferenceResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelInferenceResourceList.
func (in *ModelInferenceResourceList) DeepCopy() *ModelInferenceResourceList {
	if in == nil {
		return nil
	}
	out := new(ModelInferenceResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelInferenceResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelInferenceResourceSpec) DeepCopyInto(out *ModelInferenceResourceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelInferenceResourceSpec.
func (in *ModelInferenceResourceSpec) DeepCopy() *ModelInferenceResourceSpec {
	if in == nil {
		return nil
	}
	out := new(ModelInferenceResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelInferenceResourceStatus) DeepCopyInto(out *ModelInferenceResourceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelInferenceResourceStatus.
func (in *ModelInferenceResourceStatus) DeepCopy() *ModelInferenceResourceStatus {
	if in == nil {
		return nil
	}
	out := new(ModelInferenceResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelService) DeepCopyInto(out *ModelService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelService.
func (in *ModelService) DeepCopy() *ModelService {
	if in == nil {
		return nil
	}
	out := new(ModelService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServiceList) DeepCopyInto(out *ModelServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServiceList.
func (in *ModelServiceList) DeepCopy() *ModelServiceList {
	if in == nil {
		return nil
	}
	out := new(ModelServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServiceSpec) DeepCopyInto(out *ModelServiceSpec) {
	*out = *in
	in.Default.DeepCopyInto(&out.Default)
	if in.Canary != nil {
		in, out := &in.Canary, &out.Canary
		*out = new(CanarySpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServiceSpec.
func (in *ModelServiceSpec) DeepCopy() *ModelServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ModelServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServiceStatus) DeepCopyInto(out *ModelServiceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServiceStatus.
func (in *ModelServiceStatus) DeepCopy() *ModelServiceStatus {
	if in == nil {
		return nil
	}
	out := new(ModelServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelSpec) DeepCopyInto(out *ModelSpec) {
	*out = *in
	if in.Custom != nil {
		in, out := &in.Custom, &out.Custom
		*out = new(CustomSpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelSpec.
func (in *ModelSpec) DeepCopy() *ModelSpec {
	if in == nil {
		return nil
	}
	out := new(ModelSpec)
	in.DeepCopyInto(out)
	return out
}
