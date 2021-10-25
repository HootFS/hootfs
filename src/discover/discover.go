package discover

import (
	"context"
	"sync"

	protos "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
) // For now the max cluster size.
const clusterMaxSize = 1000

type NodeCell struct {
    id   uint64 // Node ID.
    ping bool   // Whether or not this node has pinged recently.
}

// The discover server keeps track of who currently is in the cluster.

type DiscoverServer struct {
    // These two maps form a bidirectional map between 
    // IP addresses and node ids.
    // These structures are for fast lookups when pinging
    // and adding new ids to the cluster.
    nodeMap map[string]*NodeCell
    idSet   map[uint64]bool

    // Next ID to attempt to use for a cluster node.
    nextId uint64

    // These slices are updated whenever the above maps 
    // are updated.
    // Specifically, they are what are returned to clients 
    // when an active list is requested.
    idSlice []uint64 // These can start as NIL.
    ipSlice []string

    rwLock  sync.RWMutex    
}

func NewDiscoverServer() (*DiscoverServer) {
    return &DiscoverServer{
        nodeMap: make(map[string]*NodeCell),
        idSet: make(map[uint64]bool),
        // Other values should be zeroed.
    }
}

func (d *DiscoverServer) GenSlices() {
    // This function regenerates the slices inside
    // the Discover server struct.
    // Note, in theory we could just edit the existing slices...
    // but this could lead to concurency errors with grpc.
    // So, this, while slow, is the best option.
    //
    // NOTE : this is not thread safe!!!
    // When using this function, make sure a lock has
    // been aquired.
    d.idSlice = make([]uint64, len(d.nodeMap))
    d.ipSlice = make([]string, len(d.nodeMap))

    i := 0
    for ip, nc := range d.nodeMap {
        d.idSlice[i] = nc.id
        d.ipSlice[i] = ip
        i++
    }
}

func (d *DiscoverServer) DropIdleAndReset() {
    // This function will drop all nodes which have not pinged.
    // And set all remaining nodes ping to false.
    d.rwLock.Lock()
    defer d.rwLock.Unlock()

    idleFound := false

    for ip, nc := range d.nodeMap {
        // If a ping hasn't been received.
        if !nc.ping {
            // Delete the node from the cluster maps.
            delete(d.idSet, nc.id)
            delete(d.nodeMap, ip) 
            idleFound = true
        } else {
            nc.ping = false
        }
    }

    if idleFound {
        // Only regenerate id slices if
        // a node was removed from the cluster.
        d.GenSlices()
    }
}

func (d *DiscoverServer) JoinCluster(ctx context.Context, jcr *protos.JoinClusterRequest) (*protos.JoinClusterResponse, error) {
    // Retrieve Peer IP address. 
    p, _ := peer.FromContext(ctx)
    ip := p.Addr.String()

    // Critical section for the rest of the function.
    d.rwLock.Lock()
    defer d.rwLock.Unlock()

    // Need to check if the given IP is in the cluster already.
    if _, found := d.nodeMap[ip]; found {
        return nil, status.Errorf(codes.AlreadyExists, "This IP address is already in the cluster!")
    }

    // Need to check if there is space for anymore nodes in the cluster. 
    if len(d.nodeMap) >= clusterMaxSize {
        return nil, status.Errorf(codes.ResourceExhausted, "Cluster is already full!")
    }

    // Finally, we can add this new node to the cluster!
    // First get its id.
    _, found := d.idSet[d.nextId]
    for found {
        d.nextId++
        found = d.idSet[d.nextId]
    }

    nodeId := d.nextId
    d.nextId++

    // We now have a new Id for our new node.
    d.nodeMap[ip] = &NodeCell{id: nodeId, ping: true} 
    d.idSet[nodeId] = true

    d.GenSlices()

    return &protos.JoinClusterResponse{
        NewId: nodeId,  
        ClusterIds: d.idSlice,
        ClusterIps: d.ipSlice,
    }, nil
}

func (d *DiscoverServer) GetActive(ctx context.Context, gar *protos.GetActiveRequest) (*protos.GetActiveResponse, error) {

    // TODO : verify gar.NodeKey

    // Only read lock is needed to get actives...
    d.rwLock.RLock()
    defer d.rwLock.RUnlock()

    return &protos.GetActiveResponse{
        ClusterIds: d.idSlice,
        ClusterIps: d.ipSlice,
    }, nil 
}

func (d *DiscoverServer) Ping(ctx context.Context, pr *protos.PingRequest) (*protos.PingResponse, error) {

    // TODO : verify pr.NodeKey

    // Retrieve Peer IP address. 
    p, _ := peer.FromContext(ctx)
    ip := p.Addr.String()

    d.rwLock.Lock()
    defer d.rwLock.Unlock()
    
    if _, found := d.nodeMap[ip]; !found {
        return nil, status.Errorf(codes.NotFound, "Given Node doesn't seemt to be in cluster!") 
    }

    d.nodeMap[ip].ping = true
    
    return &protos.PingResponse{}, nil
}
