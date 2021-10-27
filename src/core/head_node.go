package core

import (
	"context"
	"log"
	"net"
	"time"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	head "github.com/hootfs/hootfs/protos"
	protos "github.com/hootfs/hootfs/protos"
	cluster "github.com/hootfs/hootfs/src/core/cluster"
	hootfs "github.com/hootfs/hootfs/src/core/file_storage"
	discover "github.com/hootfs/hootfs/src/discover"
	"google.golang.org/grpc"
)

const (
	headPort          = ":50060"
	nodePingDurr      = 1
	nodeGetActiveDurr = 20
)

// File manager Server must deal with
// taking requests and pinging discovery server.

type HootFsServer struct {
	csc cluster.ClusterServiceClient
	dc  discover.DiscoverClient

	fmg  *hootfs.FileManager
	vfmg *hootfs.VirtualFileManager

	head.UnimplementedHootFsServiceServer
}

func NewHootFsServer(dip string, fmg *hootfs.FileManager,
	vfmg *hootfs.VirtualFileManager) *HootFsServer {
	return &HootFsServer{
		dc:   *discover.NewDiscoverClient(dip),
		fmg:  fmg,
		vfmg: vfmg,
	}
}

func (fms *HootFsServer) StartServer() error {
	// First start server.
	lis, err := net.Listen("tcp", headPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

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

func (s *HootFsServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {
	dirUuid, err := uuid.FromBytes(request.DirId.Value)

	if err != nil {
		return nil, err
	}

	contents, err := s.vfmg.GetDirectoryContentsAsProto(dirUuid)

	if err != nil {
		return nil, err
	}

	return &protos.GetDirectoryContentsResponse{
		Objects: contents,
	}, nil
}

func (s *HootFsServer) MakeDirectory(
	ctx context.Context, request *head.MakeDirectoryRequest) (*head.MakeDirectoryResponse, error) {
	parentUuid, err := uuid.FromBytes(request.DirId.Value)

	if err != nil {
		return nil, err
	}

	dirUuid, err := s.vfmg.CreateNewDirectory(request.DirName, parentUuid)

	if err != nil {
		return nil, err
	}

	// Broadcast directory creation to all other clients.
	for destId := range s.csc.Nodes {
		if destId != s.csc.NodeId {
			// NOTE, In the future we will need some form of error handling here.
			s.csc.SendMakeDirectory(destId, "USERID", parentUuid, dirUuid, request.DirName)
		}
	}

	return &head.MakeDirectoryResponse{
		DirId: &protos.UUID{
			Value: dirUuid[:],
		},
	}, nil
}

func (s *HootFsServer) AddNewFile(
	ctx context.Context, request *head.AddNewFileRequest) (*head.AddNewFileResponse, error) {
	parentUuid, err := uuid.FromBytes(request.DirId.Value)

	if err != nil {
		return nil, err
	}

	fileUuid, err := s.vfmg.CreateNewFile(request.FileName, parentUuid)

	if err != nil {
		return nil, err
	}

	// Send make new file request to all cluster nodes.
	for destId := range s.csc.Nodes {
		if destId != s.csc.NodeId {
			s.csc.SendAddFile(destId, "USERID", parentUuid, fileUuid, request.FileName, request.Contents)
		}
	}

	newFileInfo := hootfs.FileInfo{NamespaceId: "USERID", ObjectId: fileUuid}

	// Local machine work... could throw an error, but this is OK as long
	// as file is stored on some machine??
	s.fmg.CreateFile(request.FileName, &newFileInfo)
	s.fmg.WriteFile(&newFileInfo, request.Contents)

	return &protos.AddNewFileResponse{
		FileId: &protos.UUID{Value: fileUuid[:]},
	}, nil
}

func (s *HootFsServer) UpdateFileContents(
	ctx context.Context, request *head.UpdateFileContentsRequest) (*head.UpdateFileContentsResponse, error) {
	fileUuid, err := uuid.FromBytes(request.FileId.Value)

	if err != nil {
		return nil, err
	}

	// In theory, this file should exist on at least one machine
	// if we are updating it...
	for destId := range s.csc.Nodes {
		if destId != s.csc.NodeId {
			s.csc.SendUpdateFileContentsRequest(destId, "USERID", fileUuid, request.Contents)
		}
	}

	newFileInfo := hootfs.FileInfo{NamespaceId: "USERID", ObjectId: fileUuid}

	// If this file does not exist on this machine, an error may be thrown here.
	// This is not a big deal, since the file should exist on another machine if we are
	// updating...
	// This issue we will need to flesh out later when we have more time.

	s.fmg.WriteFile(&newFileInfo, request.Contents)
	return &protos.UpdateFileContentsResponse{}, nil
}

func (s *HootFsServer) GetFileContents(
	ctx context.Context, request *head.GetFileContentsRequest) (*head.GetFileContentsResponse, error) {
	fileUuid, err := uuid.FromBytes(request.FileId.Value)

	if err != nil {
		return nil, err
	}

	newFileInfo := hootfs.FileInfo{NamespaceId: "USERID", ObjectId: fileUuid}
	contents, err := s.fmg.ReadFile(&newFileInfo)

	// This is the case where the file is on the given machine.
	if err == nil {
		return &head.GetFileContentsResponse{Contents: contents}, nil
	}

	// TODO : In the future, we should search other machines for missing file.
	// for now we will just search this machine only... static cluster size.
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

// Both Move object and remove object require the parent IDs of the the objects being
// modified... we don't have access to this at this moment.

func (s *HootFsServer) MoveObject(
	ctx context.Context, request *head.MoveObjectRequest) (*head.MoveObjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *HootFsServer) RemoveObject(
	ctx context.Context, request *head.RemoveObjectRequest) (*head.RemoveObjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Method not implemented")
}
