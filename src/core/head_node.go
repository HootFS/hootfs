package core

import (
	"context"
    "time"
    "net"
    "log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	protos "github.com/hootfs/hootfs/protos"
	head "github.com/hootfs/hootfs/protos"
    discover "github.com/hootfs/hootfs/src/discover"
    cluster "github.com/hootfs/hootfs/src/core/cluster"
    hootfs "github.com/hootfs/hootfs/src/core/file_storage"
)

const (
    headPort = ":50060" 
    nodePingDurr = 1 
    nodeGetActiveDurr = 20
)

// File manager Server must deal with 
// taking requests and pinging discovery server.

type fileManagerServer struct {
    csc cluster.ClusterServiceClient
    dc discover.DiscoverClient
    
    vfmp *hootfs.VirtualFileMapper
    vfmg *hootfs.VirtualFileManager 

	head.UnimplementedHootFsServiceServer
}

func NewFileManagerServer(dip string, vfmp *hootfs.VirtualFileMapper, 
    vfmg *hootfs.VirtualFileManager) *fileManagerServer {
    return &fileManagerServer{
        dc: *discover.NewDiscoverClient(dip),
        vfmp: vfmp,
        vfmg: vfmg,
    }
}

func (fms *fileManagerServer) StartServer() error {
    // First start server.
    lis, err := net.Listen("tcp", headPort)
    if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

    protos.RegisterHootFsServiceServer(s, fms)

    // Join the discovery server.
    nodeId, clusterMap, err := fms.dc.JoinCluster()

    if err != nil {
        return err
    }

    fms.csc = *cluster.NewClusterServiceClient(nodeId)
    fms.csc.UpdateNodes(clusterMap)

    // Ping function.
    go func() {
        for {
            time.Sleep(nodePingDurr * time.Second)
            err := fms.dc.Ping()

            if err != nil {
                // Error case!
                // Not sure what to do here
                // if we cannot ping the discovery 
                // server.
            }
        }
    }()

    // Get Active update function.
    go func() {
        for {
            time.Sleep(nodeGetActiveDurr * time.Second)
            clusterMap, err := fms.dc.GetActive()

            if err != nil {
                // TODO
            } else {
                fms.csc.UpdateNodes(clusterMap)
            }
        }
    }()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

    // This should never be reached since Serve is blocking.
    return nil
}

func (s *fileManagerServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *fileManagerServer) MakeDirectory(
	ctx context.Context, request *head.MakeDirectoryRequest) (*head.MakeDirectoryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *fileManagerServer) AddNewFile(
	ctx context.Context, request *head.AddNewFileRequest) (*head.AddNewFileResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
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


