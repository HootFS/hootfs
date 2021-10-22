package core

import (
	"context"

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
