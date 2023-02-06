// Copyright 2020 Google LLC
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
// source: google/actions/sdk/v2/interactionmodel/prompt/content/static_content_prompt.proto

package prompt

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A placeholder for the Content part of a StaticPrompt.
type StaticContentPrompt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Only one type of content can be present in a Prompt.
	//
	// Types that are assignable to Content:
	//	*StaticContentPrompt_Card
	//	*StaticContentPrompt_Image
	//	*StaticContentPrompt_Table
	//	*StaticContentPrompt_Media
	//	*StaticContentPrompt_List
	//	*StaticContentPrompt_Collection
	//	*StaticContentPrompt_CollectionBrowse
	Content isStaticContentPrompt_Content `protobuf_oneof:"content"`
}

func (x *StaticContentPrompt) Reset() {
	*x = StaticContentPrompt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StaticContentPrompt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticContentPrompt) ProtoMessage() {}

func (x *StaticContentPrompt) ProtoReflect() protoreflect.Message {
	mi := &file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticContentPrompt.ProtoReflect.Descriptor instead.
func (*StaticContentPrompt) Descriptor() ([]byte, []int) {
	return file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescGZIP(), []int{0}
}

func (m *StaticContentPrompt) GetContent() isStaticContentPrompt_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *StaticContentPrompt) GetCard() *StaticCardPrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_Card); ok {
		return x.Card
	}
	return nil
}

func (x *StaticContentPrompt) GetImage() *StaticImagePrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_Image); ok {
		return x.Image
	}
	return nil
}

func (x *StaticContentPrompt) GetTable() *StaticTablePrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_Table); ok {
		return x.Table
	}
	return nil
}

func (x *StaticContentPrompt) GetMedia() *StaticMediaPrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_Media); ok {
		return x.Media
	}
	return nil
}

func (x *StaticContentPrompt) GetList() *StaticListPrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_List); ok {
		return x.List
	}
	return nil
}

func (x *StaticContentPrompt) GetCollection() *StaticCollectionPrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_Collection); ok {
		return x.Collection
	}
	return nil
}

func (x *StaticContentPrompt) GetCollectionBrowse() *StaticCollectionBrowsePrompt {
	if x, ok := x.GetContent().(*StaticContentPrompt_CollectionBrowse); ok {
		return x.CollectionBrowse
	}
	return nil
}

type isStaticContentPrompt_Content interface {
	isStaticContentPrompt_Content()
}

type StaticContentPrompt_Card struct {
	// A basic card.
	Card *StaticCardPrompt `protobuf:"bytes,1,opt,name=card,proto3,oneof"`
}

type StaticContentPrompt_Image struct {
	// An image.
	Image *StaticImagePrompt `protobuf:"bytes,2,opt,name=image,proto3,oneof"`
}

type StaticContentPrompt_Table struct {
	// Table card.
	Table *StaticTablePrompt `protobuf:"bytes,3,opt,name=table,proto3,oneof"`
}

type StaticContentPrompt_Media struct {
	// Response indicating a set of media to be played.
	Media *StaticMediaPrompt `protobuf:"bytes,4,opt,name=media,proto3,oneof"`
}

type StaticContentPrompt_List struct {
	// A card for presenting a list of options to select from.
	List *StaticListPrompt `protobuf:"bytes,5,opt,name=list,proto3,oneof"`
}

type StaticContentPrompt_Collection struct {
	// A card presenting a list of options to select from.
	Collection *StaticCollectionPrompt `protobuf:"bytes,6,opt,name=collection,proto3,oneof"`
}

type StaticContentPrompt_CollectionBrowse struct {
	// A card presenting a collection of web pages to open.
	CollectionBrowse *StaticCollectionBrowsePrompt `protobuf:"bytes,7,opt,name=collection_browse,json=collectionBrowse,proto3,oneof"`
}

func (*StaticContentPrompt_Card) isStaticContentPrompt_Content() {}

func (*StaticContentPrompt_Image) isStaticContentPrompt_Content() {}

func (*StaticContentPrompt_Table) isStaticContentPrompt_Content() {}

func (*StaticContentPrompt_Media) isStaticContentPrompt_Content() {}

func (*StaticContentPrompt_List) isStaticContentPrompt_Content() {}

func (*StaticContentPrompt_Collection) isStaticContentPrompt_Content() {}

func (*StaticContentPrompt_CollectionBrowse) isStaticContentPrompt_Content() {}

var File_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto protoreflect.FileDescriptor

