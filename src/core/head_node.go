package hootfs

import (
	"context"
	"errors"

	head "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fileManagerServer struct {
    mapping map[localUUID]object
}

func (s *fileManagerServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {
    dir, ok := s.mapping[toLocal(request.DirId)]

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
    
    for _, uuid := range contents {
        obj := s.mapping[uuid]

        objType := head.ObjectInfo_FILE
        if obj.isDirectory() {
            objType = head.ObjectInfo_DIRECTORY
        }

        objInfo := head.ObjectInfo{
            ObjectId: uuid.toProto(), 
            ObjectType: objType,
            ObjectName: obj.getName(),
        }

        resp.Objects = append(resp.Objects, &objInfo)
    }

	return &head.GetDirectoryContentsResponse{}, nil
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

/*
type fileMapping struct {
    mapping map[localUUID]object
}

var UnknownID = errors.New("Given ID is not defined!")

func (fm *fileMapping) getDirectoryContents(id localUUID) ([]localUUID, error) {
    dir, ok := fm.mapping[id]

    if !ok {
        return nil, UnknownID
    }

    return dir.getDirectoryContents()
}



func (fm *fileMapping) addFile(id localUUID) ([]localUUID, error) {
    return nil, nil
}
*/

// The number of bytes in a UUID (For now)
const UUIDSize uint8 = 16

// Key into the file mapping.
type localUUID struct {
    value [UUIDSize]byte
}

func (lu *localUUID) toProto() *head.UUID {
    return &head.UUID{Value: lu.value[:]}
}

func toLocal(pru *head.UUID) localUUID {
    var byts [UUIDSize]byte

    copy(byts[:], pru.Value)    

    return localUUID{value: byts}
}

type object interface {
    getName() string

    isDirectory() bool

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

func (d *directoryObject) isDirectory() bool {
    return true
}

func (d *directoryObject) getDirectoryContents() ([]localUUID, error) {
    return d.contents, nil
}

var NotAFile = errors.New("Cannot get file contents from a direcotry!")

func (d *directoryObject) getFileContents() ([]byte, error) {
    return nil, NotAFile 
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

func (fo *fileObject) isDirectory() bool {
    return false 
}

var NotADirectory = errors.New("Cannot get the contents of a file!")

func (fo *fileObject) getDirectoryContents() ([]localUUID, error) {
    return nil, NotADirectory
}

func (fo *fileObject) getFileContents() ([]byte, error) {
    // Actuall perform read!
}

