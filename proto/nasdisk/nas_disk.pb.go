// 指定的当前proto语法的版本，有2和3

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: nas_disk.proto

// 指定文件生成出来的package

package nasdisk

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type CreateNasDiskRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32                     `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg     string                    `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	NasDisk *CreateNasDiskRes_NasDisk `protobuf:"bytes,3,opt,name=nas_disk,json=nasDisk,proto3" json:"nas_disk,omitempty"`
}

func (x *CreateNasDiskRes) Reset() {
	*x = CreateNasDiskRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNasDiskRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNasDiskRes) ProtoMessage() {}

func (x *CreateNasDiskRes) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNasDiskRes.ProtoReflect.Descriptor instead.
func (*CreateNasDiskRes) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{0}
}

func (x *CreateNasDiskRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *CreateNasDiskRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *CreateNasDiskRes) GetNasDisk() *CreateNasDiskRes_NasDisk {
	if x != nil {
		return x.NasDisk
	}
	return nil
}

type CreateNasDiskReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Apikey         string `protobuf:"bytes,1,opt,name=apikey,proto3" json:"apikey,omitempty"`
	TenantId       string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	PlatformUserid string `protobuf:"bytes,3,opt,name=platform_userid,json=platformUserid,proto3" json:"platform_userid,omitempty"`
	ShareName      string `protobuf:"bytes,4,opt,name=share_name,json=shareName,proto3" json:"share_name,omitempty"`
	ShareDesc      string `protobuf:"bytes,5,opt,name=share_desc,json=shareDesc,proto3" json:"share_desc,omitempty"`
	ShareProto     string `protobuf:"bytes,6,opt,name=share_proto,json=shareProto,proto3" json:"share_proto,omitempty"`
	ShareSizeInG   int32  `protobuf:"varint,7,opt,name=share_size_in_g,json=shareSizeInG,proto3" json:"share_size_in_g,omitempty"`
	Region         string `protobuf:"bytes,8,opt,name=region,proto3" json:"region,omitempty"`
	NetworkId      string `protobuf:"bytes,9,opt,name=network_id,json=networkId,proto3" json:"network_id,omitempty"`
}

func (x *CreateNasDiskReq) Reset() {
	*x = CreateNasDiskReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNasDiskReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNasDiskReq) ProtoMessage() {}

func (x *CreateNasDiskReq) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNasDiskReq.ProtoReflect.Descriptor instead.
func (*CreateNasDiskReq) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{1}
}

func (x *CreateNasDiskReq) GetApikey() string {
	if x != nil {
		return x.Apikey
	}
	return ""
}

func (x *CreateNasDiskReq) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *CreateNasDiskReq) GetPlatformUserid() string {
	if x != nil {
		return x.PlatformUserid
	}
	return ""
}

func (x *CreateNasDiskReq) GetShareName() string {
	if x != nil {
		return x.ShareName
	}
	return ""
}

func (x *CreateNasDiskReq) GetShareDesc() string {
	if x != nil {
		return x.ShareDesc
	}
	return ""
}

func (x *CreateNasDiskReq) GetShareProto() string {
	if x != nil {
		return x.ShareProto
	}
	return ""
}

func (x *CreateNasDiskReq) GetShareSizeInG() int32 {
	if x != nil {
		return x.ShareSizeInG
	}
	return 0
}

func (x *CreateNasDiskReq) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *CreateNasDiskReq) GetNetworkId() string {
	if x != nil {
		return x.NetworkId
	}
	return ""
}

type DeleteNasDiskReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Apikey         string `protobuf:"bytes,1,opt,name=apikey,proto3" json:"apikey,omitempty"`
	TenantId       string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	PlatformUserid string `protobuf:"bytes,3,opt,name=platform_userid,json=platformUserid,proto3" json:"platform_userid,omitempty"`
	ShareId        string `protobuf:"bytes,4,opt,name=share_id,json=shareId,proto3" json:"share_id,omitempty"`
}

func (x *DeleteNasDiskReq) Reset() {
	*x = DeleteNasDiskReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNasDiskReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNasDiskReq) ProtoMessage() {}

func (x *DeleteNasDiskReq) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNasDiskReq.ProtoReflect.Descriptor instead.
func (*DeleteNasDiskReq) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteNasDiskReq) GetApikey() string {
	if x != nil {
		return x.Apikey
	}
	return ""
}

