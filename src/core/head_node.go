package hootfs

import (
	"context"

	head "github.com/hootfs/hootfs/protos"
)

type fileManagerServer struct {
}

func (s *fileManagerServer) GetDirectoryContents(
	ctx context.Context, request *head.GetDirectoryContentsRequest) (*head.GetDirectoryContentsResponse, error) {
}
