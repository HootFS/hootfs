package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	google_verifer "github.com/hootfs/hootfs/src/core/auth"
	"github.com/hootfs/hootfs/src/core/vfm"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	head "github.com/hootfs/hootfs/protos"
	cluster "github.com/hootfs/hootfs/src/core/cluster"
	hootfs "github.com/hootfs/hootfs/src/core/file_storage"
	discover "github.com/hootfs/hootfs/src/discovery/discover"

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

	fmg     *hootfs.FileManager
	vfm     vfm.VirtualFileManager
	vfmg    *hootfs.VirtualFileManager
	verify  *google_verifer.Google_verifier
	ns_root uuid.UUID

	head.UnimplementedHootFsServiceServer
}

func NewHootFsServer(dip string, fmg *hootfs.FileManager,
	vfmg *hootfs.VirtualFileManager) *HootFsServer {
	meta_vfm, err := vfm.NewMetaStore("PROD")
	if err != nil {
		return nil
	}

	return &HootFsServer{
		dc:      *discover.NewDiscoverClient(dip),
		fmg:     fmg,
		vfm:     meta_vfm,
		vfmg:    vfmg,
		verify:  nil,
		ns_root: uuid.Nil,
	}
}

func (fms *HootFsServer) GetDirectoryContentsAsProto(dirId uuid.UUID) ([]*head.ObjectInfo, error) {
	fms.vfmg.RWLock.RLock()
	defer fms.vfmg.RWLock.RUnlock()

	vd, exists := fms.vfmg.Directories[dirId]
	if !exists {
		return nil, fmt.Errorf("Directory not found!")
	}

	contents := make([]*head.ObjectInfo, len(vd.Subdirs)+len(vd.Files))

	for dirUuid := range vd.Subdirs {
		contents = append(contents, &head.ObjectInfo{
			ObjectId:   &head.UUID{Value: dirUuid[:]},
			ObjectType: head.ObjectInfo_DIRECTORY,
			ObjectName: fms.vfmg.Directories[dirUuid].Name,
		})
	}

	for fileUuid := range vd.Files {
		contents = append(contents, &head.ObjectInfo{
			ObjectId:   &head.UUID{Value: fileUuid[:]},
			ObjectType: head.ObjectInfo_FILE,
			ObjectName: fms.vfmg.Files[fileUuid].Name,
		})
	}

	return contents, nil
}

func (fms *HootFsServer) StartServer() error {
	// First start server.
	lis, err := net.Listen("tcp", headPort)
	log.Println("Starting Server")
	verifier := google_verifer.New("https://dev-dewy8ew9.us.auth0.com/userinfo")

	// Store verifier in the server.
	fms.verify = verifier

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(verifier.Authenticate)))
	s := grpc.NewServer(opts...)
	log.Println("Server started")
	head.RegisterHootFsServiceServer(s, fms)

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
				// Not sure what to do here if we cannot ping the discovery
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

/**
 * Initializes a HootFS client, returning the initialized user's namespace ID
 */
func (s *HootFsServer) InitializeHootfsClient(
	ctx context.Context, request *head.InitializeHootfsClientRequest) (*head.InitializeHootfsClientResponse, error) {

	username := s.verify.GetUsername(ctx)
	ns_stubs, err := s.vfm.GetNamespaces(vfm.User_ID(username))

	if err == vfm.ErrUserDoesNotExist {
		// If the user does not exist, we will make it a namespace.
		err = s.vfm.CreateUser(vfm.User_ID(username))
		if err != nil {
			return &head.InitializeHootfsClientResponse{NamespaceRoot: nil}, err
		}

		nsid, err := s.vfm.CreateNamespace("User Namespace",
			vfm.User_ID(username))
		if err != nil {
			return &head.InitializeHootfsClientResponse{NamespaceRoot: nil}, err
		}

		// Then we will make a root folder in the namespace.
		void, err := s.vfm.CreateFreeObjectInNamespace(nsid,
			vfm.User_ID(username), "Root", vfm.VFM_Dir_Type)
		if err != nil {
			return &head.InitializeHootfsClientResponse{NamespaceRoot: nil}, err
		}

		return &head.InitializeHootfsClientResponse{
			NamespaceRoot: &head.UUID{Value: void[:]},
		}, nil
	} else if err != nil {
		return &head.InitializeHootfsClientResponse{NamespaceRoot: nil}, err
	}

	// Otherwise, this user should have only one namespace for now.
	sole_ns := ns_stubs[0]

	ns, err := s.vfm.GetNamespaceDetails(sole_ns.NSID, vfm.User_ID(username))

	if err != nil {
		return &head.InitializeHootfsClientResponse{NamespaceRoot: nil}, err
	}

	return &head.InitializeHootfsClientResponse{
		NamespaceRoot: &head.UUID{Value: ns.RootObjects[0][:]},
	}, nil
}