func (x *DeleteNasDiskReq) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *DeleteNasDiskReq) GetPlatformUserid() string {
	if x != nil {
		return x.PlatformUserid
	}
	return ""
}

func (x *DeleteNasDiskReq) GetShareId() string {
	if x != nil {
		return x.ShareId
	}
	return ""
}

type DeleteNasDiskRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code        int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg         string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	ShareId     string `protobuf:"bytes,3,opt,name=share_id,json=shareId,proto3" json:"share_id,omitempty"`
	DeletedTime string `protobuf:"bytes,4,opt,name=deleted_time,json=deletedTime,proto3" json:"deleted_time,omitempty"`
}

func (x *DeleteNasDiskRes) Reset() {
	*x = DeleteNasDiskRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNasDiskRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNasDiskRes) ProtoMessage() {}

func (x *DeleteNasDiskRes) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNasDiskRes.ProtoReflect.Descriptor instead.
func (*DeleteNasDiskRes) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteNasDiskRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *DeleteNasDiskRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *DeleteNasDiskRes) GetShareId() string {
	if x != nil {
		return x.ShareId
	}
	return ""
}

func (x *DeleteNasDiskRes) GetDeletedTime() string {
	if x != nil {
		return x.DeletedTime
	}
	return ""
}

type GetMountClientsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Apikey         string `protobuf:"bytes,1,opt,name=apikey,proto3" json:"apikey,omitempty"`
	TenantId       string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	PlatformUserid string `protobuf:"bytes,3,opt,name=platform_userid,json=platformUserid,proto3" json:"platform_userid,omitempty"`
	ShareId        string `protobuf:"bytes,4,opt,name=share_id,json=shareId,proto3" json:"share_id,omitempty"`
}

func (x *GetMountClientsReq) Reset() {
	*x = GetMountClientsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMountClientsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMountClientsReq) ProtoMessage() {}

func (x *GetMountClientsReq) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMountClientsReq.ProtoReflect.Descriptor instead.
func (*GetMountClientsReq) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{4}
}

func (x *GetMountClientsReq) GetApikey() string {
	if x != nil {
		return x.Apikey
	}
	return ""
}

func (x *GetMountClientsReq) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *GetMountClientsReq) GetPlatformUserid() string {
	if x != nil {
		return x.PlatformUserid
	}
	return ""
}

func (x *GetMountClientsReq) GetShareId() string {
	if x != nil {
		return x.ShareId
	}
	return ""
}

type GetMountClientsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code       int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg        string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	InstanceId []string `protobuf:"bytes,3,rep,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
}

func (x *GetMountClientsRes) Reset() {
	*x = GetMountClientsRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMountClientsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMountClientsRes) ProtoMessage() {}

func (x *GetMountClientsRes) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMountClientsRes.ProtoReflect.Descriptor instead.
func (*GetMountClientsRes) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{5}
}

func (x *GetMountClientsRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GetMountClientsRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GetMountClientsRes) GetInstanceId() []string {
	if x != nil {
		return x.InstanceId
	}
	return nil
}

type CreateNasDiskRes_NasDisk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShareId         string `protobuf:"bytes,1,opt,name=share_id,json=shareId,proto3" json:"share_id,omitempty"`
	ShareName       string `protobuf:"bytes,2,opt,name=share_name,json=shareName,proto3" json:"share_name,omitempty"`
	ShareDesc       string `protobuf:"bytes,3,opt,name=share_desc,json=shareDesc,proto3" json:"share_desc,omitempty"`
	ShareProto      string `protobuf:"bytes,4,opt,name=share_proto,json=shareProto,proto3" json:"share_proto,omitempty"`
	ShareSizeInG    int32  `protobuf:"varint,5,opt,name=share_size_in_g,json=shareSizeInG,proto3" json:"share_size_in_g,omitempty"`
	Region          string `protobuf:"bytes,6,opt,name=region,proto3" json:"region,omitempty"`
	NetworkId       string `protobuf:"bytes,7,opt,name=network_id,json=networkId,proto3" json:"network_id,omitempty"`
	ShareStatus     string `protobuf:"bytes,8,opt,name=share_status,json=shareStatus,proto3" json:"share_status,omitempty"`
	ShareProgress   string `protobuf:"bytes,9,opt,name=share_progress,json=shareProgress,proto3" json:"share_progress,omitempty"`
	ShareServerId   string `protobuf:"bytes,10,opt,name=share_server_id,json=shareServerId,proto3" json:"share_server_id,omitempty"`
	ShareServerHost string `protobuf:"bytes,11,opt,name=share_server_host,json=shareServerHost,proto3" json:"share_server_host,omitempty"`
	ShareNetworkId  string `protobuf:"bytes,12,opt,name=share_network_id,json=shareNetworkId,proto3" json:"share_network_id,omitempty"`
	CreatedTime     string `protobuf:"bytes,13,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	UpdatedTime     string `protobuf:"bytes,14,opt,name=updated_time,json=updatedTime,proto3" json:"updated_time,omitempty"`
}

