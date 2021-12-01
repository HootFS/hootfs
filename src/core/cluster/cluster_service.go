package cluster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	hootpb "github.com/hootfs/hootfs/protos"
	hootfs "github.com/hootfs/hootfs/src/core/file_storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50052"
)

var ErrUnimplemented = errors.New("Unimplemented")
var ErrMessageFailed = errors.New("Message was not sent.")
var ErrInvalidId = status.Error(codes.InvalidArgument, "Could not parse specified UUID")

type ClusterServer struct {
	fmg  *hootfs.FileManager
	vfmg *hootfs.VirtualFileManager

	hootpb.UnimplementedClusterServiceServer
}

func NewClusterServer(fmg *hootfs.FileManager, vfmg *hootfs.VirtualFileManager) *ClusterServer {
	return &ClusterServer{
		fmg:  fmg,
		vfmg: vfmg,
	}
}

func (c *ClusterServer) StartServer() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	hootpb.RegisterClusterServiceServer(s, &ClusterServer{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Sets up the cluster router by starting up GRPC server
func (c *ClusterServer) Init() error {
	go c.StartServer()
	return nil
}

func (c *ClusterServer) AddNewFileCS(ctx context.Context,
	request *hootpb.AddNewFileCSRequest) (*hootpb.AddNewFileCSResponse, error) {
	new_id, err := uuid.FromBytes(request.NewFileId.Value)
	if err != nil {
		return &hootpb.AddNewFileCSResponse{}, status.Error(codes.InvalidArgument, ErrInvalidId.Error())
	}
	par_id, err := uuid.FromBytes(request.ParentDirId.Value)
	if err != nil {
		return &hootpb.AddNewFileCSResponse{}, status.Error(codes.InvalidArgument, ErrInvalidId.Error())
	}

	// Add new file to virtual file manager
	c.vfmg.AddNewFile(
		&hootfs.VirtualFile{
			Name: request.NewFileName,
			Id:   new_id},
		par_id)

	new_file_info := hootfs.FileInfo{NamespaceId: request.UserId, ObjectId: new_id}
	c.fmg.CreateFile(request.NewFileName, &new_file_info)
	if err := c.fmg.WriteFile(&new_file_info, request.Contents); err != nil {
		return &hootpb.AddNewFileCSResponse{}, status.Error(
			codes.Internal, fmt.Sprintf("Failed to write  to file: %v", err))
	}

	return &hootpb.AddNewFileCSResponse{CreatedFileId: &hootpb.UUID{Value: new_id[:]}}, nil
}

func (c *ClusterServer) MakeDirectoryCS(ctx context.Context,
	request *hootpb.MakeDirectoryCSRequest) (*hootpb.MakeDirectoryCSResponse, error) {
	dir_id, err := uuid.FromBytes(request.NewDirId.Value)
	if err != nil {
		return &hootpb.MakeDirectoryCSResponse{},
			status.Error(codes.InvalidArgument, ErrInvalidId.Error())
	}
	par_id, err := uuid.FromBytes(request.ParentDirId.Value)
	if err != nil {
		return &hootpb.MakeDirectoryCSResponse{},
			status.Error(codes.InvalidArgument, ErrInvalidId.Error())
	}
	vir_dir := hootfs.VirtualDirectory{
		Name:    request.NewDirName,
		Id:      dir_id,
		Subdirs: make(map[uuid.UUID]bool),
		Files:   make(map[uuid.UUID]bool),
	}
	c.vfmg.AddNewDirectory(&vir_dir, par_id)
	return &hootpb.MakeDirectoryCSResponse{
		CreatedDirId: &hootpb.UUID{Value: dir_id[:]},
	}, nil
}

func (c *ClusterServer) UpdateFileContentsCS(ctx context.Context,
	request *hootpb.UpdateFileContentsCSRequest) (*hootpb.UpdateFileContentsCSResponse, error) {
	file_id, err := uuid.FromBytes(request.FileId.Value)
	if err != nil {
		return &hootpb.UpdateFileContentsCSResponse{}, status.Error(codes.InvalidArgument, ErrInvalidId.Error())
	}
	file_info := hootfs.FileInfo{NamespaceId: request.UserId, ObjectId: file_id}
	if err := c.fmg.WriteFile(&file_info, request.Contents); err != nil {
		return &hootpb.UpdateFileContentsCSResponse{}, status.Error(codes.Internal, fmt.Sprintf("Unable to write file %v: %v", file_info.ObjectId, err))
	}
	return &hootpb.UpdateFileContentsCSResponse{
		UpdatedFileId: &hootpb.UUID{Value: file_id[:]}}, nil
}

func (c *ClusterServer) MoveObjectCS(ctx context.Context, request *hootpb.MoveObjectCSRequest) (*hootpb.MoveObjectCSResponse, error) {
	obj_id, err := uuid.FromBytes(request.CurrObjectId.Value)
	if err != nil {
		return &hootpb.MoveObjectCSResponse{}, status.Error(codes.InvalidArgument, ErrInvalidId.Error())
	}

	if _, exists := c.vfmg.Directories[obj_id]; exists {

	}

	return &hootpb.MoveObjectCSResponse{}, status.Error(codes.Unimplemented, ErrUnimplemented.Error())
}

func (c *ClusterServer) RemoveObjectCS(ctx context.Context, request *hootpb.RemoveObjectCSRequest) (*hootpb.RemoveObjectCSResponse, error) {
	return &hootpb.RemoveObjectCSResponse{}, status.Error(codes.Unimplemented, ErrUnimplemented.Error())
}
