// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.19.4
// source: proto/connect.proto

package proto

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

type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver       int64  `protobuf:"varint,1,opt,name=Ver,proto3" json:"Ver,omitempty"`
	Operation int64  `protobuf:"varint,2,opt,name=Operation,proto3" json:"Operation,omitempty"`
	SeqId     string `protobuf:"bytes,3,opt,name=SeqId,proto3" json:"SeqId,omitempty"`
	Body      []byte `protobuf:"bytes,4,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_connect_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_connect_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_proto_connect_proto_rawDescGZIP(), []int{0}
}

func (x *Msg) GetVer() int64 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *Msg) GetOperation() int64 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *Msg) GetSeqId() string {
	if x != nil {
		return x.SeqId
	}
	return ""
}

func (x *Msg) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type PushMsgRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64  `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Msg    []*Msg `protobuf:"bytes,2,rep,name=Msg,proto3" json:"Msg,omitempty"`
}

func (x *PushMsgRequest) Reset() {
	*x = PushMsgRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_connect_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgRequest) ProtoMessage() {}

func (x *PushMsgRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_connect_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgRequest.ProtoReflect.Descriptor instead.
func (*PushMsgRequest) Descriptor() ([]byte, []int) {
	return file_proto_connect_proto_rawDescGZIP(), []int{1}
}

func (x *PushMsgRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *PushMsgRequest) GetMsg() []*Msg {
	if x != nil {
		return x.Msg
	}
	return nil
}

type PushRoomMsgRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Room int64  `protobuf:"varint,1,opt,name=Room,proto3" json:"Room,omitempty"`
	Msg  []*Msg `protobuf:"bytes,2,rep,name=Msg,proto3" json:"Msg,omitempty"`
}

func (x *PushRoomMsgRequest) Reset() {
	*x = PushRoomMsgRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_connect_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushRoomMsgRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRoomMsgRequest) ProtoMessage() {}

func (x *PushRoomMsgRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_connect_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRoomMsgRequest.ProtoReflect.Descriptor instead.
func (*PushRoomMsgRequest) Descriptor() ([]byte, []int) {
	return file_proto_connect_proto_rawDescGZIP(), []int{2}
}

func (x *PushRoomMsgRequest) GetRoom() int64 {
	if x != nil {
		return x.Room
	}
	return 0
}

func (x *PushRoomMsgRequest) GetMsg() []*Msg {
	if x != nil {
		return x.Msg
	}
	return nil
}

type PushRoomCountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Room  int64 `protobuf:"varint,1,opt,name=Room,proto3" json:"Room,omitempty"`
	Count int64 `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"`
}

func (x *PushRoomCountRequest) Reset() {
	*x = PushRoomCountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_connect_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushRoomCountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRoomCountRequest) ProtoMessage() {}

func (x *PushRoomCountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_connect_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRoomCountRequest.ProtoReflect.Descriptor instead.
func (*PushRoomCountRequest) Descriptor() ([]byte, []int) {
	return file_proto_connect_proto_rawDescGZIP(), []int{3}
}

func (x *PushRoomCountRequest) GetRoom() int64 {
	if x != nil {
		return x.Room
	}
	return 0
}

func (x *PushRoomCountRequest) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_proto_connect_proto protoreflect.FileDescriptor

var file_proto_connect_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5f, 0x0a, 0x03,
	0x4d, 0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x56, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x56, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x65, 0x71, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x53, 0x65, 0x71, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x46, 0x0a,
	0x0e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x73, 0x67,
	0x52, 0x03, 0x4d, 0x73, 0x67, 0x22, 0x46, 0x0a, 0x12, 0x50, 0x75, 0x73, 0x68, 0x52, 0x6f, 0x6f,
	0x6d, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x52,
	0x6f, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x52, 0x6f, 0x6f, 0x6d, 0x12,
	0x1c, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x73, 0x67, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x22, 0x40, 0x0a,
	0x14, 0x50, 0x75, 0x73, 0x68, 0x52, 0x6f, 0x6f, 0x6d, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x6f, 0x6f, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42,
	0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_connect_proto_rawDescOnce sync.Once
	file_proto_connect_proto_rawDescData = file_proto_connect_proto_rawDesc
)

func file_proto_connect_proto_rawDescGZIP() []byte {
	file_proto_connect_proto_rawDescOnce.Do(func() {
		file_proto_connect_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_connect_proto_rawDescData)
	})
	return file_proto_connect_proto_rawDescData
}

var file_proto_connect_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_connect_proto_goTypes = []interface{}{
	(*Msg)(nil),                  // 0: proto.Msg
	(*PushMsgRequest)(nil),       // 1: proto.PushMsgRequest
	(*PushRoomMsgRequest)(nil),   // 2: proto.PushRoomMsgRequest
	(*PushRoomCountRequest)(nil), // 3: proto.PushRoomCountRequest
}
var file_proto_connect_proto_depIdxs = []int32{
	0, // 0: proto.PushMsgRequest.Msg:type_name -> proto.Msg
	0, // 1: proto.PushRoomMsgRequest.Msg:type_name -> proto.Msg
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_connect_proto_init() }
func file_proto_connect_proto_init() {
	if File_proto_connect_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_connect_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
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
		file_proto_connect_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgRequest); i {
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
		file_proto_connect_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushRoomMsgRequest); i {
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
		file_proto_connect_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushRoomCountRequest); i {
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
			RawDescriptor: file_proto_connect_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_connect_proto_goTypes,
		DependencyIndexes: file_proto_connect_proto_depIdxs,
		MessageInfos:      file_proto_connect_proto_msgTypes,
	}.Build()
	File_proto_connect_proto = out.File
	file_proto_connect_proto_rawDesc = nil
	file_proto_connect_proto_goTypes = nil
	file_proto_connect_proto_depIdxs = nil
}
