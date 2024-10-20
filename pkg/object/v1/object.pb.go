// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: object/v1/object.proto

package object

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Encryption_Algorithm int32

const (
	// No cryptography is being used
	Encryption_PLAINTEXT Encryption_Algorithm = 0
	// Encryption Algorithms
	Encryption_AES256_GCM Encryption_Algorithm = 110
	Encryption_AES192_GCM Encryption_Algorithm = 120
	Encryption_AES128_GCM Encryption_Algorithm = 130
	// Signature Algorithms
	Encryption_HMAC_SHA256 Encryption_Algorithm = 310
	// Sealing Algorithms (Asymmetric)
	Encryption_RSA_OAEP_SHA512 Encryption_Algorithm = 510
)

// Enum value maps for Encryption_Algorithm.
var (
	Encryption_Algorithm_name = map[int32]string{
		0:   "PLAINTEXT",
		110: "AES256_GCM",
		120: "AES192_GCM",
		130: "AES128_GCM",
		310: "HMAC_SHA256",
		510: "RSA_OAEP_SHA512",
	}
	Encryption_Algorithm_value = map[string]int32{
		"PLAINTEXT":       0,
		"AES256_GCM":      110,
		"AES192_GCM":      120,
		"AES128_GCM":      130,
		"HMAC_SHA256":     310,
		"RSA_OAEP_SHA512": 510,
	}
)

func (x Encryption_Algorithm) Enum() *Encryption_Algorithm {
	p := new(Encryption_Algorithm)
	*p = x
	return p
}

func (x Encryption_Algorithm) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Encryption_Algorithm) Descriptor() protoreflect.EnumDescriptor {
	return file_object_v1_object_proto_enumTypes[0].Descriptor()
}

func (Encryption_Algorithm) Type() protoreflect.EnumType {
	return &file_object_v1_object_proto_enumTypes[0]
}

func (x Encryption_Algorithm) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Encryption_Algorithm.Descriptor instead.
func (Encryption_Algorithm) EnumDescriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{4, 0}
}

type Compression_Algorithm int32

const (
	Compression_NONE     Compression_Algorithm = 0
	Compression_GZIP     Compression_Algorithm = 1
	Compression_COMPRESS Compression_Algorithm = 2
	Compression_DEFLATE  Compression_Algorithm = 3
	Compression_BROTLI   Compression_Algorithm = 4
)

// Enum value maps for Compression_Algorithm.
var (
	Compression_Algorithm_name = map[int32]string{
		0: "NONE",
		1: "GZIP",
		2: "COMPRESS",
		3: "DEFLATE",
		4: "BROTLI",
	}
	Compression_Algorithm_value = map[string]int32{
		"NONE":     0,
		"GZIP":     1,
		"COMPRESS": 2,
		"DEFLATE":  3,
		"BROTLI":   4,
	}
)

func (x Compression_Algorithm) Enum() *Compression_Algorithm {
	p := new(Compression_Algorithm)
	*p = x
	return p
}

func (x Compression_Algorithm) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Compression_Algorithm) Descriptor() protoreflect.EnumDescriptor {
	return file_object_v1_object_proto_enumTypes[1].Descriptor()
}

func (Compression_Algorithm) Type() protoreflect.EnumType {
	return &file_object_v1_object_proto_enumTypes[1]
}

func (x Compression_Algorithm) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Compression_Algorithm.Descriptor instead.
func (Compression_Algorithm) EnumDescriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{5, 0}
}

