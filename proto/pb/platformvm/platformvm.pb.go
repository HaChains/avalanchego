// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: platformvm/platformvm.proto

package platformvm

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

type SubnetValidatorRegistrationJustification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Preimage:
	//
	//	*SubnetValidatorRegistrationJustification_ConvertSubnetTxData
	//	*SubnetValidatorRegistrationJustification_RegisterSubnetValidatorMessage
	Preimage isSubnetValidatorRegistrationJustification_Preimage `protobuf_oneof:"preimage"`
	Filter   []byte                                              `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *SubnetValidatorRegistrationJustification) Reset() {
	*x = SubnetValidatorRegistrationJustification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platformvm_platformvm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubnetValidatorRegistrationJustification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubnetValidatorRegistrationJustification) ProtoMessage() {}

func (x *SubnetValidatorRegistrationJustification) ProtoReflect() protoreflect.Message {
	mi := &file_platformvm_platformvm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubnetValidatorRegistrationJustification.ProtoReflect.Descriptor instead.
func (*SubnetValidatorRegistrationJustification) Descriptor() ([]byte, []int) {
	return file_platformvm_platformvm_proto_rawDescGZIP(), []int{0}
}

func (m *SubnetValidatorRegistrationJustification) GetPreimage() isSubnetValidatorRegistrationJustification_Preimage {
	if m != nil {
		return m.Preimage
	}
	return nil
}

func (x *SubnetValidatorRegistrationJustification) GetConvertSubnetTxData() *SubnetIDIndex {
	if x, ok := x.GetPreimage().(*SubnetValidatorRegistrationJustification_ConvertSubnetTxData); ok {
		return x.ConvertSubnetTxData
	}
	return nil
}

func (x *SubnetValidatorRegistrationJustification) GetRegisterSubnetValidatorMessage() []byte {
	if x, ok := x.GetPreimage().(*SubnetValidatorRegistrationJustification_RegisterSubnetValidatorMessage); ok {
		return x.RegisterSubnetValidatorMessage
	}
	return nil
}

func (x *SubnetValidatorRegistrationJustification) GetFilter() []byte {
	if x != nil {
		return x.Filter
	}
	return nil
}

type isSubnetValidatorRegistrationJustification_Preimage interface {
	isSubnetValidatorRegistrationJustification_Preimage()
}

type SubnetValidatorRegistrationJustification_ConvertSubnetTxData struct {
	// Validator was added to the Subnet during the ConvertSubnetTx.
	ConvertSubnetTxData *SubnetIDIndex `protobuf:"bytes,1,opt,name=convert_subnet_tx_data,json=convertSubnetTxData,proto3,oneof"`
}

type SubnetValidatorRegistrationJustification_RegisterSubnetValidatorMessage struct {
	// Validator was registered to the Subnet after the ConvertSubnetTx.
	// The SubnetValidator is being removed from the Subnet
	RegisterSubnetValidatorMessage []byte `protobuf:"bytes,2,opt,name=register_subnet_validator_message,json=registerSubnetValidatorMessage,proto3,oneof"`
}

func (*SubnetValidatorRegistrationJustification_ConvertSubnetTxData) isSubnetValidatorRegistrationJustification_Preimage() {
}

func (*SubnetValidatorRegistrationJustification_RegisterSubnetValidatorMessage) isSubnetValidatorRegistrationJustification_Preimage() {
}

type SubnetIDIndex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubnetId []byte `protobuf:"bytes,1,opt,name=subnet_id,json=subnetId,proto3" json:"subnet_id,omitempty"`
	Index    uint32 `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *SubnetIDIndex) Reset() {
	*x = SubnetIDIndex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_platformvm_platformvm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubnetIDIndex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubnetIDIndex) ProtoMessage() {}

