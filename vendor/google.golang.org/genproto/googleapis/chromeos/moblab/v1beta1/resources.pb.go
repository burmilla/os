// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.2
// source: google/chromeos/moblab/v1beta1/resources.proto

package moblab

import (
	reflect "reflect"
	sync "sync"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The build status types.
type Build_BuildStatus int32

const (
	// No build status is specified.
	Build_BUILD_STATUS_UNSPECIFIED Build_BuildStatus = 0
	// Complete Status: The build passed.
	Build_PASS Build_BuildStatus = 1
	// Complete Status: The build failed.
	Build_FAIL Build_BuildStatus = 2
	// Intermediate Status: The build is still running.
	Build_RUNNING Build_BuildStatus = 3
	// Complete Status: The build was aborted.
	Build_ABORTED Build_BuildStatus = 4
)

// Enum value maps for Build_BuildStatus.
var (
	Build_BuildStatus_name = map[int32]string{
		0: "BUILD_STATUS_UNSPECIFIED",
		1: "PASS",
		2: "FAIL",
		3: "RUNNING",
		4: "ABORTED",
	}
	Build_BuildStatus_value = map[string]int32{
		"BUILD_STATUS_UNSPECIFIED": 0,
		"PASS":                     1,
		"FAIL":                     2,
		"RUNNING":                  3,
		"ABORTED":                  4,
	}
)

func (x Build_BuildStatus) Enum() *Build_BuildStatus {
	p := new(Build_BuildStatus)
	*p = x
	return p
}

func (x Build_BuildStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Build_BuildStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_google_chromeos_moblab_v1beta1_resources_proto_enumTypes[0].Descriptor()
}

func (Build_BuildStatus) Type() protoreflect.EnumType {
	return &file_google_chromeos_moblab_v1beta1_resources_proto_enumTypes[0]
}

func (x Build_BuildStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Build_BuildStatus.Descriptor instead.
func (Build_BuildStatus) EnumDescriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{3, 0}
}

// The build types.
type Build_BuildType int32

const (
	// Invalid build type.
	Build_BUILD_TYPE_UNSPECIFIED Build_BuildType = 0
	// The release build.
	Build_RELEASE Build_BuildType = 1
	// The firmware build.
	Build_FIRMWARE Build_BuildType = 2
)

// Enum value maps for Build_BuildType.
var (
	Build_BuildType_name = map[int32]string{
		0: "BUILD_TYPE_UNSPECIFIED",
		1: "RELEASE",
		2: "FIRMWARE",
	}
	Build_BuildType_value = map[string]int32{
		"BUILD_TYPE_UNSPECIFIED": 0,
		"RELEASE":                1,
		"FIRMWARE":               2,
	}
)

func (x Build_BuildType) Enum() *Build_BuildType {
	p := new(Build_BuildType)
	*p = x
	return p
}

func (x Build_BuildType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Build_BuildType) Descriptor() protoreflect.EnumDescriptor {
	return file_google_chromeos_moblab_v1beta1_resources_proto_enumTypes[1].Descriptor()
}

func (Build_BuildType) Type() protoreflect.EnumType {
	return &file_google_chromeos_moblab_v1beta1_resources_proto_enumTypes[1]
}

func (x Build_BuildType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Build_BuildType.Descriptor instead.
func (Build_BuildType) EnumDescriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{3, 1}
}

// Resource that represents a build target.
type BuildTarget struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the build target.
	// Format: buildTargets/{build_target}
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *BuildTarget) Reset() {
	*x = BuildTarget{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildTarget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildTarget) ProtoMessage() {}

func (x *BuildTarget) ProtoReflect() protoreflect.Message {
	mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildTarget.ProtoReflect.Descriptor instead.
func (*BuildTarget) Descriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{0}
}

func (x *BuildTarget) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Resource that represents a model. Each model belongs to a build target. For
// non-unified build, the model name is the same as its build target name.
type Model struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the model.
	// Format: buildTargets/{build_target}/models/{model}
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Model) Reset() {
	*x = Model{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model) ProtoMessage() {}

func (x *Model) ProtoReflect() protoreflect.Message {
	mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Model.ProtoReflect.Descriptor instead.
func (*Model) Descriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{1}
}

func (x *Model) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Resource that represents a chrome OS milestone.
type Milestone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the milestone.
	// Format: milestones/{milestone}
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Milestone) Reset() {
	*x = Milestone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Milestone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Milestone) ProtoMessage() {}