type Object struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version  *Version       `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Schema   *SchemaVersion `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
	Mimetype string         `protobuf:"bytes,3,opt,name=mimetype,proto3" json:"mimetype,omitempty"`
	// Access Controls
	Owner       []byte `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	Group       []byte `protobuf:"bytes,5,opt,name=group,proto3" json:"group,omitempty"`
	Permissions []byte `protobuf:"bytes,6,opt,name=permissions,proto3" json:"permissions,omitempty"`
	Acl         []*ACL `protobuf:"bytes,7,rep,name=acl,proto3" json:"acl,omitempty"`
	// Provenance Information
	WriteRegions []string   `protobuf:"bytes,8,rep,name=WriteRegions,proto3" json:"WriteRegions,omitempty"`
	Publisher    *Publisher `protobuf:"bytes,9,opt,name=publisher,proto3" json:"publisher,omitempty"`
	// Read Information
	Encryption  *Encryption  `protobuf:"bytes,10,opt,name=encryption,proto3" json:"encryption,omitempty"`
	Compression *Compression `protobuf:"bytes,11,opt,name=compression,proto3" json:"compression,omitempty"`
	// Flags
	Flags []byte `protobuf:"bytes,12,opt,name=flags,proto3" json:"flags,omitempty"`
	// Modification Timestamps
	Created  *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=created,proto3" json:"created,omitempty"`
	Modified *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=modified,proto3" json:"modified,omitempty"`
	// The actual data of the object
	Data []byte `protobuf:"bytes,15,opt,name=data,proto3" json:"data,omitempty"`
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

func (x *Object) GetVersion() *Version {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *Object) GetSchema() *SchemaVersion {
	if x != nil {
		return x.Schema
	}
	return nil
}

func (x *Object) GetMimetype() string {
	if x != nil {
		return x.Mimetype
	}
	return ""
}

func (x *Object) GetOwner() []byte {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *Object) GetGroup() []byte {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *Object) GetPermissions() []byte {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *Object) GetAcl() []*ACL {
	if x != nil {
		return x.Acl
	}
	return nil
}

func (x *Object) GetWriteRegions() []string {
	if x != nil {
		return x.WriteRegions
	}
	return nil
}

func (x *Object) GetPublisher() *Publisher {
	if x != nil {
		return x.Publisher
	}
	return nil
}

func (x *Object) GetEncryption() *Encryption {
	if x != nil {
		return x.Encryption
	}
	return nil
}

func (x *Object) GetCompression() *Compression {
	if x != nil {
		return x.Compression
	}
	return nil
}

func (x *Object) GetFlags() []byte {
	if x != nil {
		return x.Flags
	}
	return nil
}

func (x *Object) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Object) GetModified() *timestamppb.Timestamp {
	if x != nil {
		return x.Modified
	}
	return nil
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
	// The timestamp that the version was created (e.g. the last modified date).
	Created *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=created,proto3" json:"created,omitempty"`
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

func (x *Version) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

// An event type is composed of a name and a version so that the type can be looked up
// in the schema registry. The schema can then be used to validate the data inside the
// event. Schemas are optional but types are not unless the mimetype requries a schema
// for deserialization (e.g. protobuf, parquet, avro, etc.).
type SchemaVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	MajorVersion uint32 `protobuf:"varint,2,opt,name=major_version,json=majorVersion,proto3" json:"major_version,omitempty"`
	MinorVersion uint32 `protobuf:"varint,3,opt,name=minor_version,json=minorVersion,proto3" json:"minor_version,omitempty"`
	PatchVersion uint32 `protobuf:"varint,4,opt,name=patch_version,json=patchVersion,proto3" json:"patch_version,omitempty"`
}

func (x *SchemaVersion) Reset() {
	*x = SchemaVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaVersion) ProtoMessage() {}

func (x *SchemaVersion) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaVersion.ProtoReflect.Descriptor instead.
func (*SchemaVersion) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{2}
}

func (x *SchemaVersion) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SchemaVersion) GetMajorVersion() uint32 {
	if x != nil {
		return x.MajorVersion
	}
	return 0
}

func (x *SchemaVersion) GetMinorVersion() uint32 {
	if x != nil {
		return x.MinorVersion
	}
	return 0
}

func (x *SchemaVersion) GetPatchVersion() uint32 {
	if x != nil {
		return x.PatchVersion
	}
	return 0
}

