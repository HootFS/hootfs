// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/head.proto

package protos

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ResponseStatus_StatusCode int32

const (
	ResponseStatus_SUCCESS ResponseStatus_StatusCode = 0
	ResponseStatus_WARNING ResponseStatus_StatusCode = 1
	ResponseStatus_FAILURE ResponseStatus_StatusCode = 2
)

var ResponseStatus_StatusCode_name = map[int32]string{
	0: "SUCCESS",
	1: "WARNING",
	2: "FAILURE",
}

var ResponseStatus_StatusCode_value = map[string]int32{
	"SUCCESS": 0,
	"WARNING": 1,
	"FAILURE": 2,
}

func (x ResponseStatus_StatusCode) String() string {
	return proto.EnumName(ResponseStatus_StatusCode_name, int32(x))
}

func (ResponseStatus_StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{0, 0}
}

type ObjectInfo_ObjectType int32

const (
	ObjectInfo_FILE      ObjectInfo_ObjectType = 0
	ObjectInfo_DIRECTORY ObjectInfo_ObjectType = 1
)

var ObjectInfo_ObjectType_name = map[int32]string{
	0: "FILE",
	1: "DIRECTORY",
}

var ObjectInfo_ObjectType_value = map[string]int32{
	"FILE":      0,
	"DIRECTORY": 1,
}

func (x ObjectInfo_ObjectType) String() string {
	return proto.EnumName(ObjectInfo_ObjectType_name, int32(x))
}

func (ObjectInfo_ObjectType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{3, 0}
}

