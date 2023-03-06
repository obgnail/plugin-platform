// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: plugin.proto

package protocol

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

type NotificationMessage_EventType int32

const (
	NotificationMessage_TimeStamp NotificationMessage_EventType = 0
)

// Enum value maps for NotificationMessage_EventType.
var (
	NotificationMessage_EventType_name = map[int32]string{
		0: "TimeStamp",
	}
	NotificationMessage_EventType_value = map[string]int32{
		"TimeStamp": 0,
	}
)

func (x NotificationMessage_EventType) Enum() *NotificationMessage_EventType {
	p := new(NotificationMessage_EventType)
	*p = x
	return p
}

func (x NotificationMessage_EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NotificationMessage_EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_plugin_proto_enumTypes[0].Descriptor()
}

func (NotificationMessage_EventType) Type() protoreflect.EnumType {
	return &file_plugin_proto_enumTypes[0]
}

func (x NotificationMessage_EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NotificationMessage_EventType.Descriptor instead.
func (NotificationMessage_EventType) EnumDescriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{1, 0}
}

// 插件需要相应的消息请求及其对应的应答
// 插件提供的接口实现
// 插件提供的配置处理实现
type ConfigurationMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigChangeRequest *ConfigurationMessage_ConfigurationChangeMessage `protobuf:"bytes,1,opt,name=ConfigChangeRequest,proto3" json:"ConfigChangeRequest,omitempty"`
	ConfigChaneResponse *ErrorMessage                                    `protobuf:"bytes,2,opt,name=ConfigChaneResponse,proto3" json:"ConfigChaneResponse,omitempty"`
}

func (x *ConfigurationMessage) Reset() {
	*x = ConfigurationMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigurationMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigurationMessage) ProtoMessage() {}

func (x *ConfigurationMessage) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigurationMessage.ProtoReflect.Descriptor instead.
func (*ConfigurationMessage) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigurationMessage) GetConfigChangeRequest() *ConfigurationMessage_ConfigurationChangeMessage {
	if x != nil {
		return x.ConfigChangeRequest
	}
	return nil
}

func (x *ConfigurationMessage) GetConfigChaneResponse() *ErrorMessage {
	if x != nil {
		return x.ConfigChaneResponse
	}
	return nil
}

// 事件
type NotificationMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      NotificationMessage_EventType `protobuf:"varint,1,opt,name=Type,proto3,enum=protocol.NotificationMessage_EventType" json:"Type,omitempty"`
	Timestamp int64                         `protobuf:"varint,2,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Data      string                        `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *NotificationMessage) Reset() {
	*x = NotificationMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationMessage) ProtoMessage() {}

func (x *NotificationMessage) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationMessage.ProtoReflect.Descriptor instead.
func (*NotificationMessage) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *NotificationMessage) GetType() NotificationMessage_EventType {
	if x != nil {
		return x.Type
	}
	return NotificationMessage_TimeStamp
}

func (x *NotificationMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *NotificationMessage) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type PluginMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 添加消息路由数据
	Target *PluginInstanceDescriptor `protobuf:"bytes,1,opt,name=Target,proto3" json:"Target,omitempty"`
	Host   *HostDescriptor           `protobuf:"bytes,2,opt,name=Host,proto3" json:"Host,omitempty"`
	// 插件实现的各种http方法，包括对内（前端）和对外（独立的http服务）
	Http *HttpContextMessage `protobuf:"bytes,21,opt,name=Http,proto3" json:"Http,omitempty"`
	// 插件配置变动通知
	Config *ConfigurationMessage `protobuf:"bytes,22,opt,name=Config,proto3" json:"Config,omitempty"`
	// 插件接收的通知消息
	Notification *NotificationMessage `protobuf:"bytes,23,opt,name=Notification,proto3" json:"Notification,omitempty"`
}

func (x *PluginMessage) Reset() {
	*x = PluginMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginMessage) ProtoMessage() {}

func (x *PluginMessage) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginMessage.ProtoReflect.Descriptor instead.
func (*PluginMessage) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *PluginMessage) GetTarget() *PluginInstanceDescriptor {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *PluginMessage) GetHost() *HostDescriptor {
	if x != nil {
		return x.Host
	}
	return nil
}

func (x *PluginMessage) GetHttp() *HttpContextMessage {
	if x != nil {
		return x.Http
	}
	return nil
}

func (x *PluginMessage) GetConfig() *ConfigurationMessage {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *PluginMessage) GetNotification() *NotificationMessage {
	if x != nil {
		return x.Notification
	}
	return nil
}

type ConfigurationMessage_ConfigurationChangeMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewConfig  string `protobuf:"bytes,1,opt,name=NewConfig,proto3" json:"NewConfig,omitempty"`
	OrigConfig string `protobuf:"bytes,2,opt,name=OrigConfig,proto3" json:"OrigConfig,omitempty"`
}

