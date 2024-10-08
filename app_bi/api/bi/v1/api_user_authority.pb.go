// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: bi/v1/api_user_authority.proto

package v1

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

type ReqSetUserAuth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	AuthorityId uint64 `protobuf:"varint,2,opt,name=authorityId,proto3" json:"authorityId,omitempty"`
}

func (x *ReqSetUserAuth) Reset() {
	*x = ReqSetUserAuth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bi_v1_api_user_authority_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqSetUserAuth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqSetUserAuth) ProtoMessage() {}

func (x *ReqSetUserAuth) ProtoReflect() protoreflect.Message {
	mi := &file_bi_v1_api_user_authority_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqSetUserAuth.ProtoReflect.Descriptor instead.
func (*ReqSetUserAuth) Descriptor() ([]byte, []int) {
	return file_bi_v1_api_user_authority_proto_rawDescGZIP(), []int{0}
}

func (x *ReqSetUserAuth) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReqSetUserAuth) GetAuthorityId() uint64 {
	if x != nil {
		return x.AuthorityId
	}
	return 0
}

type ReqSetUserAuthorities struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID uint64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
	// 1:添加; 2:删除
	Type        uint64   `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	AuthorityId []uint64 `protobuf:"varint,3,rep,packed,name=authorityId,proto3" json:"authorityId,omitempty"`
}

func (x *ReqSetUserAuthorities) Reset() {
	*x = ReqSetUserAuthorities{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bi_v1_api_user_authority_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqSetUserAuthorities) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqSetUserAuthorities) ProtoMessage() {}

func (x *ReqSetUserAuthorities) ProtoReflect() protoreflect.Message {
	mi := &file_bi_v1_api_user_authority_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqSetUserAuthorities.ProtoReflect.Descriptor instead.
func (*ReqSetUserAuthorities) Descriptor() ([]byte, []int) {
	return file_bi_v1_api_user_authority_proto_rawDescGZIP(), []int{1}
}

func (x *ReqSetUserAuthorities) GetUserID() uint64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *ReqSetUserAuthorities) GetType() uint64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *ReqSetUserAuthorities) GetAuthorityId() []uint64 {
	if x != nil {
		return x.AuthorityId
	}
	return nil
}

type Authority struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthorityId   uint64 `protobuf:"varint,1,opt,name=authorityId,proto3" json:"authorityId,omitempty"`
	AuthorityName string `protobuf:"bytes,2,opt,name=authorityName,proto3" json:"authorityName,omitempty"`
	ParentId      uint64 `protobuf:"varint,3,opt,name=parentId,proto3" json:"parentId,omitempty"`
}

func (x *Authority) Reset() {
	*x = Authority{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bi_v1_api_user_authority_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Authority) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Authority) ProtoMessage() {}

func (x *Authority) ProtoReflect() protoreflect.Message {
	mi := &file_bi_v1_api_user_authority_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Authority.ProtoReflect.Descriptor instead.
func (*Authority) Descriptor() ([]byte, []int) {
	return file_bi_v1_api_user_authority_proto_rawDescGZIP(), []int{2}
}

func (x *Authority) GetAuthorityId() uint64 {
	if x != nil {
		return x.AuthorityId
	}
	return 0
}

func (x *Authority) GetAuthorityName() string {
	if x != nil {
		return x.AuthorityName
	}
	return ""
}

func (x *Authority) GetParentId() uint64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

var File_bi_v1_api_user_authority_proto protoreflect.FileDescriptor

var file_bi_v1_api_user_authority_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x62, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x62, 0x69, 0x2e, 0x76, 0x31, 0x22, 0x4a, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x53, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x49, 0x64, 0x22, 0x65, 0x0a, 0x15, 0x52, 0x65, 0x71, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x49, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0b, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x49, 0x64, 0x22, 0x6f, 0x0a, 0x09, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x42, 0x56, 0x0a, 0x14, 0x64,
	0x65, 0x76, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69,
	0x2e, 0x76, 0x31, 0x42, 0x09, 0x42, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x56, 0x31, 0x50, 0x01,
	0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x65, 0x61,
	0x66, 0x2d, 0x72, 0x61, 0x69, 0x6e, 0x2f, 0x72, 0x61, 0x69, 0x6e, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x61, 0x70, 0x70, 0x5f, 0x62, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x2f, 0x76, 0x31,
	0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bi_v1_api_user_authority_proto_rawDescOnce sync.Once
	file_bi_v1_api_user_authority_proto_rawDescData = file_bi_v1_api_user_authority_proto_rawDesc
)

func file_bi_v1_api_user_authority_proto_rawDescGZIP() []byte {
	file_bi_v1_api_user_authority_proto_rawDescOnce.Do(func() {
		file_bi_v1_api_user_authority_proto_rawDescData = protoimpl.X.CompressGZIP(file_bi_v1_api_user_authority_proto_rawDescData)
	})
	return file_bi_v1_api_user_authority_proto_rawDescData
}

var file_bi_v1_api_user_authority_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_bi_v1_api_user_authority_proto_goTypes = []any{
	(*ReqSetUserAuth)(nil),        // 0: bi.v1.ReqSetUserAuth
	(*ReqSetUserAuthorities)(nil), // 1: bi.v1.ReqSetUserAuthorities
	(*Authority)(nil),             // 2: bi.v1.Authority
}
var file_bi_v1_api_user_authority_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_bi_v1_api_user_authority_proto_init() }
func file_bi_v1_api_user_authority_proto_init() {
	if File_bi_v1_api_user_authority_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bi_v1_api_user_authority_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ReqSetUserAuth); i {
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
		file_bi_v1_api_user_authority_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ReqSetUserAuthorities); i {
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
		file_bi_v1_api_user_authority_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Authority); i {
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
			RawDescriptor: file_bi_v1_api_user_authority_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bi_v1_api_user_authority_proto_goTypes,
		DependencyIndexes: file_bi_v1_api_user_authority_proto_depIdxs,
		MessageInfos:      file_bi_v1_api_user_authority_proto_msgTypes,
	}.Build()
	File_bi_v1_api_user_authority_proto = out.File
	file_bi_v1_api_user_authority_proto_rawDesc = nil
	file_bi_v1_api_user_authority_proto_goTypes = nil
	file_bi_v1_api_user_authority_proto_depIdxs = nil
}