type ResponseStatus struct {
	Code                 ResponseStatus_StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=hootfs.head.ResponseStatus_StatusCode" json:"code,omitempty"`
	ResMessage           string                    `protobuf:"bytes,2,opt,name=res_message,json=resMessage,proto3" json:"res_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ResponseStatus) Reset()         { *m = ResponseStatus{} }
func (m *ResponseStatus) String() string { return proto.CompactTextString(m) }
func (*ResponseStatus) ProtoMessage()    {}
func (*ResponseStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{0}
}

func (m *ResponseStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseStatus.Unmarshal(m, b)
}
func (m *ResponseStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseStatus.Marshal(b, m, deterministic)
}
func (m *ResponseStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseStatus.Merge(m, src)
}
func (m *ResponseStatus) XXX_Size() int {
	return xxx_messageInfo_ResponseStatus.Size(m)
}
func (m *ResponseStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseStatus proto.InternalMessageInfo

func (m *ResponseStatus) GetCode() ResponseStatus_StatusCode {
	if m != nil {
		return m.Code
	}
	return ResponseStatus_SUCCESS
}

func (m *ResponseStatus) GetResMessage() string {
	if m != nil {
		return m.ResMessage
	}
	return ""
}

// UUID message structure
type UUID struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UUID) Reset()         { *m = UUID{} }
func (m *UUID) String() string { return proto.CompactTextString(m) }
func (*UUID) ProtoMessage()    {}
func (*UUID) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{1}
}

func (m *UUID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UUID.Unmarshal(m, b)
}
func (m *UUID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UUID.Marshal(b, m, deterministic)
}
func (m *UUID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UUID.Merge(m, src)
}
func (m *UUID) XXX_Size() int {
	return xxx_messageInfo_UUID.Size(m)
}
func (m *UUID) XXX_DiscardUnknown() {
	xxx_messageInfo_UUID.DiscardUnknown(m)
}

var xxx_messageInfo_UUID proto.InternalMessageInfo

func (m *UUID) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type GetDirectoryContentsRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	// The directory to get the contents of.
	// If not provided, HOME will be used.
	DirId                *UUID    `protobuf:"bytes,2,opt,name=dir_id,json=dirId,proto3" json:"dir_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetDirectoryContentsRequest) Reset()         { *m = GetDirectoryContentsRequest{} }
func (m *GetDirectoryContentsRequest) String() string { return proto.CompactTextString(m) }
func (*GetDirectoryContentsRequest) ProtoMessage()    {}
func (*GetDirectoryContentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{2}
}

func (m *GetDirectoryContentsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDirectoryContentsRequest.Unmarshal(m, b)
}
func (m *GetDirectoryContentsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDirectoryContentsRequest.Marshal(b, m, deterministic)
}
func (m *GetDirectoryContentsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDirectoryContentsRequest.Merge(m, src)
}
func (m *GetDirectoryContentsRequest) XXX_Size() int {
	return xxx_messageInfo_GetDirectoryContentsRequest.Size(m)
}
func (m *GetDirectoryContentsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDirectoryContentsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetDirectoryContentsRequest proto.InternalMessageInfo

func (m *GetDirectoryContentsRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *GetDirectoryContentsRequest) GetDirId() *UUID {
	if m != nil {
		return m.DirId
	}
	return nil
}

type ObjectInfo struct {
	// ID and type of the object.
	ObjectId             *UUID                 `protobuf:"bytes,1,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	ObjectType           ObjectInfo_ObjectType `protobuf:"varint,2,opt,name=object_type,json=objectType,proto3,enum=hootfs.head.ObjectInfo_ObjectType" json:"object_type,omitempty"`
	ObjectName           string                `protobuf:"bytes,3,opt,name=object_name,json=objectName,proto3" json:"object_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ObjectInfo) Reset()         { *m = ObjectInfo{} }
func (m *ObjectInfo) String() string { return proto.CompactTextString(m) }
func (*ObjectInfo) ProtoMessage()    {}
func (*ObjectInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{3}
}

func (m *ObjectInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ObjectInfo.Unmarshal(m, b)
}
func (m *ObjectInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ObjectInfo.Marshal(b, m, deterministic)
}
func (m *ObjectInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObjectInfo.Merge(m, src)
}
func (m *ObjectInfo) XXX_Size() int {
	return xxx_messageInfo_ObjectInfo.Size(m)
}
func (m *ObjectInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ObjectInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ObjectInfo proto.InternalMessageInfo

func (m *ObjectInfo) GetObjectId() *UUID {
	if m != nil {
		return m.ObjectId
	}
	return nil
}

func (m *ObjectInfo) GetObjectType() ObjectInfo_ObjectType {
	if m != nil {
		return m.ObjectType
	}
	return ObjectInfo_FILE
}

func (m *ObjectInfo) GetObjectName() string {
	if m != nil {
		return m.ObjectName
	}
	return ""
}

type GetDirectoryContentsResponse struct {
	Status               *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Objects              []*ObjectInfo   `protobuf:"bytes,2,rep,name=objects,proto3" json:"objects,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetDirectoryContentsResponse) Reset()         { *m = GetDirectoryContentsResponse{} }
func (m *GetDirectoryContentsResponse) String() string { return proto.CompactTextString(m) }
func (*GetDirectoryContentsResponse) ProtoMessage()    {}
func (*GetDirectoryContentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{4}
}

func (m *GetDirectoryContentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetDirectoryContentsResponse.Unmarshal(m, b)
}
func (m *GetDirectoryContentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetDirectoryContentsResponse.Marshal(b, m, deterministic)
}
func (m *GetDirectoryContentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetDirectoryContentsResponse.Merge(m, src)
}
func (m *GetDirectoryContentsResponse) XXX_Size() int {
	return xxx_messageInfo_GetDirectoryContentsResponse.Size(m)
}
func (m *GetDirectoryContentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetDirectoryContentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetDirectoryContentsResponse proto.InternalMessageInfo

func (m *GetDirectoryContentsResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetDirectoryContentsResponse) GetObjects() []*ObjectInfo {
	if m != nil {
		return m.Objects
	}
	return nil
}

type MakeDirectoryRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	// This is the ID of the directory
	// the new directory will be created into.
	// If this is not provided, HOME will be used.
	DirId *UUID `protobuf:"bytes,2,opt,name=dir_id,json=dirId,proto3" json:"dir_id,omitempty"`
	// This is the name of the new directory.
	DirName              string   `protobuf:"bytes,3,opt,name=dir_name,json=dirName,proto3" json:"dir_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MakeDirectoryRequest) Reset()         { *m = MakeDirectoryRequest{} }
func (m *MakeDirectoryRequest) String() string { return proto.CompactTextString(m) }
func (*MakeDirectoryRequest) ProtoMessage()    {}
func (*MakeDirectoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{5}
}

func (m *MakeDirectoryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MakeDirectoryRequest.Unmarshal(m, b)
}
func (m *MakeDirectoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MakeDirectoryRequest.Marshal(b, m, deterministic)
}
func (m *MakeDirectoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MakeDirectoryRequest.Merge(m, src)
}
func (m *MakeDirectoryRequest) XXX_Size() int {
	return xxx_messageInfo_MakeDirectoryRequest.Size(m)
}
func (m *MakeDirectoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MakeDirectoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MakeDirectoryRequest proto.InternalMessageInfo

func (m *MakeDirectoryRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *MakeDirectoryRequest) GetDirId() *UUID {
	if m != nil {
		return m.DirId
	}
	return nil
}

func (m *MakeDirectoryRequest) GetDirName() string {
	if m != nil {
		return m.DirName
	}
	return ""
}

type MakeDirectoryResponse struct {
	Status *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// The ID of the directory created.
	DirId                *UUID    `protobuf:"bytes,2,opt,name=dir_id,json=dirId,proto3" json:"dir_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MakeDirectoryResponse) Reset()         { *m = MakeDirectoryResponse{} }
func (m *MakeDirectoryResponse) String() string { return proto.CompactTextString(m) }
func (*MakeDirectoryResponse) ProtoMessage()    {}
func (*MakeDirectoryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{6}
}

func (m *MakeDirectoryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MakeDirectoryResponse.Unmarshal(m, b)
}
func (m *MakeDirectoryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MakeDirectoryResponse.Marshal(b, m, deterministic)
}
func (m *MakeDirectoryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MakeDirectoryResponse.Merge(m, src)
}
func (m *MakeDirectoryResponse) XXX_Size() int {
	return xxx_messageInfo_MakeDirectoryResponse.Size(m)
}
func (m *MakeDirectoryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MakeDirectoryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MakeDirectoryResponse proto.InternalMessageInfo

func (m *MakeDirectoryResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *MakeDirectoryResponse) GetDirId() *UUID {
	if m != nil {
		return m.DirId
	}
	return nil
}

type AddNewFileRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	// ID of directory to place folder in.
	DirId *UUID `protobuf:"bytes,2,opt,name=dir_id,json=dirId,proto3" json:"dir_id,omitempty"`
	// Name of file to create.
	FileName string `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// Contents of file being added...
	// If not provided, an empty file
	// will be created.
	Contents             []byte   `protobuf:"bytes,4,opt,name=contents,proto3" json:"contents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddNewFileRequest) Reset()         { *m = AddNewFileRequest{} }
func (m *AddNewFileRequest) String() string { return proto.CompactTextString(m) }
func (*AddNewFileRequest) ProtoMessage()    {}
func (*AddNewFileRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{7}
}

func (m *AddNewFileRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddNewFileRequest.Unmarshal(m, b)
}
func (m *AddNewFileRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddNewFileRequest.Marshal(b, m, deterministic)
}
func (m *AddNewFileRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddNewFileRequest.Merge(m, src)
}
func (m *AddNewFileRequest) XXX_Size() int {
	return xxx_messageInfo_AddNewFileRequest.Size(m)
}
func (m *AddNewFileRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddNewFileRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddNewFileRequest proto.InternalMessageInfo

func (m *AddNewFileRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *AddNewFileRequest) GetDirId() *UUID {
	if m != nil {
		return m.DirId
	}
	return nil
}

func (m *AddNewFileRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *AddNewFileRequest) GetContents() []byte {
	if m != nil {
		return m.Contents
	}
	return nil
}

type AddNewFileResponse struct {
	Status *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// The ID of the file created.
	FileId               *UUID    `protobuf:"bytes,2,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddNewFileResponse) Reset()         { *m = AddNewFileResponse{} }
func (m *AddNewFileResponse) String() string { return proto.CompactTextString(m) }
func (*AddNewFileResponse) ProtoMessage()    {}
func (*AddNewFileResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{8}
}

func (m *AddNewFileResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddNewFileResponse.Unmarshal(m, b)
}
func (m *AddNewFileResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddNewFileResponse.Marshal(b, m, deterministic)
}
func (m *AddNewFileResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddNewFileResponse.Merge(m, src)
}
func (m *AddNewFileResponse) XXX_Size() int {
	return xxx_messageInfo_AddNewFileResponse.Size(m)
}
func (m *AddNewFileResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddNewFileResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddNewFileResponse proto.InternalMessageInfo

func (m *AddNewFileResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *AddNewFileResponse) GetFileId() *UUID {
	if m != nil {
		return m.FileId
	}
	return nil
}

type UpdateFileContentsRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	FileId     *UUID  `protobuf:"bytes,2,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	// The given file's contents will be
	// set to contents.
	Contents             []byte   `protobuf:"bytes,3,opt,name=contents,proto3" json:"contents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateFileContentsRequest) Reset()         { *m = UpdateFileContentsRequest{} }
func (m *UpdateFileContentsRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateFileContentsRequest) ProtoMessage()    {}
func (*UpdateFileContentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{9}
}

func (m *UpdateFileContentsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateFileContentsRequest.Unmarshal(m, b)
}
func (m *UpdateFileContentsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateFileContentsRequest.Marshal(b, m, deterministic)
}
func (m *UpdateFileContentsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateFileContentsRequest.Merge(m, src)
}
func (m *UpdateFileContentsRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateFileContentsRequest.Size(m)
}
func (m *UpdateFileContentsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateFileContentsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateFileContentsRequest proto.InternalMessageInfo

func (m *UpdateFileContentsRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *UpdateFileContentsRequest) GetFileId() *UUID {
	if m != nil {
		return m.FileId
	}
	return nil
}

func (m *UpdateFileContentsRequest) GetContents() []byte {
	if m != nil {
		return m.Contents
	}
	return nil
}

type UpdateFileContentsResponse struct {
	Status               *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UpdateFileContentsResponse) Reset()         { *m = UpdateFileContentsResponse{} }
func (m *UpdateFileContentsResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateFileContentsResponse) ProtoMessage()    {}
func (*UpdateFileContentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{10}
}

func (m *UpdateFileContentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateFileContentsResponse.Unmarshal(m, b)
}
func (m *UpdateFileContentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateFileContentsResponse.Marshal(b, m, deterministic)
}
func (m *UpdateFileContentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateFileContentsResponse.Merge(m, src)
}
func (m *UpdateFileContentsResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateFileContentsResponse.Size(m)
}
func (m *UpdateFileContentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateFileContentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateFileContentsResponse proto.InternalMessageInfo

func (m *UpdateFileContentsResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

type GetFileContentsRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	// File to download.
	FileId               *UUID    `protobuf:"bytes,2,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFileContentsRequest) Reset()         { *m = GetFileContentsRequest{} }
func (m *GetFileContentsRequest) String() string { return proto.CompactTextString(m) }
func (*GetFileContentsRequest) ProtoMessage()    {}
func (*GetFileContentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{11}
}

func (m *GetFileContentsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFileContentsRequest.Unmarshal(m, b)
}
func (m *GetFileContentsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFileContentsRequest.Marshal(b, m, deterministic)
}
func (m *GetFileContentsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFileContentsRequest.Merge(m, src)
}
func (m *GetFileContentsRequest) XXX_Size() int {
	return xxx_messageInfo_GetFileContentsRequest.Size(m)
}
func (m *GetFileContentsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFileContentsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetFileContentsRequest proto.InternalMessageInfo

func (m *GetFileContentsRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *GetFileContentsRequest) GetFileId() *UUID {
	if m != nil {
		return m.FileId
	}
	return nil
}

type GetFileContentsResponse struct {
	Status *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// Contents of file. (If found)
	Contents             []byte   `protobuf:"bytes,2,opt,name=contents,proto3" json:"contents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFileContentsResponse) Reset()         { *m = GetFileContentsResponse{} }
func (m *GetFileContentsResponse) String() string { return proto.CompactTextString(m) }
func (*GetFileContentsResponse) ProtoMessage()    {}
func (*GetFileContentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{12}
}

func (m *GetFileContentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFileContentsResponse.Unmarshal(m, b)
}
func (m *GetFileContentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFileContentsResponse.Marshal(b, m, deterministic)
}
func (m *GetFileContentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFileContentsResponse.Merge(m, src)
}
func (m *GetFileContentsResponse) XXX_Size() int {
	return xxx_messageInfo_GetFileContentsResponse.Size(m)
}
func (m *GetFileContentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFileContentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetFileContentsResponse proto.InternalMessageInfo

func (m *GetFileContentsResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *GetFileContentsResponse) GetContents() []byte {
	if m != nil {
		return m.Contents
	}
	return nil
}

type MoveObjectRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	// Object to move.
	ObjectId *UUID `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	// Directory to move to.
	DirId *UUID `protobuf:"bytes,3,opt,name=dir_id,json=dirId,proto3" json:"dir_id,omitempty"`
	// New name of object.
	NewName              string   `protobuf:"bytes,4,opt,name=new_name,json=newName,proto3" json:"new_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MoveObjectRequest) Reset()         { *m = MoveObjectRequest{} }
func (m *MoveObjectRequest) String() string { return proto.CompactTextString(m) }
func (*MoveObjectRequest) ProtoMessage()    {}
func (*MoveObjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{13}
}

func (m *MoveObjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MoveObjectRequest.Unmarshal(m, b)
}
func (m *MoveObjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MoveObjectRequest.Marshal(b, m, deterministic)
}
func (m *MoveObjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MoveObjectRequest.Merge(m, src)
}
func (m *MoveObjectRequest) XXX_Size() int {
	return xxx_messageInfo_MoveObjectRequest.Size(m)
}
func (m *MoveObjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MoveObjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MoveObjectRequest proto.InternalMessageInfo

func (m *MoveObjectRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *MoveObjectRequest) GetObjectId() *UUID {
	if m != nil {
		return m.ObjectId
	}
	return nil
}

func (m *MoveObjectRequest) GetDirId() *UUID {
	if m != nil {
		return m.DirId
	}
	return nil
}

func (m *MoveObjectRequest) GetNewName() string {
	if m != nil {
		return m.NewName
	}
	return ""
}

type MoveObjectResponse struct {
	Status *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// This will be the new object ID
	// Assigned to the moved object.
	// NOTE : maybe object_id won't change on move..
	// We can decide this later.
	ObjectId             *UUID    `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MoveObjectResponse) Reset()         { *m = MoveObjectResponse{} }
func (m *MoveObjectResponse) String() string { return proto.CompactTextString(m) }
func (*MoveObjectResponse) ProtoMessage()    {}
func (*MoveObjectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{14}
}

func (m *MoveObjectResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MoveObjectResponse.Unmarshal(m, b)
}
func (m *MoveObjectResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MoveObjectResponse.Marshal(b, m, deterministic)
}
func (m *MoveObjectResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MoveObjectResponse.Merge(m, src)
}
func (m *MoveObjectResponse) XXX_Size() int {
	return xxx_messageInfo_MoveObjectResponse.Size(m)
}
func (m *MoveObjectResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MoveObjectResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MoveObjectResponse proto.InternalMessageInfo

func (m *MoveObjectResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *MoveObjectResponse) GetObjectId() *UUID {
	if m != nil {
		return m.ObjectId
	}
	return nil
}

type RemoveObjectRequest struct {
	SessionKey string `protobuf:"bytes,1,opt,name=session_key,json=sessionKey,proto3" json:"session_key,omitempty"`
	// Object to delete.
	ObjectId             *UUID    `protobuf:"bytes,2,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveObjectRequest) Reset()         { *m = RemoveObjectRequest{} }
func (m *RemoveObjectRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveObjectRequest) ProtoMessage()    {}
func (*RemoveObjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{15}
}

func (m *RemoveObjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveObjectRequest.Unmarshal(m, b)
}
func (m *RemoveObjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveObjectRequest.Marshal(b, m, deterministic)
}
func (m *RemoveObjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveObjectRequest.Merge(m, src)
}
func (m *RemoveObjectRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveObjectRequest.Size(m)
}
func (m *RemoveObjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveObjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveObjectRequest proto.InternalMessageInfo

func (m *RemoveObjectRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func (m *RemoveObjectRequest) GetObjectId() *UUID {
	if m != nil {
		return m.ObjectId
	}
	return nil
}

type RemoveObjectResponse struct {
	Status               *ResponseStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *RemoveObjectResponse) Reset()         { *m = RemoveObjectResponse{} }
func (m *RemoveObjectResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveObjectResponse) ProtoMessage()    {}
func (*RemoveObjectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9514159f75e8540f, []int{16}
}

func (m *RemoveObjectResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveObjectResponse.Unmarshal(m, b)
}
func (m *RemoveObjectResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveObjectResponse.Marshal(b, m, deterministic)
}
func (m *RemoveObjectResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveObjectResponse.Merge(m, src)
}
func (m *RemoveObjectResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveObjectResponse.Size(m)
}
func (m *RemoveObjectResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveObjectResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveObjectResponse proto.InternalMessageInfo

func (m *RemoveObjectResponse) GetStatus() *ResponseStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func init() {
	proto.RegisterEnum("hootfs.head.ResponseStatus_StatusCode", ResponseStatus_StatusCode_name, ResponseStatus_StatusCode_value)
	proto.RegisterEnum("hootfs.head.ObjectInfo_ObjectType", ObjectInfo_ObjectType_name, ObjectInfo_ObjectType_value)
	proto.RegisterType((*ResponseStatus)(nil), "hootfs.head.ResponseStatus")
	proto.RegisterType((*UUID)(nil), "hootfs.head.UUID")
	proto.RegisterType((*GetDirectoryContentsRequest)(nil), "hootfs.head.GetDirectoryContentsRequest")
	proto.RegisterType((*ObjectInfo)(nil), "hootfs.head.ObjectInfo")
	proto.RegisterType((*GetDirectoryContentsResponse)(nil), "hootfs.head.GetDirectoryContentsResponse")
	proto.RegisterType((*MakeDirectoryRequest)(nil), "hootfs.head.MakeDirectoryRequest")
	proto.RegisterType((*MakeDirectoryResponse)(nil), "hootfs.head.MakeDirectoryResponse")
	proto.RegisterType((*AddNewFileRequest)(nil), "hootfs.head.AddNewFileRequest")
	proto.RegisterType((*AddNewFileResponse)(nil), "hootfs.head.AddNewFileResponse")
	proto.RegisterType((*UpdateFileContentsRequest)(nil), "hootfs.head.UpdateFileContentsRequest")
	proto.RegisterType((*UpdateFileContentsResponse)(nil), "hootfs.head.UpdateFileContentsResponse")
	proto.RegisterType((*GetFileContentsRequest)(nil), "hootfs.head.GetFileContentsRequest")
	proto.RegisterType((*GetFileContentsResponse)(nil), "hootfs.head.GetFileContentsResponse")
	proto.RegisterType((*MoveObjectRequest)(nil), "hootfs.head.MoveObjectRequest")
	proto.RegisterType((*MoveObjectResponse)(nil), "hootfs.head.MoveObjectResponse")
	proto.RegisterType((*RemoveObjectRequest)(nil), "hootfs.head.RemoveObjectRequest")
	proto.RegisterType((*RemoveObjectResponse)(nil), "hootfs.head.RemoveObjectResponse")
}

func init() { proto.RegisterFile("protos/head.proto", fileDescriptor_9514159f75e8540f) }

var fileDescriptor_9514159f75e8540f = []byte{
	// 780 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xcf, 0x53, 0x1a, 0x49,
	0x14, 0x76, 0x00, 0xf9, 0xf1, 0x50, 0x17, 0x7a, 0xd9, 0x15, 0xd1, 0x5a, 0xb1, 0xf7, 0x17, 0xbb,
	0x07, 0xac, 0xc5, 0xdb, 0xde, 0x0c, 0x82, 0x99, 0x52, 0xb0, 0x32, 0x48, 0x52, 0x49, 0xa5, 0xca,
	0x1a, 0x99, 0x87, 0x8c, 0xca, 0x34, 0x99, 0x6e, 0xb0, 0xa8, 0xca, 0x31, 0x95, 0x53, 0xee, 0xb9,
	0xe7, 0x96, 0x3f, 0x25, 0xff, 0x55, 0x6a, 0xa6, 0x47, 0x61, 0x86, 0x1f, 0x6a, 0x11, 0x4f, 0x33,
	0xdd, 0x7c, 0xaf, 0xdf, 0xf7, 0x7d, 0xfd, 0xde, 0x1b, 0x20, 0xdd, 0xb3, 0x99, 0x60, 0x7c, 0xb7,
	0x83, 0xba, 0x51, 0x74, 0xdf, 0x49, 0xb2, 0xc3, 0x98, 0x68, 0xf3, 0xa2, 0xb3, 0x45, 0xbf, 0x28,
	0xb0, 0xa6, 0x21, 0xef, 0x31, 0x8b, 0x63, 0x43, 0xe8, 0xa2, 0xcf, 0xc9, 0xff, 0x10, 0x69, 0x31,
	0x03, 0xb3, 0x4a, 0x5e, 0x29, 0xac, 0x95, 0xfe, 0x2a, 0x8e, 0xc1, 0x8b, 0x7e, 0x68, 0x51, 0x3e,
	0xca, 0xcc, 0x40, 0xcd, 0x8d, 0x21, 0xdb, 0x90, 0xb4, 0x91, 0x9f, 0x75, 0x91, 0x73, 0xfd, 0x02,
	0xb3, 0xa1, 0xbc, 0x52, 0x48, 0x68, 0x60, 0x23, 0xaf, 0xc9, 0x1d, 0xba, 0x07, 0x30, 0x0a, 0x22,
	0x49, 0x88, 0x35, 0x9a, 0xe5, 0x72, 0xa5, 0xd1, 0x48, 0x2d, 0x39, 0x8b, 0x57, 0xfb, 0x5a, 0x5d,
	0xad, 0x1f, 0xa6, 0x14, 0x67, 0x51, 0xdd, 0x57, 0x8f, 0x9b, 0x5a, 0x25, 0x15, 0xa2, 0x5b, 0x10,
	0x69, 0x36, 0xd5, 0x03, 0x92, 0x81, 0xe5, 0x81, 0x7e, 0xdd, 0x97, 0xd4, 0x56, 0x34, 0xb9, 0xa0,
	0x1d, 0xd8, 0x3c, 0x44, 0x71, 0x60, 0xda, 0xd8, 0x12, 0xcc, 0x1e, 0x96, 0x99, 0x25, 0xd0, 0x12,
	0x5c, 0xc3, 0x77, 0x7d, 0xe4, 0xc2, 0xa1, 0xc4, 0x91, 0x73, 0x93, 0x59, 0x67, 0x57, 0x38, 0x74,
	0x43, 0x13, 0x1a, 0x78, 0x5b, 0x47, 0x38, 0x24, 0x05, 0x88, 0x1a, 0xa6, 0x7d, 0x66, 0x1a, 0x2e,
	0xdd, 0x64, 0x29, 0xed, 0x53, 0xec, 0x24, 0xd6, 0x96, 0x0d, 0xd3, 0x56, 0x0d, 0xfa, 0x4d, 0x01,
	0x38, 0x39, 0xbf, 0xc4, 0x96, 0x50, 0xad, 0x36, 0x23, 0x45, 0x48, 0x30, 0x77, 0xe5, 0xc4, 0x2a,
	0xb3, 0x62, 0xe3, 0x12, 0xa3, 0x1a, 0xa4, 0x0c, 0x49, 0x0f, 0x2f, 0x86, 0x3d, 0x69, 0xce, 0x5a,
	0x89, 0xfa, 0x22, 0x46, 0xa7, 0x7b, 0xaf, 0xa7, 0xc3, 0x1e, 0x6a, 0xc0, 0xee, 0xde, 0x1d, 0x39,
	0xde, 0x21, 0x96, 0xde, 0xc5, 0x6c, 0x58, 0xca, 0x91, 0x5b, 0x75, 0xbd, 0x8b, 0xf4, 0xcf, 0x5b,
	0x8e, 0x2e, 0x3c, 0x0e, 0x91, 0xaa, 0x7a, 0x5c, 0x49, 0x2d, 0x91, 0x55, 0x48, 0x1c, 0xa8, 0x5a,
	0xa5, 0x7c, 0x7a, 0xa2, 0xbd, 0x4e, 0x29, 0xf4, 0xa3, 0x02, 0x5b, 0xd3, 0x6d, 0x93, 0x37, 0x4c,
	0xf6, 0x20, 0xca, 0xdd, 0x9b, 0xf2, 0xa4, 0x6d, 0xce, 0x29, 0x04, 0xcd, 0x83, 0x92, 0xff, 0x20,
	0x26, 0xa9, 0xf0, 0x6c, 0x28, 0x1f, 0x2e, 0x24, 0x4b, 0xeb, 0x33, 0xe4, 0x69, 0xb7, 0x38, 0xfa,
	0x1e, 0x32, 0x35, 0xfd, 0x0a, 0xef, 0x88, 0xfc, 0xf8, 0x7b, 0x23, 0x1b, 0x10, 0x77, 0x90, 0x63,
	0x86, 0xc5, 0x0c, 0xd3, 0x76, 0xdd, 0x1a, 0xc0, 0x2f, 0x81, 0xec, 0x8b, 0xc8, 0x7f, 0x78, 0x29,
	0x7d, 0x56, 0x20, 0xbd, 0x6f, 0x18, 0x75, 0xbc, 0xa9, 0x9a, 0xd7, 0xf8, 0x04, 0x9a, 0x37, 0x21,
	0xd1, 0x36, 0xaf, 0x71, 0x5c, 0x74, 0xdc, 0xd9, 0x70, 0x54, 0x93, 0x1c, 0xc4, 0x5b, 0xde, 0x7d,
	0x67, 0x23, 0x6e, 0x2f, 0xdd, 0xad, 0x69, 0x1f, 0xc8, 0x38, 0xb1, 0x45, 0xec, 0xf8, 0x17, 0x62,
	0x2e, 0x87, 0x79, 0x74, 0xa3, 0x0e, 0x42, 0x35, 0xe8, 0x07, 0x05, 0x36, 0x9a, 0x3d, 0x43, 0x17,
	0xe8, 0xe4, 0x7d, 0x74, 0x13, 0x3f, 0x22, 0x95, 0x4f, 0x7d, 0x38, 0xa0, 0xfe, 0x05, 0xe4, 0xa6,
	0xb1, 0x58, 0xc0, 0x05, 0x8a, 0xf0, 0xeb, 0x21, 0x8a, 0xa7, 0x56, 0x45, 0x2f, 0x61, 0x7d, 0x22,
	0xcd, 0x22, 0x97, 0x37, 0xee, 0x52, 0x28, 0xe0, 0xd2, 0x57, 0x05, 0xd2, 0x35, 0x36, 0x40, 0xd9,
	0xcf, 0x0f, 0x96, 0xe3, 0x1b, 0x98, 0xa1, 0xfb, 0x07, 0xe6, 0xa8, 0xda, 0xc3, 0xf7, 0x77, 0xb8,
	0x85, 0x37, 0xb2, 0xd8, 0x23, 0xb2, 0xc3, 0x2d, 0xbc, 0x71, 0x3b, 0x7c, 0x08, 0x64, 0x9c, 0xea,
	0x22, 0x96, 0x3c, 0x92, 0x3f, 0x6d, 0xc3, 0xcf, 0x1a, 0x76, 0x9f, 0xdc, 0x27, 0x7a, 0x04, 0x19,
	0x7f, 0x9e, 0x05, 0x44, 0x96, 0x3e, 0x2d, 0xc3, 0xea, 0x73, 0xc6, 0x44, 0x95, 0x37, 0xd0, 0x1e,
	0x98, 0x2d, 0x24, 0x57, 0x90, 0x99, 0xf6, 0xa5, 0x20, 0x05, 0xdf, 0x71, 0x73, 0xbe, 0xc1, 0xb9,
	0x7f, 0x1e, 0x80, 0xf4, 0x38, 0xbf, 0x84, 0x55, 0xdf, 0x40, 0x26, 0x3b, 0xbe, 0xd8, 0x69, 0x9f,
	0x8a, 0x1c, 0x9d, 0x07, 0xf1, 0xce, 0xad, 0x01, 0x8c, 0xc6, 0x1a, 0xf9, 0xcd, 0x17, 0x31, 0x31,
	0x88, 0x73, 0xdb, 0x33, 0x7f, 0xf7, 0x8e, 0x43, 0x20, 0x93, 0x73, 0x82, 0xf8, 0xff, 0x2c, 0xcd,
	0x1c, 0x67, 0xb9, 0xbf, 0xef, 0xc5, 0x79, 0x69, 0xde, 0xc2, 0x4f, 0x81, 0xa6, 0x26, 0xbf, 0x07,
	0xbd, 0x9c, 0x96, 0xe0, 0x8f, 0xf9, 0xa0, 0x91, 0x27, 0xa3, 0xd6, 0x08, 0x78, 0x32, 0xd1, 0xde,
	0x01, 0x4f, 0xa6, 0xf4, 0x54, 0x03, 0x56, 0xc6, 0xcb, 0x90, 0xe4, 0x03, 0xe5, 0x36, 0xd1, 0x09,
	0xb9, 0x9d, 0x39, 0x08, 0x79, 0xe8, 0xb3, 0x9d, 0x37, 0xdb, 0x17, 0xa6, 0xe8, 0xf4, 0xcf, 0x8b,
	0x2d, 0xd6, 0xdd, 0x95, 0xf0, 0xdb, 0x87, 0xfc, 0x6f, 0x7b, 0x1e, 0x75, 0x9f, 0x7b, 0xdf, 0x03,
	0x00, 0x00, 0xff, 0xff, 0x69, 0xce, 0x48, 0x35, 0xec, 0x0a, 0x00, 0x00,
}