func (x *Milestone) ProtoReflect() protoreflect.Message {
	mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Milestone.ProtoReflect.Descriptor instead.
func (*Milestone) Descriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{2}
}

func (x *Milestone) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Resource that represents a build for the given build target, model, milestone
// and build version.
type Build struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the build.
	// Format: buildTargets/{build_target}/models/{model}/builds/{build}
	// Example: buildTargets/octopus/models/bobba/builds/1234.0.0
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The milestone that owns the build.
	// Format: milestones/{milestone}
	Milestone string `protobuf:"bytes,2,opt,name=milestone,proto3" json:"milestone,omitempty"`
	// The build version of the build, e.g. 1234.0.0.
	BuildVersion string `protobuf:"bytes,3,opt,name=build_version,json=buildVersion,proto3" json:"build_version,omitempty"`
	// The status of the build.
	Status Build_BuildStatus `protobuf:"varint,4,opt,name=status,proto3,enum=google.chromeos.moblab.v1beta1.Build_BuildStatus" json:"status,omitempty"`
	// The type of the build.
	Type Build_BuildType `protobuf:"varint,5,opt,name=type,proto3,enum=google.chromeos.moblab.v1beta1.Build_BuildType" json:"type,omitempty"`
	// The branch of the build.
	Branch string `protobuf:"bytes,6,opt,name=branch,proto3" json:"branch,omitempty"`
	// The read write firmware version of the software that is flashed to the chip
	// on the Chrome OS device.
	RwFirmwareVersion string `protobuf:"bytes,7,opt,name=rw_firmware_version,json=rwFirmwareVersion,proto3" json:"rw_firmware_version,omitempty"`
}

func (x *Build) Reset() {
	*x = Build{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Build) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Build) ProtoMessage() {}

func (x *Build) ProtoReflect() protoreflect.Message {
	mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Build.ProtoReflect.Descriptor instead.
func (*Build) Descriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{3}
}

func (x *Build) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Build) GetMilestone() string {
	if x != nil {
		return x.Milestone
	}
	return ""
}

func (x *Build) GetBuildVersion() string {
	if x != nil {
		return x.BuildVersion
	}
	return ""
}

func (x *Build) GetStatus() Build_BuildStatus {
	if x != nil {
		return x.Status
	}
	return Build_BUILD_STATUS_UNSPECIFIED
}

func (x *Build) GetType() Build_BuildType {
	if x != nil {
		return x.Type
	}
	return Build_BUILD_TYPE_UNSPECIFIED
}

func (x *Build) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *Build) GetRwFirmwareVersion() string {
	if x != nil {
		return x.RwFirmwareVersion
	}
	return ""
}

