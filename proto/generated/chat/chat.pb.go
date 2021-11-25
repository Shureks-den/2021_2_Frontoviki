// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: chat.proto

package chat

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

type DialogIdentifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdFrom int64 `protobuf:"varint,1,opt,name=idFrom,proto3" json:"idFrom,omitempty"`
	IdTo   int64 `protobuf:"varint,2,opt,name=idTo,proto3" json:"idTo,omitempty"`
	IdAdv  int64 `protobuf:"varint,3,opt,name=idAdv,proto3" json:"idAdv,omitempty"`
}

func (x *DialogIdentifier) Reset() {
	*x = DialogIdentifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DialogIdentifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DialogIdentifier) ProtoMessage() {}

func (x *DialogIdentifier) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DialogIdentifier.ProtoReflect.Descriptor instead.
func (*DialogIdentifier) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *DialogIdentifier) GetIdFrom() int64 {
	if x != nil {
		return x.IdFrom
	}
	return 0
}

func (x *DialogIdentifier) GetIdTo() int64 {
	if x != nil {
		return x.IdTo
	}
	return 0
}

func (x *DialogIdentifier) GetIdAdv() int64 {
	if x != nil {
		return x.IdAdv
	}
	return 0
}

type UserIdentifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdFrom int64 `protobuf:"varint,1,opt,name=idFrom,proto3" json:"idFrom,omitempty"`
}

func (x *UserIdentifier) Reset() {
	*x = UserIdentifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIdentifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIdentifier) ProtoMessage() {}

func (x *UserIdentifier) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIdentifier.ProtoReflect.Descriptor instead.
func (*UserIdentifier) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *UserIdentifier) GetIdFrom() int64 {
	if x != nil {
		return x.IdFrom
	}
	return 0
}

type FilterParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset int64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  int64 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *FilterParams) Reset() {
	*x = FilterParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilterParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterParams) ProtoMessage() {}

func (x *FilterParams) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterParams.ProtoReflect.Descriptor instead.
func (*FilterParams) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{2}
}

func (x *FilterParams) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *FilterParams) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdFrom    int64                  `protobuf:"varint,1,opt,name=IdFrom,proto3" json:"IdFrom,omitempty"`
	IdTo      int64                  `protobuf:"varint,2,opt,name=IdTo,proto3" json:"IdTo,omitempty"`
	IdAdv     int64                  `protobuf:"varint,3,opt,name=IdAdv,proto3" json:"IdAdv,omitempty"`
	Msg       string                 `protobuf:"bytes,4,opt,name=Msg,proto3" json:"Msg,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{3}
}

func (x *Message) GetIdFrom() int64 {
	if x != nil {
		return x.IdFrom
	}
	return 0
}

func (x *Message) GetIdTo() int64 {
	if x != nil {
		return x.IdTo
	}
	return 0
}

func (x *Message) GetIdAdv() int64 {
	if x != nil {
		return x.IdAdv
	}
	return 0
}

func (x *Message) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Message) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type Nothing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dummy bool `protobuf:"varint,1,opt,name=dummy,proto3" json:"dummy,omitempty"`
}

func (x *Nothing) Reset() {
	*x = Nothing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nothing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nothing) ProtoMessage() {}

func (x *Nothing) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nothing.ProtoReflect.Descriptor instead.
func (*Nothing) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{4}
}

func (x *Nothing) GetDummy() bool {
	if x != nil {
		return x.Dummy
	}
	return false
}

type GetHistoryArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DI *DialogIdentifier `protobuf:"bytes,1,opt,name=DI,proto3" json:"DI,omitempty"`
	FP *FilterParams     `protobuf:"bytes,2,opt,name=FP,proto3" json:"FP,omitempty"`
}

func (x *GetHistoryArg) Reset() {
	*x = GetHistoryArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetHistoryArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetHistoryArg) ProtoMessage() {}

func (x *GetHistoryArg) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetHistoryArg.ProtoReflect.Descriptor instead.
func (*GetHistoryArg) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{5}
}

func (x *GetHistoryArg) GetDI() *DialogIdentifier {
	if x != nil {
		return x.DI
	}
	return nil
}

func (x *GetHistoryArg) GetFP() *FilterParams {
	if x != nil {
		return x.FP
	}
	return nil
}

type Dialog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id1       int64                  `protobuf:"varint,1,opt,name=Id1,proto3" json:"Id1,omitempty"`
	Id2       int64                  `protobuf:"varint,2,opt,name=Id2,proto3" json:"Id2,omitempty"`
	IdAdv     int64                  `protobuf:"varint,3,opt,name=IdAdv,proto3" json:"IdAdv,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
}

func (x *Dialog) Reset() {
	*x = Dialog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dialog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dialog) ProtoMessage() {}

func (x *Dialog) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dialog.ProtoReflect.Descriptor instead.
func (*Dialog) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{6}
}

func (x *Dialog) GetId1() int64 {
	if x != nil {
		return x.Id1
	}
	return 0
}

func (x *Dialog) GetId2() int64 {
	if x != nil {
		return x.Id2
	}
	return 0
}

func (x *Dialog) GetIdAdv() int64 {
	if x != nil {
		return x.IdAdv
	}
	return 0
}

func (x *Dialog) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type Messages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	M []*Message `protobuf:"bytes,1,rep,name=m,proto3" json:"m,omitempty"`
}

func (x *Messages) Reset() {
	*x = Messages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Messages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Messages) ProtoMessage() {}

func (x *Messages) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Messages.ProtoReflect.Descriptor instead.
func (*Messages) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{7}
}

func (x *Messages) GetM() []*Message {
	if x != nil {
		return x.M
	}
	return nil
}

