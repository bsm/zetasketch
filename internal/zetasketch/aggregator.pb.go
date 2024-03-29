//
// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file contains messages for representing the internal state of
// an aggregation algorithm, common properties of all aggregation
// algorithms and common per-result element properties. Algorithms
// specific properties should be added as extensions to different
// proto files in the same directory.
//
// Adding a new algorithm requires the following steps:
//   1. Add a new value with a descriptive name to the AggregatorType enum.
//   2. Add an extension with the same tag as the enum value to
//      AggregatorStateProto to hold the serialized state of the new
//      algorithm.
//   3. [optional] Add an extension with the same tag as the enum
//      value to AggregatorValueStatsProto to hold meta data for each
//      element in the result set.
//   4. [optional] Add an extension with the same tag as the enum value to
//      AggregatorStatsProto to hold additional run-time statistics for
//      the aggregator.
//
// Each algorithm will have its own extension, rather than a single
// range for all extensions since it's easy to make a mistake.
//
// Messages defined in this file may be stored on disk, so the
// aggregation library should be able to parse all historic versions
// of the serialized data and it should be able to merge data with
// different serialization formats.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.0
// source: aggregator.proto

package zetasketch

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Enumeration of all supported aggregation algorithms. Values should
// start from 100.
type AggregatorType int32

const (
	// Sum all values added to the aggregator.
	AggregatorType_SUM AggregatorType = 100
	// Computes a cardinality estimation using the HyperLogLog++ algorithm.
	AggregatorType_HYPERLOGLOG_PLUS_UNIQUE AggregatorType = 112
)

// Enum value maps for AggregatorType.
var (
	AggregatorType_name = map[int32]string{
		100: "SUM",
		112: "HYPERLOGLOG_PLUS_UNIQUE",
	}
	AggregatorType_value = map[string]int32{
		"SUM":                     100,
		"HYPERLOGLOG_PLUS_UNIQUE": 112,
	}
)

func (x AggregatorType) Enum() *AggregatorType {
	p := new(AggregatorType)
	*p = x
	return p
}

func (x AggregatorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AggregatorType) Descriptor() protoreflect.EnumDescriptor {
	return file_aggregator_proto_enumTypes[0].Descriptor()
}

func (AggregatorType) Type() protoreflect.EnumType {
	return &file_aggregator_proto_enumTypes[0]
}

func (x AggregatorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *AggregatorType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = AggregatorType(num)
	return nil
}

// Deprecated: Use AggregatorType.Descriptor instead.
func (AggregatorType) EnumDescriptor() ([]byte, []int) {
	return file_aggregator_proto_rawDescGZIP(), []int{0}
}

// Each value corresponds to a C++ type T and its corresponding
// DefaultOps<T> instantiation. A ValueOps implementation returning
// something other than UNKNOWN for a given value is promising that the value
// of the type corresponding to the value, and that the Ops implementation
// performs identical operations as DefaultOps<T> for that type.
type DefaultOpsType_Id int32

const (
	DefaultOpsType_UNKNOWN DefaultOpsType_Id = 0
	// int8, DefaultOps<int8>
	// SerializeToString writes the single 2s-complement byte.
	DefaultOpsType_INT8 DefaultOpsType_Id = 1
	// int16, DefaultOps<int16>
	// SerializeToString writes the two little-endian 2s-complement bytes.
	DefaultOpsType_INT16 DefaultOpsType_Id = 2
	// int32, DefaultOps<int32>
	// SerializeToString uses varint encoding of the 2s complement in 32 bits -
	// i.e. the result for negative integers is 5 bytes long, not 10.
	DefaultOpsType_INT32 DefaultOpsType_Id = 3
	// int64, DefaultOps<int64>
	// SerializeToString uses varint encoding of the 2s complement.
	DefaultOpsType_INT64 DefaultOpsType_Id = 4
	// uint8, DefaultOps<uint8>
	// SerializeToString writes the single byte.
	DefaultOpsType_UINT8 DefaultOpsType_Id = 5
	// uint16, DefaultOps<uint16>
	// SerializeToString writes the two little-endian bytes.
	DefaultOpsType_UINT16 DefaultOpsType_Id = 6
	// uint32, DefaultOps<uint32>
	// SerializeToString uses varint encoding.
	DefaultOpsType_UINT32 DefaultOpsType_Id = 7
	// uint64, DefaultOps<uint64>
	// SerializeToString uses varint encoding.
	DefaultOpsType_UINT64 DefaultOpsType_Id = 8
	// float, DefaultOps<float>
	// SerializeToString encodes the 4 little endian IEEE754 bytes.
	DefaultOpsType_FLOAT DefaultOpsType_Id = 9
	// double, DefaultOps<double>
	// SerializeToString encodes the 8 little endian IEEE754 bytes.
	DefaultOpsType_DOUBLE DefaultOpsType_Id = 10
	// string, DefaultOps<string>
	// SerializeToString just copies the bytes.
	DefaultOpsType_BYTES_OR_UTF8_STRING DefaultOpsType_Id = 11
)