func (s *HootFsServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {
	dirUuid, err := uuid.FromBytes(request.DirId.Value)

	if err != nil {
		return &head.GetDirectoryContentsResponse{},
			status.Error(codes.InvalidArgument, cluster.ErrInvalidId.Error())
	}

	username := s.verify.GetUsername(ctx)
	contents, err := s.vfm.GetObjectDetails(vfm.VO_ID(dirUuid),
		vfm.User_ID(username))

	if err != nil {
		return &head.GetDirectoryContentsResponse{}, err
	}

	var sos []*head.ObjectInfo
	subcontents, err := contents.GetSubObjects()

	for _, sc := range subcontents {
		ot := head.ObjectInfo_DIRECTORY
		if sc.Type == vfm.VFM_File_Type {
			ot = head.ObjectInfo_FILE
		}

		sos = append(sos, &head.ObjectInfo{
			ObjectId:   &head.UUID{Value: sc.Id[:]},
			ObjectType: ot,
			ObjectName: sc.Name,
		})
	}

	return &head.GetDirectoryContentsResponse{
		Objects: sos,
	}, nil

	// contents, err := s.GetDirectoryContentsAsProto(dirUuid)
	// if err != nil {
	// 	return &head.GetDirectoryContentsResponse{}, err
	// }

	// if err != nil {
	// 	return &head.GetDirectoryContentsResponse{},
	// 		status.Error(codes.Internal, fmt.Sprintf("Unable to get directory contents: %v", err))
	// }

	// return &head.GetDirectoryContentsResponse{
	// 	Objects: contents,
	// }, nil
}

func (s *HootFsServer) MakeDirectory(
	ctx context.Context, request *head.MakeDirectoryRequest) (*head.MakeDirectoryResponse, error) {
	parentUuid, err := uuid.FromBytes(request.DirId.Value)

	if err != nil {
		return &head.MakeDirectoryResponse{},
			status.Error(codes.InvalidArgument, cluster.ErrInvalidId.Error())
	}

	username := s.verify.GetUsername(ctx)

	// Directories only need to be made in the permanent store.
	void, err := s.vfm.CreateObject(vfm.VO_ID(parentUuid),
		vfm.User_ID(username), request.DirName, vfm.VFM_Dir_Type)

	if err != nil {
		return &head.MakeDirectoryResponse{}, err
	}

	return &head.MakeDirectoryResponse{
		DirId: &head.UUID{Value: void[:]},
	}, nil

	// dirUuid, err := s.vfmg.CreateNewDirectory(request.DirName, parentUuid)
	// if err != nil {
	// 	return &head.MakeDirectoryResponse{},
	// 		status.Error(codes.Internal, fmt.Sprintf("Failed to create directory: %v", err))
	// }

	// // Broadcast directory creation to all other clients.
	// for destId := range s.csc.Nodes {
	// 	if destId != s.csc.NodeId {
	// 		// NOTE, In the future we will need some form of error handling here.
	// 		s.csc.SendMakeDirectory(destId, "USERID", parentUuid, dirUuid, request.DirName)
	// 	}
	// }

	// return &head.MakeDirectoryResponse{
	// 	DirId: &head.UUID{
	// 		Value: dirUuid[:],
	// 	},
	// }, nil
}

func (s *HootFsServer) AddNewFile(
	ctx context.Context, request *head.AddNewFileRequest) (*head.AddNewFileResponse, error) {
	parentUuid, err := uuid.FromBytes(request.DirId.Value)

	if err != nil {
		return &head.AddNewFileResponse{},
			status.Error(codes.InvalidArgument, cluster.ErrInvalidId.Error())
	}

	username := s.verify.GetUsername(ctx)

	// Create the file in the permanent store.
	void, err := s.vfm.CreateObject(vfm.VO_ID(parentUuid),
		vfm.User_ID(username), request.FileName, vfm.VFM_File_Type)

	if err != nil {
		return &head.AddNewFileResponse{}, nil
	}

	// fileUuid, err := s.vfmg.CreateNewFile(request.FileName, parentUuid)

	// if err != nil {
	// 	return &head.AddNewFileResponse{},
	// 		status.Error(codes.Internal, fmt.Sprintf("Error creating new file: %v", err))
	// }

	// Send make new file request to all cluster nodes.
	for destId := range s.csc.Nodes {
		if destId != s.csc.NodeId {
			s.csc.SendAddFile(destId, username, parentUuid,
				uuid.UUID(void), request.FileName, request.Contents)
		}
	}

	newFileInfo := hootfs.FileInfo{NamespaceId: username,
		ObjectId: uuid.UUID(void)}

	// Local machine work... could throw an error, but this is OK as long
	// as file is stored on some machine??
	s.fmg.CreateFile(request.FileName, &newFileInfo)
	s.fmg.WriteFile(&newFileInfo, request.Contents)

	return &head.AddNewFileResponse{
		FileId: &head.UUID{Value: void[:]},
	}, nil
}

func (s *HootFsServer) UpdateFileContents(
	ctx context.Context, request *head.UpdateFileContentsRequest) (*head.UpdateFileContentsResponse, error) {
	fileUuid, err := uuid.FromBytes(request.FileId.Value)

	if err != nil {
		return &head.UpdateFileContentsResponse{},
			status.Error(codes.InvalidArgument, cluster.ErrInvalidId.Error())
	}

	username := s.verify.GetUsername(ctx)

	// In theory, this file should exist on at least one machine
	// if we are updating it...
	for destId := range s.csc.Nodes {
		if destId != s.csc.NodeId {
			s.csc.SendUpdateFileContentsRequest(destId, username, fileUuid, request.Contents)
		}
	}

	newFileInfo := hootfs.FileInfo{NamespaceId: username, ObjectId: fileUuid}

	// If this file does not exist on this machine, an error may be thrown here.
	// This is not a big deal, since the file should exist on another machine
	// if we are
	// updating...
	// This issue we will need to flesh out later when we have more time.

	s.fmg.WriteFile(&newFileInfo, request.Contents)
	return &head.UpdateFileContentsResponse{}, nil
}

func (s *HootFsServer) GetFileContents(
	ctx context.Context, request *head.GetFileContentsRequest) (*head.GetFileContentsResponse, error) {
	fileUuid, err := uuid.FromBytes(request.FileId.Value)

	if err != nil {
		return &head.GetFileContentsResponse{}, status.Error(codes.InvalidArgument, cluster.ErrInvalidId.Error())
	}

	username := s.verify.GetUsername(ctx)

	newFileInfo := hootfs.FileInfo{NamespaceId: username, ObjectId: fileUuid}
	contents, err := s.fmg.ReadFile(&newFileInfo)

	// This is the case where the file is on the given machine.
	if err == nil {
		return &head.GetFileContentsResponse{Contents: contents}, nil
	}

	// TODO : In the future, we should search other machines for missing file.
	// for now we will just search this machine only... static cluster size.
	return &head.GetFileContentsResponse{}, status.Error(codes.Unimplemented, "Method not implemented")
}

// Both Move object and remove object require the parent IDs of the the objects
// being modified... we don't have access to this at this moment.

func (s *HootFsServer) MoveObject(
	ctx context.Context, request *head.MoveObjectRequest) (*head.MoveObjectResponse, error) {
	return &head.MoveObjectResponse{},
		status.Error(codes.Unimplemented, "Method not implemented")
}

func (s *HootFsServer) RemoveObject(
	ctx context.Context, request *head.RemoveObjectRequest) (*head.RemoveObjectResponse, error) {
	return &head.RemoveObjectResponse{},
		status.Error(codes.Unimplemented, "Method not implemented")
}