func (x *ConfigurationMessage_ConfigurationChangeMessage) Reset() {
	*x = ConfigurationMessage_ConfigurationChangeMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigurationMessage_ConfigurationChangeMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigurationMessage_ConfigurationChangeMessage) ProtoMessage() {}

func (x *ConfigurationMessage_ConfigurationChangeMessage) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigurationMessage_ConfigurationChangeMessage.ProtoReflect.Descriptor instead.
func (*ConfigurationMessage_ConfigurationChangeMessage) Descriptor() ([]byte, []int) {
	return file_plugin_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ConfigurationMessage_ConfigurationChangeMessage) GetNewConfig() string {
	if x != nil {
		return x.NewConfig
	}
	return ""
}

func (x *ConfigurationMessage_ConfigurationChangeMessage) GetOrigConfig() string {
	if x != nil {
		return x.OrigConfig
	}
	return ""
}

var File_plugin_proto protoreflect.FileDescriptor

var file_plugin_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x02, 0x0a, 0x14, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x6b, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x13, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x48, 0x0a, 0x13,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x43, 0x68, 0x61, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x13, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x43, 0x68, 0x61, 0x6e, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x5a, 0x0a, 0x1a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x4e, 0x65, 0x77, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4e, 0x65, 0x77, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x4f, 0x72, 0x69, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x4f, 0x72, 0x69, 0x67, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x22, 0xa0, 0x01, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3b, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x1a, 0x0a, 0x09, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x6d, 0x70, 0x10, 0x00, 0x22, 0xa6, 0x02, 0x0a, 0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3a, 0x0a, 0x06, 0x54, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x06, 0x54, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x12, 0x2c, 0x0a, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x48, 0x6f, 0x73,
	0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x48, 0x6f, 0x73,
	0x74, 0x12, 0x30, 0x0a, 0x04, 0x48, 0x74, 0x74, 0x70, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x48,
	0x74, 0x74, 0x70, 0x12, 0x36, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x16, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x41, 0x0a, 0x0c, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x17, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0d,
	0x5a, 0x0b, 0x2e, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugin_proto_rawDescOnce sync.Once
	file_plugin_proto_rawDescData = file_plugin_proto_rawDesc
)

func file_plugin_proto_rawDescGZIP() []byte {
	file_plugin_proto_rawDescOnce.Do(func() {
		file_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_proto_rawDescData)
	})
	return file_plugin_proto_rawDescData
}

var file_plugin_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_plugin_proto_goTypes = []interface{}{
	(NotificationMessage_EventType)(0),                      // 0: protocol.NotificationMessage.EventType
	(*ConfigurationMessage)(nil),                            // 1: protocol.ConfigurationMessage
	(*NotificationMessage)(nil),                             // 2: protocol.NotificationMessage
	(*PluginMessage)(nil),                                   // 3: protocol.PluginMessage
	(*ConfigurationMessage_ConfigurationChangeMessage)(nil), // 4: protocol.ConfigurationMessage.ConfigurationChangeMessage
	(*ErrorMessage)(nil),                                    // 5: protocol.ErrorMessage
	(*PluginInstanceDescriptor)(nil),                        // 6: protocol.PluginInstanceDescriptor
	(*HostDescriptor)(nil),                                  // 7: protocol.HostDescriptor
	(*HttpContextMessage)(nil),                              // 8: protocol.HttpContextMessage
}
var file_plugin_proto_depIdxs = []int32{
	4, // 0: protocol.ConfigurationMessage.ConfigChangeRequest:type_name -> protocol.ConfigurationMessage.ConfigurationChangeMessage
	5, // 1: protocol.ConfigurationMessage.ConfigChaneResponse:type_name -> protocol.ErrorMessage
	0, // 2: protocol.NotificationMessage.Type:type_name -> protocol.NotificationMessage.EventType
	6, // 3: protocol.PluginMessage.Target:type_name -> protocol.PluginInstanceDescriptor
	7, // 4: protocol.PluginMessage.Host:type_name -> protocol.HostDescriptor
	8, // 5: protocol.PluginMessage.Http:type_name -> protocol.HttpContextMessage
	1, // 6: protocol.PluginMessage.Config:type_name -> protocol.ConfigurationMessage
	2, // 7: protocol.PluginMessage.Notification:type_name -> protocol.NotificationMessage
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_plugin_proto_init() }
func file_plugin_proto_init() {
	if File_plugin_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigurationMessage); i {
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
		file_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationMessage); i {
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
		file_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginMessage); i {
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
		file_plugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigurationMessage_ConfigurationChangeMessage); i {
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
			RawDescriptor: file_plugin_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_plugin_proto_goTypes,
		DependencyIndexes: file_plugin_proto_depIdxs,
		EnumInfos:         file_plugin_proto_enumTypes,
		MessageInfos:      file_plugin_proto_msgTypes,
	}.Build()
	File_plugin_proto = out.File
	file_plugin_proto_rawDesc = nil
	file_plugin_proto_goTypes = nil
	file_plugin_proto_depIdxs = nil
}