// Enum value maps for DefaultOpsType_Id.
var (
	DefaultOpsType_Id_name = map[int32]string{
		0:  "UNKNOWN",
		1:  "INT8",
		2:  "INT16",
		3:  "INT32",
		4:  "INT64",
		5:  "UINT8",
		6:  "UINT16",
		7:  "UINT32",
		8:  "UINT64",
		9:  "FLOAT",
		10: "DOUBLE",
		11: "BYTES_OR_UTF8_STRING",
	}
	DefaultOpsType_Id_value = map[string]int32{
		"UNKNOWN":              0,
		"INT8":                 1,
		"INT16":                2,
		"INT32":                3,
		"INT64":                4,
		"UINT8":                5,
		"UINT16":               6,
		"UINT32":               7,
		"UINT64":               8,
		"FLOAT":                9,
		"DOUBLE":               10,
		"BYTES_OR_UTF8_STRING": 11,
	}
)

func (x DefaultOpsType_Id) Enum() *DefaultOpsType_Id {
	p := new(DefaultOpsType_Id)
	*p = x
	return p
}

func (x DefaultOpsType_Id) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DefaultOpsType_Id) Descriptor() protoreflect.EnumDescriptor {
	return file_aggregator_proto_enumTypes[1].Descriptor()
}

func (DefaultOpsType_Id) Type() protoreflect.EnumType {
	return &file_aggregator_proto_enumTypes[1]
}

func (x DefaultOpsType_Id) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *DefaultOpsType_Id) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = DefaultOpsType_Id(num)
	return nil
}

// Deprecated: Use DefaultOpsType_Id.Descriptor instead.
func (DefaultOpsType_Id) EnumDescriptor() ([]byte, []int) {
	return file_aggregator_proto_rawDescGZIP(), []int{0, 0}
}

// Never instantiated, just for scoping an enum and associated options.
type DefaultOpsType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DefaultOpsType) Reset() {
	*x = DefaultOpsType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DefaultOpsType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DefaultOpsType) ProtoMessage() {}

func (x *DefaultOpsType) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DefaultOpsType.ProtoReflect.Descriptor instead.
func (*DefaultOpsType) Descriptor() ([]byte, []int) {
	return file_aggregator_proto_rawDescGZIP(), []int{0}
}

// This message contains common "public" properties of an aggregation
// algorithm. Add additional fields here only if they make sense for
// all algorithms.
type AggregatorStatsProto struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// Total number of values added to this aggregator.
	NumValues *int64 `protobuf:"varint,1,req,name=num_values,json=numValues" json:"num_values,omitempty"`
}

func (x *AggregatorStatsProto) Reset() {
	*x = AggregatorStatsProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatorStatsProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatorStatsProto) ProtoMessage() {}

func (x *AggregatorStatsProto) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatorStatsProto.ProtoReflect.Descriptor instead.
func (*AggregatorStatsProto) Descriptor() ([]byte, []int) {
	return file_aggregator_proto_rawDescGZIP(), []int{1}
}

