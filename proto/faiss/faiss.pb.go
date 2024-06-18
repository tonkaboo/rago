// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: faiss.proto

package faiss

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

// SimilarityRequest represents a request containing a vector embedding.
type SimilarityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Embedding []float32 `protobuf:"fixed32,1,rep,packed,name=embedding,proto3" json:"embedding,omitempty"` // A list of float32 numbers representing the embedding vector.
}

func (x *SimilarityRequest) Reset() {
	*x = SimilarityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_faiss_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimilarityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimilarityRequest) ProtoMessage() {}

func (x *SimilarityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_faiss_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimilarityRequest.ProtoReflect.Descriptor instead.
func (*SimilarityRequest) Descriptor() ([]byte, []int) {
	return file_faiss_proto_rawDescGZIP(), []int{0}
}

func (x *SimilarityRequest) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

// SimilarityReply represents a response containing similar chunk contents or an error.
type SimilarityReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkContents []string `protobuf:"bytes,1,rep,name=chunk_contents,json=chunkContents,proto3" json:"chunk_contents,omitempty"` // Similar chunk contents returned by the server.
	Error         string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`                                      // Error message in case of failure.
}

func (x *SimilarityReply) Reset() {
	*x = SimilarityReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_faiss_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimilarityReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimilarityReply) ProtoMessage() {}

func (x *SimilarityReply) ProtoReflect() protoreflect.Message {
	mi := &file_faiss_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimilarityReply.ProtoReflect.Descriptor instead.
func (*SimilarityReply) Descriptor() ([]byte, []int) {
	return file_faiss_proto_rawDescGZIP(), []int{1}
}

func (x *SimilarityReply) GetChunkContents() []string {
	if x != nil {
		return x.ChunkContents
	}
	return nil
}

func (x *SimilarityReply) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_faiss_proto protoreflect.FileDescriptor

var file_faiss_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x66, 0x61, 0x69, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x66,
	0x61, 0x69, 0x73, 0x73, 0x22, 0x31, 0x0a, 0x11, 0x53, 0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x69,
	0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x09, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x22, 0x4e, 0x0a, 0x0f, 0x53, 0x69, 0x6d, 0x69, 0x6c,
	0x61, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0d, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x57, 0x0a, 0x0c, 0x46, 0x61, 0x69, 0x73, 0x73,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x53,
	0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x12, 0x18, 0x2e, 0x66,
	0x61, 0x69, 0x73, 0x73, 0x2e, 0x53, 0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x69, 0x74, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x66, 0x61, 0x69, 0x73, 0x73, 0x2e, 0x53,
	0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x42, 0x12, 0x5a, 0x10, 0x67, 0x65, 0x6d, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66,
	0x61, 0x69, 0x73, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_faiss_proto_rawDescOnce sync.Once
	file_faiss_proto_rawDescData = file_faiss_proto_rawDesc
)

func file_faiss_proto_rawDescGZIP() []byte {
	file_faiss_proto_rawDescOnce.Do(func() {
		file_faiss_proto_rawDescData = protoimpl.X.CompressGZIP(file_faiss_proto_rawDescData)
	})
	return file_faiss_proto_rawDescData
}

var file_faiss_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_faiss_proto_goTypes = []interface{}{
	(*SimilarityRequest)(nil), // 0: faiss.SimilarityRequest
	(*SimilarityReply)(nil),   // 1: faiss.SimilarityReply
}
var file_faiss_proto_depIdxs = []int32{
	0, // 0: faiss.FaissService.FindSimilarChunks:input_type -> faiss.SimilarityRequest
	1, // 1: faiss.FaissService.FindSimilarChunks:output_type -> faiss.SimilarityReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_faiss_proto_init() }
func file_faiss_proto_init() {
	if File_faiss_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_faiss_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SimilarityRequest); i {
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
		file_faiss_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SimilarityReply); i {
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
			RawDescriptor: file_faiss_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_faiss_proto_goTypes,
		DependencyIndexes: file_faiss_proto_depIdxs,
		MessageInfos:      file_faiss_proto_msgTypes,
	}.Build()
	File_faiss_proto = out.File
	file_faiss_proto_rawDesc = nil
	file_faiss_proto_goTypes = nil
	file_faiss_proto_depIdxs = nil
}