type Dialogs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	D []*Dialog `protobuf:"bytes,1,rep,name=d,proto3" json:"d,omitempty"`
}

func (x *Dialogs) Reset() {
	*x = Dialogs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dialogs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dialogs) ProtoMessage() {}

func (x *Dialogs) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dialogs.ProtoReflect.Descriptor instead.
func (*Dialogs) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{8}
}

func (x *Dialogs) GetD() []*Dialog {
	if x != nil {
		return x.D
	}
	return nil
}

var File_chat_proto protoreflect.FileDescriptor

var file_chat_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68,
	0x61, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x10, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x64, 0x46, 0x72, 0x6f,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x69, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x12,
	0x12, 0x0a, 0x04, 0x69, 0x64, 0x54, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x69,
	0x64, 0x54, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x64, 0x41, 0x64, 0x76, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x69, 0x64, 0x41, 0x64, 0x76, 0x22, 0x28, 0x0a, 0x0e, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x69,
	0x64, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x69, 0x64, 0x46,
	0x72, 0x6f, 0x6d, 0x22, 0x3c, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x22, 0x97, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x49, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x49,
	0x64, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x49, 0x64, 0x54, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x49, 0x64, 0x54, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x64, 0x41,
	0x64, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x49, 0x64, 0x41, 0x64, 0x76, 0x12,
	0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73,
	0x67, 0x12, 0x38, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x1f, 0x0a, 0x07, 0x4e,
	0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x64, 0x75, 0x6d, 0x6d, 0x79, 0x22, 0x5b, 0x0a, 0x0d,
	0x47, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x41, 0x72, 0x67, 0x12, 0x26, 0x0a,
	0x02, 0x44, 0x49, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x52, 0x02, 0x44, 0x49, 0x12, 0x22, 0x0a, 0x02, 0x46, 0x50, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x52, 0x02, 0x46, 0x50, 0x22, 0x7c, 0x0a, 0x06, 0x44, 0x69, 0x61,
	0x6c, 0x6f, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x49, 0x64, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x64, 0x32, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x49, 0x64, 0x32, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x64, 0x41, 0x64, 0x76,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x49, 0x64, 0x41, 0x64, 0x76, 0x12, 0x38, 0x0a,
	0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x27, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x01, 0x6d, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x01, 0x6d,
	0x22, 0x25, 0x0a, 0x07, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x73, 0x12, 0x1a, 0x0a, 0x01, 0x64,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x44, 0x69,
	0x61, 0x6c, 0x6f, 0x67, 0x52, 0x01, 0x64, 0x32, 0xc4, 0x01, 0x0a, 0x04, 0x43, 0x68, 0x61, 0x74,
	0x12, 0x31, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x13,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x41, 0x72, 0x67, 0x1a, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x0d, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x12, 0x2e, 0x0a, 0x05, 0x43,
	0x6c, 0x65, 0x61, 0x72, 0x12, 0x16, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x44, 0x69, 0x61, 0x6c,
	0x6f, 0x67, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x1a, 0x0d, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x12, 0x31, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x73, 0x12, 0x14, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x1a,
	0x0d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x0a,
	0x5a, 0x08, 0x2e, 0x2f, 0x2e, 0x3b, 0x63, 0x68, 0x61, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_chat_proto_rawDescOnce sync.Once
	file_chat_proto_rawDescData = file_chat_proto_rawDesc
)

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_proto_rawDescData)
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_chat_proto_goTypes = []interface{}{
	(*DialogIdentifier)(nil),      // 0: chat.DialogIdentifier
	(*UserIdentifier)(nil),        // 1: chat.UserIdentifier
	(*FilterParams)(nil),          // 2: chat.FilterParams
	(*Message)(nil),               // 3: chat.Message
	(*Nothing)(nil),               // 4: chat.Nothing
	(*GetHistoryArg)(nil),         // 5: chat.GetHistoryArg
	(*Dialog)(nil),                // 6: chat.Dialog
	(*Messages)(nil),              // 7: chat.Messages
	(*Dialogs)(nil),               // 8: chat.Dialogs
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_chat_proto_depIdxs = []int32{
	9,  // 0: chat.Message.CreatedAt:type_name -> google.protobuf.Timestamp
	0,  // 1: chat.GetHistoryArg.DI:type_name -> chat.DialogIdentifier
	2,  // 2: chat.GetHistoryArg.FP:type_name -> chat.FilterParams
	9,  // 3: chat.Dialog.CreatedAt:type_name -> google.protobuf.Timestamp
	3,  // 4: chat.Messages.m:type_name -> chat.Message
	6,  // 5: chat.Dialogs.d:type_name -> chat.Dialog
	5,  // 6: chat.Chat.GetHistory:input_type -> chat.GetHistoryArg
	3,  // 7: chat.Chat.Create:input_type -> chat.Message
	0,  // 8: chat.Chat.Clear:input_type -> chat.DialogIdentifier
	1,  // 9: chat.Chat.GetDialogs:input_type -> chat.UserIdentifier
	7,  // 10: chat.Chat.GetHistory:output_type -> chat.Messages
	4,  // 11: chat.Chat.Create:output_type -> chat.Nothing
	4,  // 12: chat.Chat.Clear:output_type -> chat.Nothing
	8,  // 13: chat.Chat.GetDialogs:output_type -> chat.Dialogs
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DialogIdentifier); i {
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
		file_chat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIdentifier); i {
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
		file_chat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilterParams); i {
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
		file_chat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_chat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Nothing); i {
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
		file_chat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetHistoryArg); i {
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
		file_chat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dialog); i {
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
		file_chat_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Messages); i {
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
		file_chat_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dialogs); i {
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
			RawDescriptor: file_chat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_rawDesc = nil
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}
