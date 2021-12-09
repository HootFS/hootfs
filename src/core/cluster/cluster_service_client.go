package cluster

import (
	"context"
	"errors"
	"log"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/hootfs/hootfs/protos"
	hootpb "github.com/hootfs/hootfs/protos"

	"google.golang.org/grpc"
)

var ErrIncorrectFileCreated = errors.New("The wrong file ID was returned")

type ClusterServiceClient struct {
	// This will be used to make all operations on the cluster
	// service client atomic. This is because the nodes mapping may be updated
	// in parallel.
	rwLock sync.RWMutex

	Nodes  map[uint64]string
	NodeId uint64
}

func NewClusterServiceClient(nodeId uint64) *ClusterServiceClient {
	return &ClusterServiceClient{
		NodeId: nodeId,
	}
}

func (c *ClusterServiceClient) UpdateNodes(nodes map[uint64]string) {
	c.rwLock.Lock()
	c.Nodes = nodes
	c.rwLock.Unlock()
}

func (c *ClusterServiceClient) SendAddFile(destId uint64, userId string, parentDirId uuid.UUID, newFileId uuid.UUID, filename string, contents []byte) error {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	log.Println("Dialing ", c.Nodes[destId], " ", port)
	c.rwLock.RLock()
	conn, err := grpc.Dial(c.Nodes[destId]+port, opts...)
	c.rwLock.RUnlock()

	if err != nil {
		return err
	}
	defer conn.Close()

	request := hootpb.AddNewFileCSRequest{
		NodeKey:     "",
		UserId:      userId,
		ParentDirId: &protos.UUID{Value: parentDirId[:]},
		NewFileId:   &protos.UUID{Value: newFileId[:]},
		NewFileName: filename,
		Contents:    contents,
	}

	client := hootpb.NewClusterServiceClient(conn)
	resp, err := client.AddNewFileCS(context.Background(), &request)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(resp.GetCreatedFileId().Value, newFileId[:]) {
		return ErrIncorrectFileCreated
	}

	return ErrUnimplemented
}

func (c *ClusterServiceClient) SendMakeDirectory(destId uint64, userId string, parentDirId uuid.UUID, newDirId uuid.UUID, dirname string) error {
	var opts []grpc.DialOption

	c.rwLock.RLock()
	conn, err := grpc.Dial(c.Nodes[destId]+port, opts...)
	c.rwLock.RUnlock()

	if err != nil {
		return err
	}
	defer conn.Close()

	request := hootpb.MakeDirectoryCSRequest{
		NodeKey:     "",
		UserId:      userId,
		ParentDirId: &protos.UUID{Value: parentDirId[:]},
		NewDirId:    &protos.UUID{Value: newDirId[:]},
		NewDirName:  dirname,
	}

	client := hootpb.NewClusterServiceClient(conn)
	resp, err := client.MakeDirectoryCS(context.Background(), &request)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(resp.GetCreatedDirId().Value, newDirId[:]) {
		return ErrIncorrectFileCreated
	}

	return ErrUnimplemented
}

func (c *ClusterServiceClient) SendUpdateFileContentsRequest(destId uint64, userId string, fileId uuid.UUID, contents []byte) error {
	var opts []grpc.DialOption

	c.rwLock.RLock()
	conn, err := grpc.Dial(c.Nodes[destId]+port, opts...)
	c.rwLock.RUnlock()
	if err != nil {
		return err
	}
	defer conn.Close()

	request := hootpb.UpdateFileContentsCSRequest{
		NodeKey:  "",
		UserId:   userId,
		FileId:   &protos.UUID{Value: fileId[:]},
		Contents: contents,
	}

	client := hootpb.NewClusterServiceClient(conn)
	resp, err := client.UpdateFileContentsCS(context.Background(), &request)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(resp.UpdatedFileId.Value, fileId[:]) {
		return ErrIncorrectFileCreated
	}

	return ErrUnimplemented
}

func (c *ClusterServiceClient) SendMoveObject(destId uint64, userId string, currentObjId uuid.UUID, newParentId uuid.UUID, newName string) error {
	var opts []grpc.DialOption

	c.rwLock.RLock()
	conn, err := grpc.Dial(c.Nodes[destId]+port, opts...)
	c.rwLock.RUnlock()

	if err != nil {
		return err
	}
	defer conn.Close()

	request := hootpb.MoveObjectCSRequest{
		NodeKey:      "",
		UserId:       userId,
		CurrObjectId: &protos.UUID{Value: currentObjId[:]},
		DirId:        &protos.UUID{Value: newParentId[:]},
		NewName:      newName,
	}

	client := hootpb.NewClusterServiceClient(conn)
	resp, err := client.MoveObjectCS(context.Background(), &request)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(resp.MovedObjectId.Value, currentObjId[:]) {
		return ErrIncorrectFileCreated
	}

	return ErrUnimplemented
}

func (c *ClusterServiceClient) SendRemoveObject(destId uint64, userId string, objId uuid.UUID) error {
	var opts []grpc.DialOption

	c.rwLock.RLock()
	conn, err := grpc.Dial(c.Nodes[destId]+port, opts...)
	c.rwLock.RUnlock()

	if err != nil {
		return err
	}
	defer conn.Close()

	request := hootpb.RemoveObjectCSRequest{
		NodeKey:  "",
		UserId:   userId,
		ObjectId: &protos.UUID{Value: objId[:]},
	}

	client := hootpb.NewClusterServiceClient(conn)
	resp, err := client.RemoveObjectCS(context.Background(), &request)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(resp.RemovedObjectId.Value, objId[:]) {
		return ErrIncorrectFileCreated
	}

	return ErrUnimplemented
}