// Resource that represents a build artifact stored in Google Cloud Storage for
// the given build target, model, build version and bucket.
type BuildArtifact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The resource name of the build artifact.
	// Format:
	// buildTargets/{build_target}/models/{model}/builds/{build}/artifacts/{artifact}
	// Example:
	// buildTargets/octopus/models/bobba/builds/1234.0.0/artifacts/chromeos-moblab-peng-staging
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The build metadata of the build artifact.
	Build string `protobuf:"bytes,2,opt,name=build,proto3" json:"build,omitempty"`
	// The bucket that stores the build artifact.
	Bucket string `protobuf:"bytes,3,opt,name=bucket,proto3" json:"bucket,omitempty"`
	// The path of the build artifact in the bucket.
	Path string `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	// The number of objects in the build artifact folder. The object number can
	// be used to calculated the stage progress by comparing the source build
	// artifact with the destination build artifact.
	ObjectCount uint32 `protobuf:"varint,5,opt,name=object_count,json=objectCount,proto3" json:"object_count,omitempty"`
}

func (x *BuildArtifact) Reset() {
	*x = BuildArtifact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildArtifact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildArtifact) ProtoMessage() {}

func (x *BuildArtifact) ProtoReflect() protoreflect.Message {
	mi := &file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildArtifact.ProtoReflect.Descriptor instead.
func (*BuildArtifact) Descriptor() ([]byte, []int) {
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP(), []int{4}
}

func (x *BuildArtifact) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BuildArtifact) GetBuild() string {
	if x != nil {
		return x.Build
	}
	return ""
}

func (x *BuildArtifact) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *BuildArtifact) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *BuildArtifact) GetObjectCount() uint32 {
	if x != nil {
		return x.ObjectCount
	}
	return 0
}

var File_google_chromeos_moblab_v1beta1_resources_proto protoreflect.FileDescriptor

var file_google_chromeos_moblab_v1beta1_resources_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2f, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2e, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6e, 0x0a, 0x0b, 0x42,
	0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x4b,
	0xea, 0x41, 0x48, 0x0a, 0x29, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f, 0x62,
	0x6c, 0x61, 0x62, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x1b,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x2f, 0x7b, 0x62, 0x75,
	0x69, 0x6c, 0x64, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x7d, 0x22, 0x71, 0x0a, 0x05, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x54, 0xea, 0x41, 0x51, 0x0a, 0x23, 0x63,
	0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x6f, 0x64,
	0x65, 0x6c, 0x12, 0x2a, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73,
	0x2f, 0x7b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x7d, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x7b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x7d, 0x22, 0x65,
	0x0a, 0x09, 0x4d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a,
	0x44, 0xea, 0x41, 0x41, 0x0a, 0x27, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f,
	0x62, 0x6c, 0x61, 0x62, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12, 0x16, 0x6d,
	0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x73, 0x2f, 0x7b, 0x6d, 0x69, 0x6c, 0x65, 0x73,
	0x74, 0x6f, 0x6e, 0x65, 0x7d, 0x22, 0xe8, 0x04, 0x0a, 0x05, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x4a, 0x0a, 0x09, 0x6d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2c, 0xfa, 0x41, 0x29, 0x0a, 0x27, 0x63, 0x68, 0x72,
	0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x6c, 0x65, 0x73,
	0x74, 0x6f, 0x6e, 0x65, 0x52, 0x09, 0x6d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x49, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x31, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x42, 0x75, 0x69, 0x6c,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x43, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e,
	0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x42,
	0x75, 0x69, 0x6c, 0x64, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x2e, 0x0a, 0x13,
	0x72, 0x77, 0x5f, 0x66, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x77, 0x46, 0x69, 0x72,
	0x6d, 0x77, 0x61, 0x72, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x59, 0x0a, 0x0b,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a, 0x18, 0x42,
	0x55, 0x49, 0x4c, 0x44, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x41, 0x53,
	0x53, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x02, 0x12, 0x0b, 0x0a,
	0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x42,
	0x4f, 0x52, 0x54, 0x45, 0x44, 0x10, 0x04, 0x22, 0x42, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x4c, 0x45, 0x41, 0x53, 0x45, 0x10, 0x01, 0x12, 0x0c, 0x0a,
	0x08, 0x46, 0x49, 0x52, 0x4d, 0x57, 0x41, 0x52, 0x45, 0x10, 0x02, 0x3a, 0x63, 0xea, 0x41, 0x60,
	0x0a, 0x23, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x39, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x73, 0x2f, 0x7b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x7d, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x7b, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x7d, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x73, 0x2f, 0x7b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x7d,
	0x22, 0xb5, 0x02, 0x0a, 0x0d, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61,
	0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x05, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x28, 0xfa, 0x41, 0x25, 0x0a, 0x23, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x52,
	0x05, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x3a, 0x80, 0x01, 0xea, 0x41, 0x7d, 0x0a, 0x2b, 0x63, 0x68, 0x72,
	0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x12, 0x4e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x2f, 0x7b, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x7d, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x7b, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x7d, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x73, 0x2f, 0x7b, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x7d, 0x2f, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73, 0x2f, 0x7b, 0x61,
	0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x7d, 0x42, 0x7e, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e,
	0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42, 0x0e,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x01,
	0x50, 0x01, 0x5a, 0x44, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e,
	0x67, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65,
	0x6f, 0x73, 0x2f, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x3b, 0x6d, 0x6f, 0x62, 0x6c, 0x61, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_google_chromeos_moblab_v1beta1_resources_proto_rawDescOnce sync.Once
	file_google_chromeos_moblab_v1beta1_resources_proto_rawDescData = file_google_chromeos_moblab_v1beta1_resources_proto_rawDesc
)

func file_google_chromeos_moblab_v1beta1_resources_proto_rawDescGZIP() []byte {
	file_google_chromeos_moblab_v1beta1_resources_proto_rawDescOnce.Do(func() {
		file_google_chromeos_moblab_v1beta1_resources_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_chromeos_moblab_v1beta1_resources_proto_rawDescData)
	})
	return file_google_chromeos_moblab_v1beta1_resources_proto_rawDescData
}

var file_google_chromeos_moblab_v1beta1_resources_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_google_chromeos_moblab_v1beta1_resources_proto_goTypes = []interface{}{
	(Build_BuildStatus)(0), // 0: google.chromeos.moblab.v1beta1.Build.BuildStatus
	(Build_BuildType)(0),   // 1: google.chromeos.moblab.v1beta1.Build.BuildType
	(*BuildTarget)(nil),    // 2: google.chromeos.moblab.v1beta1.BuildTarget
	(*Model)(nil),          // 3: google.chromeos.moblab.v1beta1.Model
	(*Milestone)(nil),      // 4: google.chromeos.moblab.v1beta1.Milestone
	(*Build)(nil),          // 5: google.chromeos.moblab.v1beta1.Build
	(*BuildArtifact)(nil),  // 6: google.chromeos.moblab.v1beta1.BuildArtifact
}
var file_google_chromeos_moblab_v1beta1_resources_proto_depIdxs = []int32{
	0, // 0: google.chromeos.moblab.v1beta1.Build.status:type_name -> google.chromeos.moblab.v1beta1.Build.BuildStatus
	1, // 1: google.chromeos.moblab.v1beta1.Build.type:type_name -> google.chromeos.moblab.v1beta1.Build.BuildType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_google_chromeos_moblab_v1beta1_resources_proto_init() }
func file_google_chromeos_moblab_v1beta1_resources_proto_init() {
	if File_google_chromeos_moblab_v1beta1_resources_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildTarget); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Model); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Milestone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Build); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildArtifact); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_chromeos_moblab_v1beta1_resources_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_chromeos_moblab_v1beta1_resources_proto_goTypes,
		DependencyIndexes: file_google_chromeos_moblab_v1beta1_resources_proto_depIdxs,
		EnumInfos:         file_google_chromeos_moblab_v1beta1_resources_proto_enumTypes,
		MessageInfos:      file_google_chromeos_moblab_v1beta1_resources_proto_msgTypes,
	}.Build()
	File_google_chromeos_moblab_v1beta1_resources_proto = out.File
	file_google_chromeos_moblab_v1beta1_resources_proto_rawDesc = nil
	file_google_chromeos_moblab_v1beta1_resources_proto_goTypes = nil
	file_google_chromeos_moblab_v1beta1_resources_proto_depIdxs = nil
}
