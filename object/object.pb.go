// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: object/v1/object.proto

package object

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// An Object is a generic representation of a replicated piece of data. At the top level
// it contains enough metadata in order to detect changes during anti-entropy,
// specifically, given two objects, which is the later object (or does it need to be
// created or deleted). When used in VersionVectors only the metadata of the object is
// supplied. When passed via Updates, then the full data of the object is populated.
//
// For research purposes, this anti-entropy mechanism tracks the region and owner of
// each object to determine provinence and global interactions. Because this is side
// channel information that may not necessarily be stored with the object data, it is
// recommended that an objects table be kept for fast lookups of object versions.
// Additionally, it is recommended that a history table be maintained locally so that
// Object versions can be rolled back to previous states where necessary.
type Object struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The object metadata that must be populated on both VersionVectors and Updates
	Key       []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`             // A unique key/id that represents the object across the namespace of the object type
	Namespace string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"` // The namespace of the object, used to disambiguate keys or different object types
	Version   *Version `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`     // A version vector representing this objects current or latest version
	Region    string   `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`       // The region code where the data originated
	Owner     string   `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`         // The replica that created the object (identified as "pid:name" if name exists)
	// The object data that is only populated on Updates.
	Data []byte `protobuf:"bytes,10,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Object) Reset() {
	*x = Object{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Object) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Object) ProtoMessage() {}

func (x *Object) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Object.ProtoReflect.Descriptor instead.
func (*Object) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{0}
}

func (x *Object) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Object) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Object) GetVersion() *Version {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *Object) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Object) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Object) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// Implements a geo-distributed version as a Lamport Scalar
type Version struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pid       uint64   `protobuf:"varint,1,opt,name=pid,proto3" json:"pid,omitempty"`             // Process ID - used to deconflict ties in the version number.
	Version   uint64   `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`     // Montonically increasing version number.
	Region    string   `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`        // The region where the change occurred to track multi-region handling.
	Parent    *Version `protobuf:"bytes,4,opt,name=parent,proto3" json:"parent,omitempty"`        // In order to get a complete version history, identify the predessor; for compact data transfer parent should not be defined in parent version.
	Tombstone bool     `protobuf:"varint,5,opt,name=tombstone,proto3" json:"tombstone,omitempty"` // Set to true in order to mark the object as deleted
}

func (x *Version) Reset() {
	*x = Version{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Version) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Version) ProtoMessage() {}

func (x *Version) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Version.ProtoReflect.Descriptor instead.
func (*Version) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{1}
}

func (x *Version) GetPid() uint64 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *Version) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Version) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Version) GetParent() *Version {
	if x != nil {
		return x.Parent
	}
	return nil
}

func (x *Version) GetTombstone() bool {
	if x != nil {
		return x.Tombstone
	}
	return false
}

var File_object_v1_object_proto protoreflect.FileDescriptor

var file_object_v1_object_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x68, 0x6f, 0x6e, 0x75, 0x2e, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x22, 0xad, 0x01, 0x0a, 0x06, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x68, 0x6f, 0x6e, 0x75, 0x2e, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x9c, 0x01, 0x0a, 0x07, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x68, 0x6f, 0x6e, 0x75, 0x2e,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x6f, 0x6d,
	0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x74, 0x6f,
	0x6d, 0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c,
	0x69, 0x6f, 0x2f, 0x68, 0x6f, 0x6e, 0x75, 0x2f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_object_v1_object_proto_rawDescOnce sync.Once
	file_object_v1_object_proto_rawDescData = file_object_v1_object_proto_rawDesc
)

func file_object_v1_object_proto_rawDescGZIP() []byte {
	file_object_v1_object_proto_rawDescOnce.Do(func() {
		file_object_v1_object_proto_rawDescData = protoimpl.X.CompressGZIP(file_object_v1_object_proto_rawDescData)
	})
	return file_object_v1_object_proto_rawDescData
}

var file_object_v1_object_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_object_v1_object_proto_goTypes = []interface{}{
	(*Object)(nil),  // 0: honu.object.v1.Object
	(*Version)(nil), // 1: honu.object.v1.Version
}
var file_object_v1_object_proto_depIdxs = []int32{
	1, // 0: honu.object.v1.Object.version:type_name -> honu.object.v1.Version
	1, // 1: honu.object.v1.Version.parent:type_name -> honu.object.v1.Version
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_object_v1_object_proto_init() }
func file_object_v1_object_proto_init() {
	if File_object_v1_object_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_object_v1_object_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Object); i {
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
		file_object_v1_object_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Version); i {
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
			RawDescriptor: file_object_v1_object_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_object_v1_object_proto_goTypes,
		DependencyIndexes: file_object_v1_object_proto_depIdxs,
		MessageInfos:      file_object_v1_object_proto_msgTypes,
	}.Build()
	File_object_v1_object_proto = out.File
	file_object_v1_object_proto_rawDesc = nil
	file_object_v1_object_proto_goTypes = nil
	file_object_v1_object_proto_depIdxs = nil
}