func (x *SubnetIDIndex) ProtoReflect() protoreflect.Message {
	mi := &file_platformvm_platformvm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubnetIDIndex.ProtoReflect.Descriptor instead.
func (*SubnetIDIndex) Descriptor() ([]byte, []int) {
	return file_platformvm_platformvm_proto_rawDescGZIP(), []int{1}
}

func (x *SubnetIDIndex) GetSubnetId() []byte {
	if x != nil {
		return x.SubnetId
	}
	return nil
}

func (x *SubnetIDIndex) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

var File_platformvm_platformvm_proto protoreflect.FileDescriptor

var file_platformvm_platformvm_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x76, 0x6d, 0x2f, 0x70, 0x6c, 0x61,
	0x74, 0x66, 0x6f, 0x72, 0x6d, 0x76, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x76, 0x6d, 0x22, 0xed, 0x01, 0x0a, 0x28, 0x53, 0x75,
	0x62, 0x6e, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x75, 0x73, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x50, 0x0a, 0x16, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72,
	0x74, 0x5f, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x74, 0x78, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x76, 0x6d, 0x2e, 0x53, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x44, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x48, 0x00, 0x52, 0x13, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x53, 0x75, 0x62, 0x6e,
	0x65, 0x74, 0x54, 0x78, 0x44, 0x61, 0x74, 0x61, 0x12, 0x4b, 0x0a, 0x21, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x1e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53,
	0x75, 0x62, 0x6e, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x42, 0x0a, 0x0a,
	0x08, 0x70, 0x72, 0x65, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x42, 0x0a, 0x0d, 0x53, 0x75, 0x62,
	0x6e, 0x65, 0x74, 0x49, 0x44, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x75,
	0x62, 0x6e, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x73,
	0x75, 0x62, 0x6e, 0x65, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x35, 0x5a,
	0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x76, 0x61, 0x2d,
	0x6c, 0x61, 0x62, 0x73, 0x2f, 0x61, 0x76, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x67, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x62, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x76, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_platformvm_platformvm_proto_rawDescOnce sync.Once
	file_platformvm_platformvm_proto_rawDescData = file_platformvm_platformvm_proto_rawDesc
)

func file_platformvm_platformvm_proto_rawDescGZIP() []byte {
	file_platformvm_platformvm_proto_rawDescOnce.Do(func() {
		file_platformvm_platformvm_proto_rawDescData = protoimpl.X.CompressGZIP(file_platformvm_platformvm_proto_rawDescData)
	})
	return file_platformvm_platformvm_proto_rawDescData
}

var file_platformvm_platformvm_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_platformvm_platformvm_proto_goTypes = []interface{}{
	(*SubnetValidatorRegistrationJustification)(nil), // 0: platformvm.SubnetValidatorRegistrationJustification
	(*SubnetIDIndex)(nil),                            // 1: platformvm.SubnetIDIndex
}
var file_platformvm_platformvm_proto_depIdxs = []int32{
	1, // 0: platformvm.SubnetValidatorRegistrationJustification.convert_subnet_tx_data:type_name -> platformvm.SubnetIDIndex
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_platformvm_platformvm_proto_init() }
func file_platformvm_platformvm_proto_init() {
	if File_platformvm_platformvm_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_platformvm_platformvm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubnetValidatorRegistrationJustification); i {
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
		file_platformvm_platformvm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubnetIDIndex); i {
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
	file_platformvm_platformvm_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*SubnetValidatorRegistrationJustification_ConvertSubnetTxData)(nil),
		(*SubnetValidatorRegistrationJustification_RegisterSubnetValidatorMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_platformvm_platformvm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_platformvm_platformvm_proto_goTypes,
		DependencyIndexes: file_platformvm_platformvm_proto_depIdxs,
		MessageInfos:      file_platformvm_platformvm_proto_msgTypes,
	}.Build()
	File_platformvm_platformvm_proto = out.File
	file_platformvm_platformvm_proto_rawDesc = nil
	file_platformvm_platformvm_proto_goTypes = nil
	file_platformvm_platformvm_proto_depIdxs = nil
}