func (x *AggregatorStatsProto) GetNumValues() int64 {
	if x != nil && x.NumValues != nil {
		return *x.NumValues
	}
	return 0
}

// Serialized state of an aggregator. Add additional fields here only
// if they make sense for all algorithms and if it doesn't make sense to
// expose them to the users of the library, e.g. encoding version.
type AggregatorStateProto struct {
	state           protoimpl.MessageState
	sizeCache       protoimpl.SizeCache
	unknownFields   protoimpl.UnknownFields
	extensionFields protoimpl.ExtensionFields

	// The type of the aggregator.
	Type      *AggregatorType `protobuf:"varint,1,req,name=type,enum=zetasketch.AggregatorType" json:"type,omitempty"`
	NumValues *int64          `protobuf:"varint,2,req,name=num_values,json=numValues" json:"num_values,omitempty"`
	// Version of the encoded internal state. On a per-aggregator basis, set this
	// field to indicate that the format of the aggregator encoding has changed
	// such that the library has to decide how to decode. Do NOT change the
	// default value, as this affects all aggregators.
	EncodingVersion *int32 `protobuf:"varint,3,opt,name=encoding_version,json=encodingVersion,def=1" json:"encoding_version,omitempty"`
	// Specifies the value type for the aggregation.
	//
	// If the value type is one supported by the DefaultOps<T> template, and that
	// set of operations (or a compatible implementation) was used, then this will
	// be a value of the DefaultOpsType.Id enum.
	//
	// Otherwise, this is a globally unique number corresponding to the value and
	// Ops implementation (e.g. the CL number in which the implementation is
	// defined). Values for custom types should be greater than 1000. Implementors
	// should consider registering a name for their custom type in
	// custom-value-type.proto, to facilitate easier discovery and better error
	// messages when conflicting types are merged.
	ValueType *int32 `protobuf:"varint,4,opt,name=value_type,json=valueType" json:"value_type,omitempty"`
}

// Default values for AggregatorStateProto fields.
const (
	Default_AggregatorStateProto_EncodingVersion = int32(1)
)

func (x *AggregatorStateProto) Reset() {
	*x = AggregatorStateProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatorStateProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatorStateProto) ProtoMessage() {}

func (x *AggregatorStateProto) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatorStateProto.ProtoReflect.Descriptor instead.
func (*AggregatorStateProto) Descriptor() ([]byte, []int) {
	return file_aggregator_proto_rawDescGZIP(), []int{2}
}

func (x *AggregatorStateProto) GetType() AggregatorType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return AggregatorType_SUM
}

func (x *AggregatorStateProto) GetNumValues() int64 {
	if x != nil && x.NumValues != nil {
		return *x.NumValues
	}
	return 0
}

func (x *AggregatorStateProto) GetEncodingVersion() int32 {
	if x != nil && x.EncodingVersion != nil {
		return *x.EncodingVersion
	}
	return Default_AggregatorStateProto_EncodingVersion
}

func (x *AggregatorStateProto) GetValueType() int32 {
	if x != nil && x.ValueType != nil {
		return *x.ValueType
	}
	return 0
}

// Additional metadata for each element in the result iterator.
type AggregatorValueStatsProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AggregatorValueStatsProto) Reset() {
	*x = AggregatorValueStatsProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aggregator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatorValueStatsProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatorValueStatsProto) ProtoMessage() {}

func (x *AggregatorValueStatsProto) ProtoReflect() protoreflect.Message {
	mi := &file_aggregator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatorValueStatsProto.ProtoReflect.Descriptor instead.
func (*AggregatorValueStatsProto) Descriptor() ([]byte, []int) {
	return file_aggregator_proto_rawDescGZIP(), []int{3}
}

var file_aggregator_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*DefaultOpsType_Id)(nil),
		Field:         132643189,
		Name:          "zetasketch.DefaultOpsType.unsigned_counterpart",
		Tag:           "varint,132643189,opt,name=unsigned_counterpart,enum=zetasketch.DefaultOpsType_Id",
		Filename:      "aggregator.proto",
	},
}