type ACL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId    []byte `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Permissions []byte `protobuf:"bytes,2,opt,name=permissions,proto3" json:"permissions,omitempty"`
}

func (x *ACL) Reset() {
	*x = ACL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ACL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ACL) ProtoMessage() {}

func (x *ACL) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ACL.ProtoReflect.Descriptor instead.
func (*ACL) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{3}
}

func (x *ACL) GetClientId() []byte {
	if x != nil {
		return x.ClientId
	}
	return nil
}

func (x *ACL) GetPermissions() []byte {
	if x != nil {
		return x.Permissions
	}
	return nil
}

// Metadata about the cryptography used to secure the event.
type Encryption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKeyId         string               `protobuf:"bytes,1,opt,name=public_key_id,json=publicKeyId,proto3" json:"public_key_id,omitempty"`
	EncryptionKey       []byte               `protobuf:"bytes,2,opt,name=encryption_key,json=encryptionKey,proto3" json:"encryption_key,omitempty"`
	HmacSecret          []byte               `protobuf:"bytes,3,opt,name=hmac_secret,json=hmacSecret,proto3" json:"hmac_secret,omitempty"`
	Signature           []byte               `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	SealingAlgorithm    Encryption_Algorithm `protobuf:"varint,5,opt,name=sealing_algorithm,json=sealingAlgorithm,proto3,enum=object.v1.Encryption_Algorithm" json:"sealing_algorithm,omitempty"`
	EncryptionAlgorithm Encryption_Algorithm `protobuf:"varint,6,opt,name=encryption_algorithm,json=encryptionAlgorithm,proto3,enum=object.v1.Encryption_Algorithm" json:"encryption_algorithm,omitempty"`
	SignatureAlgorithm  Encryption_Algorithm `protobuf:"varint,7,opt,name=signature_algorithm,json=signatureAlgorithm,proto3,enum=object.v1.Encryption_Algorithm" json:"signature_algorithm,omitempty"`
}

func (x *Encryption) Reset() {
	*x = Encryption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Encryption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Encryption) ProtoMessage() {}

func (x *Encryption) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Encryption.ProtoReflect.Descriptor instead.
func (*Encryption) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{4}
}

func (x *Encryption) GetPublicKeyId() string {
	if x != nil {
		return x.PublicKeyId
	}
	return ""
}

func (x *Encryption) GetEncryptionKey() []byte {
	if x != nil {
		return x.EncryptionKey
	}
	return nil
}

func (x *Encryption) GetHmacSecret() []byte {
	if x != nil {
		return x.HmacSecret
	}
	return nil
}

func (x *Encryption) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Encryption) GetSealingAlgorithm() Encryption_Algorithm {
	if x != nil {
		return x.SealingAlgorithm
	}
	return Encryption_PLAINTEXT
}

func (x *Encryption) GetEncryptionAlgorithm() Encryption_Algorithm {
	if x != nil {
		return x.EncryptionAlgorithm
	}
	return Encryption_PLAINTEXT
}

func (x *Encryption) GetSignatureAlgorithm() Encryption_Algorithm {
	if x != nil {
		return x.SignatureAlgorithm
	}
	return Encryption_PLAINTEXT
}

// Metadata about compression used to reduce the storage size of the event.
type Compression struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algorithm Compression_Algorithm `protobuf:"varint,1,opt,name=algorithm,proto3,enum=object.v1.Compression_Algorithm" json:"algorithm,omitempty"`
	Level     int64                 `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"`
}

func (x *Compression) Reset() {
	*x = Compression{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Compression) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Compression) ProtoMessage() {}

func (x *Compression) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Compression.ProtoReflect.Descriptor instead.
func (*Compression) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{5}
}

func (x *Compression) GetAlgorithm() Compression_Algorithm {
	if x != nil {
		return x.Algorithm
	}
	return Compression_NONE
}

func (x *Compression) GetLevel() int64 {
	if x != nil {
		return x.Level
	}
	return 0
}

