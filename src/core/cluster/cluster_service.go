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
)

const (
	port = ":50052"
)

var ErrUnimplemented = errors.New("Unimplemented")
var ErrMessageFailed = errors.New("Message was not sent.")

type ClusterServer struct {
	// Is this atomic though????
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
	s := grpc.NewServer()

	hootpb.RegisterClusterServiceServer(s, &ClusterServer{})
	log.Printf("Server lsitening at %v", lis.Addr())

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
	// For tabbing
	new_id, err := uuid.FromBytes(request.NewFileId.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to get proper file ID")
	}
	par_id, err := uuid.FromBytes(request.ParentDirId.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to get proper parent ID")
	}

	// Add new file to virtual file manager
	c.vfmg.AddNewFile(
		hootfs.VirtualFile{
			Name: request.NewFileName,
			Id:   new_id},
		par_id)

	new_file_info := hootfs.FileInfo{NamespaceId: request.UserId, ObjectId: new_id}
	c.fmg.CreateFile(request.NewFileName, &new_file_info)
	if err := c.fmg.WriteFile(&new_file_info, request.Contents); err != nil {
		return nil, fmt.Errorf("Failed to write to file: %v", err)
	}

	return &hootpb.AddNewFileCSResponse{CreatedFileId: &hootpb.UUID{Value: new_id[:]}}, nil
}

func (c *ClusterServer) MakeDirectoryCS(ctx context.Context,
	request *hootpb.MakeDirectoryCSRequest) (*hootpb.MakeDirectoryCSResponse, error) {

	return nil, ErrUnimplemented
}

func (c *ClusterServer) UpdateFileContentsCS(ctx context.Context,
	request *hootpb.UpdateFileContentsCSRequest) (*hootpb.UpdateFileContentsCSResponse, error) {
	return nil, ErrUnimplemented
}

func (c *ClusterServer) MoveObjectCS(ctx context.Context, request *hootpb.MoveObjectCSRequest) (*hootpb.MoveObjectCSResponse, error) {
	return nil, ErrUnimplemented
}

func (c *ClusterServer) RemoveObjectCS(ctx context.Context, request *hootpb.RemoveObjectCSRequest) (*hootpb.RemoveObjectCSResponse, error) {
	return nil, ErrUnimplemented
}
