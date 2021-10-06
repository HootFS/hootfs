package hootfs

import (
	"context"
	"errors"

	head "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fileManagerServer struct {
}

func (s *fileManagerServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {

	return &head.GetDirectoryContentsResponse{}, status.Error(codes.Unimplemented,
		"Unimplemented")
}

func (s *fileManagerServer) MakeDirectory(
	ctx context.Context, request *head.MakeDirectoryRequest) (*head.MakeDirectoryResponse, error) {
	return &head.MakeDirectoryResponse{}, status.Error(codes.Unimplemented, "")
}

func (s *fileManagerServer) AddNewFile(
	ctx context.Context, request *head.AddNewFileRequest) (*head.AddNewFileResponse, error) {
	return &head.AddNewFileResponse{}, nil
}

func (s *fileManagerServer) UpdateFileContents(
	ctx context.Context, request *head.UpdateFileContentsRequest) (*head.UpdateFileContentsResponse, error) {
	return &head.UpdateFileContentsResponse{}, nil
}

func (s *fileManagerServer) GetFileContents(
	ctx context.Context, request *head.GetFileContentsRequest) (*head.GetFileContentsResponse, error) {
	return &head.GetFileContentsResponse{}, nil
}
func (s *fileManagerServer) MoveObject(
	ctx context.Context, request *head.MoveObjectRequest) (*head.MoveObjectResponse, error) {
	return &head.MoveObjectResponse{}, nil
}
func (s *fileManagerServer) RemoveObject(
	ctx context.Context, request *head.RemoveObjectRequest) (*head.RemoveObjectResponse, error) {
	return &head.RemoveObjectResponse{}, nil
}

// The number of bytes in a UUID (For now)
const UUIDSize uint8 = 2

// Key into the file mapping.
type localUUID struct {
    value [UUIDSize]byte
}

type object interface {
    getName() string

    // Directory actions.
    getDirectoryContents() ([]localUUID, error)

    // File actions.
    getFileContents() ([]byte, error)
}


// Object representing a directory.
type directoryObject struct {
    name string
    contents []localUUID
}

func (d *directoryObject) getName() string {
    return d.name
}

func (d *directoryObject) getDirectoryContents() ([]localUUID, error) {
    return d.contents, nil
}

func (d *directoryObject) getFileContents() ([]byte, error) {
    return nil, errors.New("Cannot get file contents from a direcotry!")
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

func (fo *fileObject) getDirectoryContents() ([]localUUID, error) {
    return nil, errors.New("Cannot get the contents of a file!"
}

func (fo *fileObject getFileContents() ([]byte, error) {
    // Actuall perform read!
}

