// 指定的当前proto语法的版本，有2和3

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: image.proto

// 指定文件生成出来的package

package image

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

type Image struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImageId              string  `protobuf:"bytes,1,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
	ImageName            string  `protobuf:"bytes,2,opt,name=image_name,json=imageName,proto3" json:"image_name,omitempty"`
	ImageDiskformat      string  `protobuf:"bytes,3,opt,name=image_diskformat,json=imageDiskformat,proto3" json:"image_diskformat,omitempty"`
	ImageContainerformat string  `protobuf:"bytes,4,opt,name=image_containerformat,json=imageContainerformat,proto3" json:"image_containerformat,omitempty"`
	ImageSizeInG         float32 `protobuf:"fixed32,5,opt,name=image_size_in_g,json=imageSizeInG,proto3" json:"image_size_in_g,omitempty"`
}

func (x *Image) Reset() {
	*x = Image{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Image) ProtoMessage() {}

func (x *Image) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Image.ProtoReflect.Descriptor instead.
func (*Image) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{0}
}

func (x *Image) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

func (x *Image) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

func (x *Image) GetImageDiskformat() string {
	if x != nil {
		return x.ImageDiskformat
	}
	return ""
}

func (x *Image) GetImageContainerformat() string {
	if x != nil {
		return x.ImageContainerformat
	}
	return ""
}

func (x *Image) GetImageSizeInG() float32 {
	if x != nil {
		return x.ImageSizeInG
	}
	return 0
}

type ListImagesReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Apikey         string `protobuf:"bytes,1,opt,name=apikey,proto3" json:"apikey,omitempty"`
	TenantId       string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	PlatformUserid string `protobuf:"bytes,3,opt,name=platform_userid,json=platformUserid,proto3" json:"platform_userid,omitempty"`
	StartImageId   string `protobuf:"bytes,4,opt,name=start_image_id,json=startImageId,proto3" json:"start_image_id,omitempty"`
	PageSize       int32  `protobuf:"varint,5,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *ListImagesReq) Reset() {
	*x = ListImagesReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListImagesReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListImagesReq) ProtoMessage() {}

func (x *ListImagesReq) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListImagesReq.ProtoReflect.Descriptor instead.
func (*ListImagesReq) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{1}
}

func (x *ListImagesReq) GetApikey() string {
	if x != nil {
		return x.Apikey
	}
	return ""
}

func (x *ListImagesReq) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *ListImagesReq) GetPlatformUserid() string {
	if x != nil {
		return x.PlatformUserid
	}
	return ""
}

func (x *ListImagesReq) GetStartImageId() string {
	if x != nil {
		return x.StartImageId
	}
	return ""
}

func (x *ListImagesReq) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ListImagesRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code        int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg         string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Images      []*Image `protobuf:"bytes,3,rep,name=images,proto3" json:"images,omitempty"`
	NextImageId string   `protobuf:"bytes,4,opt,name=next_image_id,json=nextImageId,proto3" json:"next_image_id,omitempty"`
}

func (x *ListImagesRes) Reset() {
	*x = ListImagesRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListImagesRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListImagesRes) ProtoMessage() {}

func (x *ListImagesRes) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListImagesRes.ProtoReflect.Descriptor instead.
func (*ListImagesRes) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{2}
}

func (x *ListImagesRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ListImagesRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *ListImagesRes) GetImages() []*Image {
	if x != nil {
		return x.Images
	}
	return nil
}

func (x *ListImagesRes) GetNextImageId() string {
	if x != nil {
		return x.NextImageId
	}
	return ""
}

type GetImageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Apikey         string `protobuf:"bytes,1,opt,name=apikey,proto3" json:"apikey,omitempty"`
	TenantId       string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	PlatformUserid string `protobuf:"bytes,3,opt,name=platform_userid,json=platformUserid,proto3" json:"platform_userid,omitempty"`
	ImageId        string `protobuf:"bytes,4,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty"`
}

func (x *GetImageReq) Reset() {
	*x = GetImageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageReq) ProtoMessage() {}

func (x *GetImageReq) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageReq.ProtoReflect.Descriptor instead.
func (*GetImageReq) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{3}
}

func (x *GetImageReq) GetApikey() string {
	if x != nil {
		return x.Apikey
	}
	return ""
}

func (x *GetImageReq) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *GetImageReq) GetPlatformUserid() string {
	if x != nil {
		return x.PlatformUserid
	}
	return ""
}

func (x *GetImageReq) GetImageId() string {
	if x != nil {
		return x.ImageId
	}
	return ""
}

type GetImageRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code  int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg   string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Image *Image `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *GetImageRes) Reset() {
	*x = GetImageRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_image_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetImageRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetImageRes) ProtoMessage() {}

func (x *GetImageRes) ProtoReflect() protoreflect.Message {
	mi := &file_image_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetImageRes.ProtoReflect.Descriptor instead.
func (*GetImageRes) Descriptor() ([]byte, []int) {
	return file_image_proto_rawDescGZIP(), []int{4}
}

func (x *GetImageRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GetImageRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GetImageRes) GetImage() *Image {
	if x != nil {
		return x.Image
	}
	return nil
}

var File_image_proto protoreflect.FileDescriptor

var file_image_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x22, 0xc8, 0x01, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x69, 0x73, 0x6b, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x12, 0x33, 0x0a, 0x15, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x14, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x25, 0x0a, 0x0f, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x0c, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x49, 0x6e, 0x47, 0x22,
	0xb0, 0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12,
	0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x22, 0x7f, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x24, 0x0a, 0x06, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12,
	0x22, 0x0a, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x65, 0x78, 0x74, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x49, 0x64, 0x22, 0x86, 0x01, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x6b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x69,
	0x64, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x64, 0x22, 0x57, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x12, 0x22, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x05,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x32, 0x7c, 0x0a, 0x0c, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x73, 0x12, 0x14, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x12,
	0x32, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x12, 0x2e, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a,
	0x12, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_image_proto_rawDescOnce sync.Once
	file_image_proto_rawDescData = file_image_proto_rawDesc
)

func file_image_proto_rawDescGZIP() []byte {
	file_image_proto_rawDescOnce.Do(func() {
		file_image_proto_rawDescData = protoimpl.X.CompressGZIP(file_image_proto_rawDescData)
	})
	return file_image_proto_rawDescData
}

var file_image_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_image_proto_goTypes = []interface{}{
	(*Image)(nil),         // 0: image.Image
	(*ListImagesReq)(nil), // 1: image.ListImagesReq
	(*ListImagesRes)(nil), // 2: image.ListImagesRes
	(*GetImageReq)(nil),   // 3: image.GetImageReq
	(*GetImageRes)(nil),   // 4: image.GetImageRes
}
var file_image_proto_depIdxs = []int32{
	0, // 0: image.ListImagesRes.images:type_name -> image.Image
	0, // 1: image.GetImageRes.image:type_name -> image.Image
	1, // 2: image.ImageService.ListImages:input_type -> image.ListImagesReq
	3, // 3: image.ImageService.GetImage:input_type -> image.GetImageReq
	2, // 4: image.ImageService.ListImages:output_type -> image.ListImagesRes
	4, // 5: image.ImageService.GetImage:output_type -> image.GetImageRes
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_image_proto_init() }
func file_image_proto_init() {
	if File_image_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_image_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Image); i {
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
		file_image_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListImagesReq); i {
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
		file_image_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListImagesRes); i {
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
		file_image_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageReq); i {
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
		file_image_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetImageRes); i {
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
			RawDescriptor: file_image_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_image_proto_goTypes,
		DependencyIndexes: file_image_proto_depIdxs,
		MessageInfos:      file_image_proto_msgTypes,
	}.Build()
	File_image_proto = out.File
	file_image_proto_rawDesc = nil
	file_image_proto_goTypes = nil
	file_image_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ImageServiceClient is the client API for ImageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImageServiceClient interface {
	//获取镜像列表
	ListImages(ctx context.Context, in *ListImagesReq, opts ...grpc.CallOption) (*ListImagesRes, error)
	//获取镜像信息
	GetImage(ctx context.Context, in *GetImageReq, opts ...grpc.CallOption) (*GetImageRes, error)
}

type imageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewImageServiceClient(cc grpc.ClientConnInterface) ImageServiceClient {
	return &imageServiceClient{cc}
}

func (c *imageServiceClient) ListImages(ctx context.Context, in *ListImagesReq, opts ...grpc.CallOption) (*ListImagesRes, error) {
	out := new(ListImagesRes)
	err := c.cc.Invoke(ctx, "/image.ImageService/ListImages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageServiceClient) GetImage(ctx context.Context, in *GetImageReq, opts ...grpc.CallOption) (*GetImageRes, error) {
	out := new(GetImageRes)
	err := c.cc.Invoke(ctx, "/image.ImageService/GetImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImageServiceServer is the server API for ImageService service.
type ImageServiceServer interface {
	//获取镜像列表
	ListImages(context.Context, *ListImagesReq) (*ListImagesRes, error)
	//获取镜像信息
	GetImage(context.Context, *GetImageReq) (*GetImageRes, error)
}

// UnimplementedImageServiceServer can be embedded to have forward compatible implementations.
type UnimplementedImageServiceServer struct {
}

func (*UnimplementedImageServiceServer) ListImages(context.Context, *ListImagesReq) (*ListImagesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListImages not implemented")
}
func (*UnimplementedImageServiceServer) GetImage(context.Context, *GetImageReq) (*GetImageRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetImage not implemented")
}

func RegisterImageServiceServer(s *grpc.Server, srv ImageServiceServer) {
	s.RegisterService(&_ImageService_serviceDesc, srv)
}

func _ImageService_ListImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListImagesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).ListImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image.ImageService/ListImages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).ListImages(ctx, req.(*ListImagesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageService_GetImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetImageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).GetImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/image.ImageService/GetImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).GetImage(ctx, req.(*GetImageReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _ImageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "image.ImageService",
	HandlerType: (*ImageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListImages",
			Handler:    _ImageService_ListImages_Handler,
		},
		{
			MethodName: "GetImage",
			Handler:    _ImageService_GetImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "image.proto",
}