func (x *CreateNasDiskRes_NasDisk) Reset() {
	*x = CreateNasDiskRes_NasDisk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nas_disk_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNasDiskRes_NasDisk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNasDiskRes_NasDisk) ProtoMessage() {}

func (x *CreateNasDiskRes_NasDisk) ProtoReflect() protoreflect.Message {
	mi := &file_nas_disk_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNasDiskRes_NasDisk.ProtoReflect.Descriptor instead.
func (*CreateNasDiskRes_NasDisk) Descriptor() ([]byte, []int) {
	return file_nas_disk_proto_rawDescGZIP(), []int{0, 0}
}

func (x *CreateNasDiskRes_NasDisk) GetShareId() string {
	if x != nil {
		return x.ShareId
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareName() string {
	if x != nil {
		return x.ShareName
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareDesc() string {
	if x != nil {
		return x.ShareDesc
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareProto() string {
	if x != nil {
		return x.ShareProto
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareSizeInG() int32 {
	if x != nil {
		return x.ShareSizeInG
	}
	return 0
}

func (x *CreateNasDiskRes_NasDisk) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetNetworkId() string {
	if x != nil {
		return x.NetworkId
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareStatus() string {
	if x != nil {
		return x.ShareStatus
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareProgress() string {
	if x != nil {
		return x.ShareProgress
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareServerId() string {
	if x != nil {
		return x.ShareServerId
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareServerHost() string {
	if x != nil {
		return x.ShareServerHost
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetShareNetworkId() string {
	if x != nil {
		return x.ShareNetworkId
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetCreatedTime() string {
	if x != nil {
		return x.CreatedTime
	}
	return ""
}

func (x *CreateNasDiskRes_NasDisk) GetUpdatedTime() string {
	if x != nil {
		return x.UpdatedTime
	}
	return ""
}

var File_nas_disk_proto protoreflect.FileDescriptor

var file_nas_disk_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6e, 0x61, 0x73, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x6e, 0x61, 0x73, 0x64, 0x69, 0x73, 0x6b, 0x22, 0xe8, 0x04, 0x0a, 0x10, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x12, 0x3c, 0x0a, 0x08, 0x6e, 0x61, 0x73, 0x5f, 0x64, 0x69, 0x73, 0x6b,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6e, 0x61, 0x73, 0x64, 0x69, 0x73, 0x6b,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65,
	0x73, 0x2e, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x07, 0x6e, 0x61, 0x73, 0x44, 0x69,
	0x73, 0x6b, 0x1a, 0xef, 0x03, 0x0a, 0x07, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x19,
	0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x68, 0x61, 0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x44, 0x65, 0x73, 0x63, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x25, 0x0a, 0x0f, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x49, 0x6e, 0x47, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x26, 0x0a, 0x0f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x48, 0x6f, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0xad, 0x02, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e,
	0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65,
	0x79, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x27,
	0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x55, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f,
	0x64, 0x65, 0x73, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x44, 0x65, 0x73, 0x63, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x25, 0x0a, 0x0f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x67, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x49, 0x6e, 0x47, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x49, 0x64, 0x22, 0x8b, 0x01, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e,
	0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65,
	0x79, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x27,
	0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x55, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x49, 0x64, 0x22, 0x76, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44,
	0x69, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x08,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x8d, 0x01, 0x0a, 0x12, 0x47,
	0x65, 0x74, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12,
	0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x65, 0x49, 0x64, 0x22, 0x5b, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x32, 0xeb, 0x01, 0x0a, 0x0e, 0x4e, 0x61, 0x73, 0x44,
	0x69, 0x73, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x19, 0x2e, 0x6e, 0x61,
	0x73, 0x64, 0x69, 0x73, 0x6b, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44,
	0x69, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x6e, 0x61, 0x73, 0x64, 0x69, 0x73, 0x6b,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65,
	0x73, 0x12, 0x45, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44, 0x69,
	0x73, 0x6b, 0x12, 0x19, 0x2e, 0x6e, 0x61, 0x73, 0x64, 0x69, 0x73, 0x6b, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x4e, 0x61, 0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e,
	0x6e, 0x61, 0x73, 0x64, 0x69, 0x73, 0x6b, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x61,
	0x73, 0x44, 0x69, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x12, 0x4b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4d,
	0x6f, 0x75, 0x6e, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1b, 0x2e, 0x6e, 0x61,
	0x73, 0x64, 0x69, 0x73, 0x6b, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x6e, 0x61, 0x73, 0x64, 0x69,
	0x73, 0x6b, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nas_disk_proto_rawDescOnce sync.Once
	file_nas_disk_proto_rawDescData = file_nas_disk_proto_rawDesc
)

func file_nas_disk_proto_rawDescGZIP() []byte {
	file_nas_disk_proto_rawDescOnce.Do(func() {
		file_nas_disk_proto_rawDescData = protoimpl.X.CompressGZIP(file_nas_disk_proto_rawDescData)
	})
	return file_nas_disk_proto_rawDescData
}

var file_nas_disk_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_nas_disk_proto_goTypes = []interface{}{
	(*CreateNasDiskRes)(nil),         // 0: nasdisk.CreateNasDiskRes
	(*CreateNasDiskReq)(nil),         // 1: nasdisk.CreateNasDiskReq
	(*DeleteNasDiskReq)(nil),         // 2: nasdisk.DeleteNasDiskReq
	(*DeleteNasDiskRes)(nil),         // 3: nasdisk.DeleteNasDiskRes
	(*GetMountClientsReq)(nil),       // 4: nasdisk.GetMountClientsReq
	(*GetMountClientsRes)(nil),       // 5: nasdisk.GetMountClientsRes
	(*CreateNasDiskRes_NasDisk)(nil), // 6: nasdisk.CreateNasDiskRes.NasDisk
}
var file_nas_disk_proto_depIdxs = []int32{
	6, // 0: nasdisk.CreateNasDiskRes.nas_disk:type_name -> nasdisk.CreateNasDiskRes.NasDisk
	1, // 1: nasdisk.NasDiskService.CreateNasDisk:input_type -> nasdisk.CreateNasDiskReq
	2, // 2: nasdisk.NasDiskService.DeleteNasDisk:input_type -> nasdisk.DeleteNasDiskReq
	4, // 3: nasdisk.NasDiskService.GetMountClients:input_type -> nasdisk.GetMountClientsReq
	0, // 4: nasdisk.NasDiskService.CreateNasDisk:output_type -> nasdisk.CreateNasDiskRes
	3, // 5: nasdisk.NasDiskService.DeleteNasDisk:output_type -> nasdisk.DeleteNasDiskRes
	5, // 6: nasdisk.NasDiskService.GetMountClients:output_type -> nasdisk.GetMountClientsRes
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_nas_disk_proto_init() }
func file_nas_disk_proto_init() {
	if File_nas_disk_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nas_disk_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNasDiskRes); i {
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
		file_nas_disk_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNasDiskReq); i {
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
		file_nas_disk_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNasDiskReq); i {
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
		file_nas_disk_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNasDiskRes); i {
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
		file_nas_disk_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMountClientsReq); i {
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
		file_nas_disk_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMountClientsRes); i {
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
		file_nas_disk_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNasDiskRes_NasDisk); i {
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
			RawDescriptor: file_nas_disk_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nas_disk_proto_goTypes,
		DependencyIndexes: file_nas_disk_proto_depIdxs,
		MessageInfos:      file_nas_disk_proto_msgTypes,
	}.Build()
	File_nas_disk_proto = out.File
	file_nas_disk_proto_rawDesc = nil
	file_nas_disk_proto_goTypes = nil
	file_nas_disk_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NasDiskServiceClient is the client API for NasDiskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NasDiskServiceClient interface {
	//创建NAS盘
	CreateNasDisk(ctx context.Context, in *CreateNasDiskReq, opts ...grpc.CallOption) (*CreateNasDiskRes, error)
	//删除NAS盘
	DeleteNasDisk(ctx context.Context, in *DeleteNasDiskReq, opts ...grpc.CallOption) (*DeleteNasDiskRes, error)
	//查看挂载客户端
	GetMountClients(ctx context.Context, in *GetMountClientsReq, opts ...grpc.CallOption) (*GetMountClientsRes, error)
}

type nasDiskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNasDiskServiceClient(cc grpc.ClientConnInterface) NasDiskServiceClient {
	return &nasDiskServiceClient{cc}
}

func (c *nasDiskServiceClient) CreateNasDisk(ctx context.Context, in *CreateNasDiskReq, opts ...grpc.CallOption) (*CreateNasDiskRes, error) {
	out := new(CreateNasDiskRes)
	err := c.cc.Invoke(ctx, "/nasdisk.NasDiskService/CreateNasDisk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nasDiskServiceClient) DeleteNasDisk(ctx context.Context, in *DeleteNasDiskReq, opts ...grpc.CallOption) (*DeleteNasDiskRes, error) {
	out := new(DeleteNasDiskRes)
	err := c.cc.Invoke(ctx, "/nasdisk.NasDiskService/DeleteNasDisk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nasDiskServiceClient) GetMountClients(ctx context.Context, in *GetMountClientsReq, opts ...grpc.CallOption) (*GetMountClientsRes, error) {
	out := new(GetMountClientsRes)
	err := c.cc.Invoke(ctx, "/nasdisk.NasDiskService/GetMountClients", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NasDiskServiceServer is the server API for NasDiskService service.
type NasDiskServiceServer interface {
	//创建NAS盘
	CreateNasDisk(context.Context, *CreateNasDiskReq) (*CreateNasDiskRes, error)
	//删除NAS盘
	DeleteNasDisk(context.Context, *DeleteNasDiskReq) (*DeleteNasDiskRes, error)
	//查看挂载客户端
	GetMountClients(context.Context, *GetMountClientsReq) (*GetMountClientsRes, error)
}

// UnimplementedNasDiskServiceServer can be embedded to have forward compatible implementations.
type UnimplementedNasDiskServiceServer struct {
}

func (*UnimplementedNasDiskServiceServer) CreateNasDisk(context.Context, *CreateNasDiskReq) (*CreateNasDiskRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNasDisk not implemented")
}
func (*UnimplementedNasDiskServiceServer) DeleteNasDisk(context.Context, *DeleteNasDiskReq) (*DeleteNasDiskRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNasDisk not implemented")
}
func (*UnimplementedNasDiskServiceServer) GetMountClients(context.Context, *GetMountClientsReq) (*GetMountClientsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMountClients not implemented")
}

func RegisterNasDiskServiceServer(s *grpc.Server, srv NasDiskServiceServer) {
	s.RegisterService(&_NasDiskService_serviceDesc, srv)
}

func _NasDiskService_CreateNasDisk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNasDiskReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NasDiskServiceServer).CreateNasDisk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nasdisk.NasDiskService/CreateNasDisk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NasDiskServiceServer).CreateNasDisk(ctx, req.(*CreateNasDiskReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NasDiskService_DeleteNasDisk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNasDiskReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NasDiskServiceServer).DeleteNasDisk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nasdisk.NasDiskService/DeleteNasDisk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NasDiskServiceServer).DeleteNasDisk(ctx, req.(*DeleteNasDiskReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NasDiskService_GetMountClients_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMountClientsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NasDiskServiceServer).GetMountClients(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nasdisk.NasDiskService/GetMountClients",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NasDiskServiceServer).GetMountClients(ctx, req.(*GetMountClientsReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _NasDiskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nasdisk.NasDiskService",
	HandlerType: (*NasDiskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNasDisk",
			Handler:    _NasDiskService_CreateNasDisk_Handler,
		},
		{
			MethodName: "DeleteNasDisk",
			Handler:    _NasDiskService_DeleteNasDisk_Handler,
		},
		{
			MethodName: "GetMountClients",
			Handler:    _NasDiskService_GetMountClients_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nas_disk.proto",
}