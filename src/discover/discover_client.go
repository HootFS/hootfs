package discover

import (
	"context"

	"github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc"
)

type DiscoverClient struct {
    discover_ip string    // Ip address of the discovery server to comnnect to.
}

func toIdMap(protoMap []*protos.NodeInfo) map[uint64]string {
    idMap := make(map[uint64]string)

    for _, ni := range protoMap {
        idMap[ni.NodeId] = ni.NodeIp
    }

    return idMap
}

func (dc *DiscoverClient) JoinCluster() (uint64, map[uint64]string, error) {
    // TODO : add some form of node security here.
    // As of now, no node key is given to this node.

	var opts []grpc.DialOption


	conn, err := grpc.Dial(dc.discover_ip + discover_port, opts...)
	if err != nil {
		return 0, nil, err
	}
	defer conn.Close()


    request := protos.JoinClusterRequest{ClusterKey: ""}
    client := protos.NewDiscoverServiceClient(conn)
    resp, err := client.JoinCluster(context.Background(), &request)

    if err != nil {
        return 0, nil, err
    }

    return resp.NewId, toIdMap(resp.ClusterMap), nil
}


func (dc *DiscoverClient) GetActive() (map[uint64]string, error) {
	var opts []grpc.DialOption


	conn, err := grpc.Dial(dc.discover_ip + discover_port, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
    
    request := protos.GetActiveRequest{NodeKey: ""}
    client := protos.NewDiscoverServiceClient(conn)
    resp, err := client.GetActive(context.Background(), &request)

    if err != nil {
        return nil, err
    }

    return toIdMap(resp.ClusterMap), nil
}

func (dc *DiscoverClient) Ping() error {
	var opts []grpc.DialOption


	conn, err := grpc.Dial(dc.discover_ip + discover_port, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

    request := protos.PingRequest{NodeKey: ""}
    client := protos.NewDiscoverServiceClient(conn)
    _, err = client.Ping(context.Background(), &request)

    return err
}