var file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDesc = []byte{
	0x0a, 0x51, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x2d, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d,
	0x70, 0x74, 0x1a, 0x4e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x5f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x5b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x62, 0x72, 0x6f, 0x77,
	0x73, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x54, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2f, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x63, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x4f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72,
	0x6f, 0x6d, 0x70, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x4e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74,
	0x61, 0x74, 0x69, 0x63, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x4f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x74,
	0x61, 0x74, 0x69, 0x63, 0x5f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x4f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x6d,
	0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x05, 0x0a, 0x13, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x12, 0x55, 0x0a, 0x04, 0x63, 0x61, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x61, 0x72, 0x64, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x48,
	0x00, 0x52, 0x04, 0x63, 0x61, 0x72, 0x64, 0x12, 0x58, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x40, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x48, 0x00, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x12, 0x58, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x40, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x6d,
	0x70, 0x74, 0x48, 0x00, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x58, 0x0a, 0x05, 0x6d,
	0x65, 0x64, 0x69, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x40, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e,
	0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69,
	0x63, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x48, 0x00, 0x52, 0x05,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x55, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x3f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x6d, 0x70, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72,
	0x6f, 0x6d, 0x70, 0x74, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x67, 0x0a, 0x0a,
	0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x45, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x7a, 0x0a, 0x11, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x4b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x48, 0x00, 0x52,
	0x10, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x72, 0x6f, 0x77, 0x73,
	0x65, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0xa4, 0x01, 0x0a,
	0x31, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x76, 0x32, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x6d,
	0x70, 0x74, 0x42, 0x18, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x53,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72,
	0x67, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x64,
	0x6b, 0x2f, 0x76, 0x32, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x3b, 0x70, 0x72, 0x6f,
	0x6d, 0x70, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescOnce sync.Once
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescData = file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDesc
)

func file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescGZIP() []byte {
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescOnce.Do(func() {
		file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescData)
	})
	return file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDescData
}

var file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_goTypes = []interface{}{
	(*StaticContentPrompt)(nil),          // 0: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt
	(*StaticCardPrompt)(nil),             // 1: google.actions.sdk.v2.interactionmodel.prompt.StaticCardPrompt
	(*StaticImagePrompt)(nil),            // 2: google.actions.sdk.v2.interactionmodel.prompt.StaticImagePrompt
	(*StaticTablePrompt)(nil),            // 3: google.actions.sdk.v2.interactionmodel.prompt.StaticTablePrompt
	(*StaticMediaPrompt)(nil),            // 4: google.actions.sdk.v2.interactionmodel.prompt.StaticMediaPrompt
	(*StaticListPrompt)(nil),             // 5: google.actions.sdk.v2.interactionmodel.prompt.StaticListPrompt
	(*StaticCollectionPrompt)(nil),       // 6: google.actions.sdk.v2.interactionmodel.prompt.StaticCollectionPrompt
	(*StaticCollectionBrowsePrompt)(nil), // 7: google.actions.sdk.v2.interactionmodel.prompt.StaticCollectionBrowsePrompt
}
var file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_depIdxs = []int32{
	1, // 0: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.card:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticCardPrompt
	2, // 1: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.image:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticImagePrompt
	3, // 2: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.table:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticTablePrompt
	4, // 3: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.media:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticMediaPrompt
	5, // 4: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.list:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticListPrompt
	6, // 5: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.collection:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticCollectionPrompt
	7, // 6: google.actions.sdk.v2.interactionmodel.prompt.StaticContentPrompt.collection_browse:type_name -> google.actions.sdk.v2.interactionmodel.prompt.StaticCollectionBrowsePrompt
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() {
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_init()
}
func file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_init() {
	if File_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto != nil {
		return
	}
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_card_prompt_proto_init()
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_collection_browse_prompt_proto_init()
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_collection_prompt_proto_init()
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_image_prompt_proto_init()
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_list_prompt_proto_init()
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_media_prompt_proto_init()
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_table_prompt_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StaticContentPrompt); i {
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
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*StaticContentPrompt_Card)(nil),
		(*StaticContentPrompt_Image)(nil),
		(*StaticContentPrompt_Table)(nil),
		(*StaticContentPrompt_Media)(nil),
		(*StaticContentPrompt_List)(nil),
		(*StaticContentPrompt_Collection)(nil),
		(*StaticContentPrompt_CollectionBrowse)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_goTypes,
		DependencyIndexes: file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_depIdxs,
		MessageInfos:      file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_msgTypes,
	}.Build()
	File_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto = out.File
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_rawDesc = nil
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_goTypes = nil
	file_google_actions_sdk_v2_interactionmodel_prompt_content_static_content_prompt_proto_depIdxs = nil
}
