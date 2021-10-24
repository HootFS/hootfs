package cluster

import (
	"context"
	"errors"
	"reflect"

	"github.com/google/uuid"
	"github.com/hootfs/hootfs/protos"
	hootpb "github.com/hootfs/hootfs/protos"

	"google.golang.org/grpc"
)

var ErrIncorrectFileCreated = errors.New("The wrong file ID was returned")

type ClusterServiceClient struct {
	nodes  map[int]string
	nodeId int
}

func (c *ClusterServiceClient) SendAddFile(destId int, userId string, parentDirId uuid.UUID, newFileId uuid.UUID, filename string, contents []byte) error {
	var opts []grpc.DialOption

	conn, err := grpc.Dial(c.nodes[destId]+port, opts...)
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

func (c *ClusterServiceClient) SendMakeDirectory(destId int, userId string, parentDirId uuid.UUID, newDirId uuid.UUID, dirname string) error {
	var opts []grpc.DialOption

	conn, err := grpc.Dial(c.nodes[destId]+port, opts...)
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

func (c *ClusterServiceClient) SendUpdateFileContentsRequest(destId int, userId string, fileId uuid.UUID, contents []byte) error {
	var opts []grpc.DialOption

	conn, err := grpc.Dial(c.nodes[destId]+port, opts...)
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

func (c *ClusterServiceClient) SendMoveObject(destId int, userId string, currentObjId uuid.UUID, newParentId uuid.UUID, newName string) error {
	var opts []grpc.DialOption

	conn, err := grpc.Dial(c.nodes[destId]+port, opts...)
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

func (c *ClusterServiceClient) SendRemoveObject(destId int, userId string, objId uuid.UUID) error {
	var opts []grpc.DialOption

	conn, err := grpc.Dial(c.nodes[destId]+port, opts...)
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