// Extension fields to descriptorpb.EnumValueOptions.
var (
	// Meant to be used on Id values, which represent types. Specifies the
	// unsigned counterpart to the type.
	//
	// optional zetasketch.DefaultOpsType.Id unsigned_counterpart = 132643189;
	E_DefaultOpsType_UnsignedCounterpart = &file_aggregator_proto_extTypes[0]
)

var File_aggregator_proto protoreflect.FileDescriptor

var file_aggregator_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x7a, 0x65, 0x74, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x63, 0x68, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xdb, 0x02, 0x0a, 0x0e, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4f, 0x70, 0x73, 0x54,
	0x79, 0x70, 0x65, 0x22, 0xd0, 0x01, 0x0a, 0x02, 0x49, 0x64, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x04, 0x49, 0x4e, 0x54, 0x38, 0x10,
	0x01, 0x1a, 0x06, 0xa8, 0x97, 0xff, 0xf9, 0x03, 0x05, 0x12, 0x11, 0x0a, 0x05, 0x49, 0x4e, 0x54,
	0x31, 0x36, 0x10, 0x02, 0x1a, 0x06, 0xa8, 0x97, 0xff, 0xf9, 0x03, 0x06, 0x12, 0x11, 0x0a, 0x05,
	0x49, 0x4e, 0x54, 0x33, 0x32, 0x10, 0x03, 0x1a, 0x06, 0xa8, 0x97, 0xff, 0xf9, 0x03, 0x07, 0x12,
	0x11, 0x0a, 0x05, 0x49, 0x4e, 0x54, 0x36, 0x34, 0x10, 0x04, 0x1a, 0x06, 0xa8, 0x97, 0xff, 0xf9,
	0x03, 0x08, 0x12, 0x09, 0x0a, 0x05, 0x55, 0x49, 0x4e, 0x54, 0x38, 0x10, 0x05, 0x12, 0x0a, 0x0a,
	0x06, 0x55, 0x49, 0x4e, 0x54, 0x31, 0x36, 0x10, 0x06, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x49, 0x4e,
	0x54, 0x33, 0x32, 0x10, 0x07, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x49, 0x4e, 0x54, 0x36, 0x34, 0x10,
	0x08, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4c, 0x4f, 0x41, 0x54, 0x10, 0x09, 0x12, 0x0a, 0x0a, 0x06,
	0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10, 0x0a, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x59, 0x54, 0x45,
	0x53, 0x5f, 0x4f, 0x52, 0x5f, 0x55, 0x54, 0x46, 0x38, 0x5f, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47,
	0x10, 0x0b, 0x22, 0x04, 0x08, 0x0c, 0x10, 0x0c, 0x2a, 0x0c, 0x55, 0x54, 0x46, 0x31, 0x36, 0x5f,
	0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x32, 0x76, 0x0a, 0x14, 0x75, 0x6e, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x61, 0x72, 0x74, 0x12, 0x21,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xf5, 0xf2, 0x9f, 0x3f, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x7a, 0x65, 0x74,
	0x61, 0x73, 0x6b, 0x65, 0x74, 0x63, 0x68, 0x2e, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4f,
	0x70, 0x73, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x49, 0x64, 0x52, 0x13, 0x75, 0x6e, 0x73, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x61, 0x72, 0x74, 0x22, 0x48,
	0x0a, 0x14, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x02, 0x28, 0x03, 0x52, 0x09, 0x6e, 0x75, 0x6d, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x2a, 0x04, 0x08, 0x6c, 0x10, 0x70, 0x2a, 0x05, 0x08, 0x71, 0x10,
	0x8d, 0x01, 0x2a, 0x04, 0x08, 0x70, 0x10, 0x71, 0x22, 0xc5, 0x01, 0x0a, 0x14, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0e, 0x32,
	0x1a, 0x2e, 0x7a, 0x65, 0x74, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x63, 0x68, 0x2e, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x02, 0x28, 0x03, 0x52, 0x09, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x12, 0x2c, 0x0a, 0x10, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x3a, 0x01, 0x31, 0x52, 0x0f, 0x65,
	0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d,
	0x0a, 0x0a, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65, 0x2a, 0x04, 0x08,
	0x64, 0x10, 0x70, 0x2a, 0x05, 0x08, 0x71, 0x10, 0x8d, 0x01, 0x2a, 0x04, 0x08, 0x70, 0x10, 0x71,
	0x22, 0x1b, 0x0a, 0x19, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x49, 0x0a,
	0x0e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x07, 0x0a, 0x03, 0x53, 0x55, 0x4d, 0x10, 0x64, 0x12, 0x1b, 0x0a, 0x17, 0x48, 0x59, 0x50, 0x45,
	0x52, 0x4c, 0x4f, 0x47, 0x4c, 0x4f, 0x47, 0x5f, 0x50, 0x4c, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x49,
	0x51, 0x55, 0x45, 0x10, 0x70, 0x22, 0x04, 0x08, 0x00, 0x10, 0x00, 0x22, 0x04, 0x08, 0x65, 0x10,
	0x6f, 0x22, 0x05, 0x08, 0x71, 0x10, 0x8c, 0x01, 0x42, 0x21, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x7a, 0x65,
	0x74, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x63, 0x68, 0xf8, 0x01, 0x01,
}

