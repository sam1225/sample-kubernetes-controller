//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResource) DeepCopyInto(out *MyAppResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResource.
func (in *MyAppResource) DeepCopy() *MyAppResource {
	if in == nil {
		return nil
	}
	out := new(MyAppResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MyAppResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResourceList) DeepCopyInto(out *MyAppResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MyAppResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResourceList.
func (in *MyAppResourceList) DeepCopy() *MyAppResourceList {
	if in == nil {
		return nil
	}
	out := new(MyAppResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MyAppResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResourceSpec) DeepCopyInto(out *MyAppResourceSpec) {
	*out = *in
	out.Resources = in.Resources
	out.Image = in.Image
	out.UI = in.UI
	out.Redis = in.Redis
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResourceSpec.
func (in *MyAppResourceSpec) DeepCopy() *MyAppResourceSpec {
	if in == nil {
		return nil
	}
	out := new(MyAppResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResourceStatus) DeepCopyInto(out *MyAppResourceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResourceStatus.
func (in *MyAppResourceStatus) DeepCopy() *MyAppResourceStatus {
	if in == nil {
		return nil
	}
	out := new(MyAppResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodImage) DeepCopyInto(out *PodImage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodImage.
func (in *PodImage) DeepCopy() *PodImage {
	if in == nil {
		return nil
	}
	out := new(PodImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodResources) DeepCopyInto(out *PodResources) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodResources.
func (in *PodResources) DeepCopy() *PodResources {
	if in == nil {
		return nil
	}
	out := new(PodResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodUI) DeepCopyInto(out *PodUI) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodUI.
func (in *PodUI) DeepCopy() *PodUI {
	if in == nil {
		return nil
	}
	out := new(PodUI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RedisFlag) DeepCopyInto(out *RedisFlag) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedisFlag.
func (in *RedisFlag) DeepCopy() *RedisFlag {
	if in == nil {
		return nil
	}
	out := new(RedisFlag)
	in.DeepCopyInto(out)
	return out
}
