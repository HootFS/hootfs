package discover

import (
	"context"

	protos "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The discover server keeps track of who currently is in the cluster.
// This is done in two ways.
// 1) The discover server can be sent an intentional join or leave request from a machine.
// 2) The discover server will periodically ping all "active" nodes to make sure none have crashed.

type DiscoverServer struct {
       
}

func  (ds *DiscoverServer) JoinCluster(context.Context, *protos.JoinClusterRequest) (*protos.JoinClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}

func (ds *DiscoverServer) LeaveCluster(context.Context, *protos.LeaveClusterRequest) (*protos.LeaveClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveCluster not implemented")
}
