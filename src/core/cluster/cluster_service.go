package cluster

import (
	"context"
	"errors"
	"log"
	"net"

	hootpb "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc"
    hootfs "github.com/hootfs/hootfs/src/core/file_storage"
)

const (
	port = ":50052"
)

var ErrUnimplemented = errors.New("Unimplemented")
var ErrMessageFailed = errors.New("Message was not sent.")

type ClusterServer struct {
    vfmp *hootfs.VirtualFileMapper
    vfmg *hootfs.VirtualFileManager 

	hootpb.UnimplementedClusterServiceServer
}

func NewClusterServer(vfmp *hootfs.VirtualFileMapper, vfmg *hootfs.VirtualFileManager) *ClusterServer {
    return &ClusterServer{
        vfmp: vfmp,
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

	return nil, ErrUnimplemented
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