var (
	file_aggregator_proto_rawDescOnce sync.Once
	file_aggregator_proto_rawDescData = file_aggregator_proto_rawDesc
)

func file_aggregator_proto_rawDescGZIP() []byte {
	file_aggregator_proto_rawDescOnce.Do(func() {
		file_aggregator_proto_rawDescData = protoimpl.X.CompressGZIP(file_aggregator_proto_rawDescData)
	})
	return file_aggregator_proto_rawDescData
}

var file_aggregator_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_aggregator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_aggregator_proto_goTypes = []interface{}{
	(AggregatorType)(0),                   // 0: zetasketch.AggregatorType
	(DefaultOpsType_Id)(0),                // 1: zetasketch.DefaultOpsType.Id
	(*DefaultOpsType)(nil),                // 2: zetasketch.DefaultOpsType
	(*AggregatorStatsProto)(nil),          // 3: zetasketch.AggregatorStatsProto
	(*AggregatorStateProto)(nil),          // 4: zetasketch.AggregatorStateProto
	(*AggregatorValueStatsProto)(nil),     // 5: zetasketch.AggregatorValueStatsProto
	(*descriptorpb.EnumValueOptions)(nil), // 6: google.protobuf.EnumValueOptions
}
var file_aggregator_proto_depIdxs = []int32{
	0, // 0: zetasketch.AggregatorStateProto.type:type_name -> zetasketch.AggregatorType
	6, // 1: zetasketch.DefaultOpsType.unsigned_counterpart:extendee -> google.protobuf.EnumValueOptions
	1, // 2: zetasketch.DefaultOpsType.unsigned_counterpart:type_name -> zetasketch.DefaultOpsType.Id
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	2, // [2:3] is the sub-list for extension type_name
	1, // [1:2] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_aggregator_proto_init() }
func file_aggregator_proto_init() {
	if File_aggregator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aggregator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DefaultOpsType); i {
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
		file_aggregator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatorStatsProto); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
		file_aggregator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatorStateProto); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			case 3:
				return &v.extensionFields
			default:
				return nil
			}
		}
		file_aggregator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatorValueStatsProto); i {
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
			RawDescriptor: file_aggregator_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_aggregator_proto_goTypes,
		DependencyIndexes: file_aggregator_proto_depIdxs,
		EnumInfos:         file_aggregator_proto_enumTypes,
		MessageInfos:      file_aggregator_proto_msgTypes,
		ExtensionInfos:    file_aggregator_proto_extTypes,
	}.Build()
	File_aggregator_proto = out.File
	file_aggregator_proto_rawDesc = nil
	file_aggregator_proto_goTypes = nil
	file_aggregator_proto_depIdxs = nil
}