// Information about the publisher of the event for provenance and auditing purposes.
type Publisher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublisherId []byte `protobuf:"bytes,1,opt,name=publisher_id,json=publisherId,proto3" json:"publisher_id,omitempty"`
	ClientId    []byte `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	IpAddr      string `protobuf:"bytes,3,opt,name=ip_addr,json=ipAddr,proto3" json:"ip_addr,omitempty"`
	UserAgent   string `protobuf:"bytes,4,opt,name=user_agent,json=userAgent,proto3" json:"user_agent,omitempty"`
}

func (x *Publisher) Reset() {
	*x = Publisher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_object_v1_object_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publisher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publisher) ProtoMessage() {}

func (x *Publisher) ProtoReflect() protoreflect.Message {
	mi := &file_object_v1_object_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publisher.ProtoReflect.Descriptor instead.
func (*Publisher) Descriptor() ([]byte, []int) {
	return file_object_v1_object_proto_rawDescGZIP(), []int{6}
}

func (x *Publisher) GetPublisherId() []byte {
	if x != nil {
		return x.PublisherId
	}
	return nil
}

func (x *Publisher) GetClientId() []byte {
	if x != nil {
		return x.ClientId
	}
	return nil
}

func (x *Publisher) GetIpAddr() string {
	if x != nil {
		return x.IpAddr
	}
	return ""
}

func (x *Publisher) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}

var File_object_v1_object_proto protoreflect.FileDescriptor

var file_object_v1_object_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5, 0x04, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x2c, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x0a,
	0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12,
	0x1a, 0x0a, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x70, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x20, 0x0a, 0x03, 0x61, 0x63, 0x6c,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x43, 0x4c, 0x52, 0x03, 0x61, 0x63, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x32, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73,
	0x68, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x0a, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a,
	0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x0b, 0x63, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x12, 0x36, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08,
	0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xcd, 0x01, 0x0a,
	0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x06,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x6f, 0x6d, 0x62,
	0x73, 0x74, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x74, 0x6f, 0x6d,
	0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22, 0x92, 0x01, 0x0a,
	0x0d, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x6d, 0x61, 0x6a, 0x6f, 0x72,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x69, 0x6e, 0x6f, 0x72,
	0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c,
	0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0c, 0x70, 0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x22, 0x44, 0x0a, 0x03, 0x41, 0x43, 0x4c, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xff, 0x03, 0x0a, 0x0a, 0x45, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e,
	0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65,
	0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x6d, 0x61, 0x63, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x68, 0x6d, 0x61, 0x63, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x12, 0x4c, 0x0a, 0x11, 0x73, 0x65, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x10, 0x73, 0x65,
	0x61, 0x6c, 0x69, 0x6e, 0x67, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x52,
	0x0a, 0x14, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x13, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74,
	0x68, 0x6d, 0x12, 0x50, 0x0a, 0x13, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1f, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d,
	0x52, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72,
	0x69, 0x74, 0x68, 0x6d, 0x22, 0x73, 0x0a, 0x09, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x12, 0x0d, 0x0a, 0x09, 0x50, 0x4c, 0x41, 0x49, 0x4e, 0x54, 0x45, 0x58, 0x54, 0x10, 0x00,
	0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x45, 0x53, 0x32, 0x35, 0x36, 0x5f, 0x47, 0x43, 0x4d, 0x10, 0x6e,
	0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x45, 0x53, 0x31, 0x39, 0x32, 0x5f, 0x47, 0x43, 0x4d, 0x10, 0x78,
	0x12, 0x0f, 0x0a, 0x0a, 0x41, 0x45, 0x53, 0x31, 0x32, 0x38, 0x5f, 0x47, 0x43, 0x4d, 0x10, 0x82,
	0x01, 0x12, 0x10, 0x0a, 0x0b, 0x48, 0x4d, 0x41, 0x43, 0x5f, 0x53, 0x48, 0x41, 0x32, 0x35, 0x36,
	0x10, 0xb6, 0x02, 0x12, 0x14, 0x0a, 0x0f, 0x52, 0x53, 0x41, 0x5f, 0x4f, 0x41, 0x45, 0x50, 0x5f,
	0x53, 0x48, 0x41, 0x35, 0x31, 0x32, 0x10, 0xfe, 0x03, 0x22, 0xab, 0x01, 0x0a, 0x0b, 0x43, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3e, 0x0a, 0x09, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x09,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x22,
	0x46, 0x0a, 0x09, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x08, 0x0a, 0x04,
	0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x5a, 0x49, 0x50, 0x10, 0x01,
	0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4d, 0x50, 0x52, 0x45, 0x53, 0x53, 0x10, 0x02, 0x12, 0x0b,
	0x0a, 0x07, 0x44, 0x45, 0x46, 0x4c, 0x41, 0x54, 0x45, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x42,
	0x52, 0x4f, 0x54, 0x4c, 0x49, 0x10, 0x04, 0x22, 0x83, 0x01, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x73, 0x68, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x12, 0x1d,
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_object_v1_object_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_object_v1_object_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_object_v1_object_proto_goTypes = []any{
	(Encryption_Algorithm)(0),     // 0: object.v1.Encryption.Algorithm
	(Compression_Algorithm)(0),    // 1: object.v1.Compression.Algorithm
	(*Object)(nil),                // 2: object.v1.Object
	(*Version)(nil),               // 3: object.v1.Version
	(*SchemaVersion)(nil),         // 4: object.v1.SchemaVersion
	(*ACL)(nil),                   // 5: object.v1.ACL
	(*Encryption)(nil),            // 6: object.v1.Encryption
	(*Compression)(nil),           // 7: object.v1.Compression
	(*Publisher)(nil),             // 8: object.v1.Publisher
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_object_v1_object_proto_depIdxs = []int32{
	3,  // 0: object.v1.Object.version:type_name -> object.v1.Version
	4,  // 1: object.v1.Object.schema:type_name -> object.v1.SchemaVersion
	5,  // 2: object.v1.Object.acl:type_name -> object.v1.ACL
	8,  // 3: object.v1.Object.publisher:type_name -> object.v1.Publisher
	6,  // 4: object.v1.Object.encryption:type_name -> object.v1.Encryption
	7,  // 5: object.v1.Object.compression:type_name -> object.v1.Compression
	9,  // 6: object.v1.Object.created:type_name -> google.protobuf.Timestamp
	9,  // 7: object.v1.Object.modified:type_name -> google.protobuf.Timestamp
	3,  // 8: object.v1.Version.parent:type_name -> object.v1.Version
	9,  // 9: object.v1.Version.created:type_name -> google.protobuf.Timestamp
	0,  // 10: object.v1.Encryption.sealing_algorithm:type_name -> object.v1.Encryption.Algorithm
	0,  // 11: object.v1.Encryption.encryption_algorithm:type_name -> object.v1.Encryption.Algorithm
	0,  // 12: object.v1.Encryption.signature_algorithm:type_name -> object.v1.Encryption.Algorithm
	1,  // 13: object.v1.Compression.algorithm:type_name -> object.v1.Compression.Algorithm
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_object_v1_object_proto_init() }
func file_object_v1_object_proto_init() {
	if File_object_v1_object_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_object_v1_object_proto_msgTypes[0].Exporter = func(v any, i int) any {
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
		file_object_v1_object_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
		file_object_v1_object_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SchemaVersion); i {
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
		file_object_v1_object_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ACL); i {
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
		file_object_v1_object_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Encryption); i {
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
		file_object_v1_object_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Compression); i {
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
		file_object_v1_object_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*Publisher); i {
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
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_object_v1_object_proto_goTypes,
		DependencyIndexes: file_object_v1_object_proto_depIdxs,
		EnumInfos:         file_object_v1_object_proto_enumTypes,
		MessageInfos:      file_object_v1_object_proto_msgTypes,
	}.Build()
	File_object_v1_object_proto = out.File
	file_object_v1_object_proto_rawDesc = nil
	file_object_v1_object_proto_goTypes = nil
	file_object_v1_object_proto_depIdxs = nil
}
