package core

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	head "github.com/hootfs/hootfs/protos"
)

type fileManagerServer struct {
	head.UnimplementedHootFsServiceServer
	directories map[*localUUID]directoryObject
	files       map[*localUUID]fileObject
}

func (s *fileManagerServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {
	dir, ok := s.directories[protoToLocalUUID(request.DirId)]

	if !ok {
		return &head.GetDirectoryContentsResponse{}, status.Error(codes.InvalidArgument,
			"Given ID is not defined!")
	}

	contents, err := dir.getDirectoryContents()

	if err != nil {
		return &head.GetDirectoryContentsResponse{}, status.Error(codes.InvalidArgument,
			"File is not a directory!")
	}

	// Success.
	resp := head.GetDirectoryContentsResponse{}

	for _, id := range *contents {
		objInfo := head.ObjectInfo{
			ObjectId: id.toProto(),
		}

		child_dir, exists := s.directories[&id]
		if exists {
			objInfo.ObjectName = child_dir.name
			objInfo.ObjectType = head.ObjectInfo_DIRECTORY
		}

		child_file, exists := s.files[&id]
		if exists {
			objInfo.ObjectName = child_file.name
			objInfo.ObjectType = head.ObjectInfo_FILE
		}

		resp.Objects = append(resp.Objects, &objInfo)
	}

	return &head.GetDirectoryContentsResponse{}, nil
}

func (s *fileManagerServer) MakeDirectory(
	ctx context.Context, request *head.MakeDirectoryRequest) (*head.MakeDirectoryResponse, error) {
	parent_dir, exists := s.directories[protoToLocalUUID(request.DirId)]

	if !exists {
		return nil, status.Error(codes.InvalidArgument, "Requested parent directory does not exist")
	}

	// Make directory somewhere
	directory_uuid, err := uuid.NewUUID()
	if err != nil {
		// Report error and return.
	}
	new_uuid := localUUID{value: directory_uuid}
	parent_dir.directories = append(parent_dir.directories, new_uuid)

	// s.mapping[uuid.]
	return &head.MakeDirectoryResponse{DirId: new_uuid.toProto()}, nil
}

func (s *fileManagerServer) AddNewFile(
	ctx context.Context, request *head.AddNewFileRequest) (*head.AddNewFileResponse, error) {
	parent_dir, exists := s.directories[protoToLocalUUID(request.DirId)]

	if !exists {
		return nil, status.Error(codes.InvalidArgument, "Requested parent directory does not exist")
	}

	file_uuid, err := uuid.NewUUID()
	if err != nil {
		// Report error and return
	}
	new_uuid := localUUID{value: file_uuid}
	parent_dir.files = append(parent_dir.files, new_uuid)

	return &head.AddNewFileResponse{}, nil
}

func (s *fileManagerServer) UpdateFileContents(
	ctx context.Context, request *head.UpdateFileContentsRequest) (*head.UpdateFileContentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *fileManagerServer) GetFileContents(
	ctx context.Context, request *head.GetFileContentsRequest) (*head.GetFileContentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *fileManagerServer) MoveObject(
	ctx context.Context, request *head.MoveObjectRequest) (*head.MoveObjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *fileManagerServer) RemoveObject(
	ctx context.Context, request *head.RemoveObjectRequest) (*head.RemoveObjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

// The number of bytes in a UUID (For now)
const UUIDSize uint8 = 16

// Key into the file mapping.
type localUUID struct {
	value [UUIDSize]byte
}

func (lu *localUUID) toProto() *head.UUID {
	return &head.UUID{Value: lu.value[:]}
}

func protoToLocalUUID(pru *head.UUID) *localUUID {
	var byts [UUIDSize]byte

	copy(byts[:], pru.Value)
	return &localUUID{value: byts}
}

// Object representing a directory.
type directoryObject struct {
	name        string
	directories []localUUID
	files       []localUUID
}

func (d *directoryObject) getName() string {
	return d.name
}

func (d *directoryObject) getDirectoryContents() (*[]localUUID, error) {
	out := []localUUID{}
	copy(out, d.directories)
	out = append(out, d.files...)
	return &out, nil
}

type fileObject struct {
	name string

	// NOTE, this will be taken out later.
	// right now since everything is stored on one
	// machine, just a local path is needed to locate it.
	localPath string
}

func (fo *fileObject) getName() string {
	return fo.name
}

func (fo *fileObject) getFileContents() ([]byte, error) {
	return nil, nil
	// Actually perform read!
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	head.RegisterHootFsServiceServer(s, &fileManagerServer{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
