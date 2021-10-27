package discover

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	protos "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	discoverPort   = ":50053"
	clusterMaxSize = 1000
	pingDurr       = 5
)

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

	// This holds equal data to the two maps above.
	// However, it is in a form which easily transmittable
	// via grpc.
	protoActiveMap []*protos.NodeInfo

	rwLock sync.RWMutex

	// Embedded unimplemented service.
	protos.UnimplementedDiscoverServiceServer
}

func NewDiscoverServer() *DiscoverServer {
	return &DiscoverServer{
		nodeMap: make(map[string]*NodeCell),
		idSet:   make(map[uint64]bool),
		// Other values should be zeroed.
	}
}

func getIp(addr string) string {
	return addr[:strings.IndexByte(addr, ':')]
}

func (d *DiscoverServer) StartServer() {
	// Start up the discover server!
	lis, err := net.Listen("tcp", discoverPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	protos.RegisterDiscoverServiceServer(s, d)

	// Start ping checks.
	go func() {
		for {
			time.Sleep(pingDurr * time.Second)
			d.DropIdleAndReset()
		}
	}()

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
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
	d.protoActiveMap = make([]*protos.NodeInfo, len(d.nodeMap))

	i := 0
	for ip, nc := range d.nodeMap {
		d.protoActiveMap[i] = &protos.NodeInfo{
			NodeId: nc.id,
			NodeIp: ip,
		}
		i++
	}
}

func (d *DiscoverServer) DropIdleAndReset() {
	// This function will drop all nodes which have not pinged.
	// And set all remaining nodes ping to false.
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	idleFound := false

	log.Println("Filtering Idle nodes.")

	for ip, nc := range d.nodeMap {
		// If a ping hasn't been received.
		if !nc.ping {
			log.Printf("Node %d (%s) found idle.", nc.id, ip)
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
	ip := getIp(p.Addr.String())

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
		_, found = d.idSet[d.nextId]
	}

	nodeId := d.nextId
	d.nextId++

	// We now have a new Id for our new node.
	d.nodeMap[ip] = &NodeCell{id: nodeId, ping: true}
	d.idSet[nodeId] = true

	d.GenSlices()

	return &protos.JoinClusterResponse{
		NewId:      nodeId,
		ClusterMap: d.protoActiveMap,
	}, nil
}

func (d *DiscoverServer) GetActive(ctx context.Context, gar *protos.GetActiveRequest) (*protos.GetActiveResponse, error) {

	// TODO : verify gar.NodeKey

	// Only read lock is needed to get actives...
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()

	return &protos.GetActiveResponse{
		ClusterMap: d.protoActiveMap,
	}, nil
}

func (d *DiscoverServer) Ping(ctx context.Context, pr *protos.PingRequest) (*protos.PingResponse, error) {

	// TODO : verify pr.NodeKey

	// Retrieve Peer IP address.
	p, _ := peer.FromContext(ctx)
	ip := getIp(p.Addr.String())

	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	if _, found := d.nodeMap[ip]; !found {
		return nil, status.Errorf(codes.NotFound, "Given Node doesn't seemt to be in cluster!")
	}

	d.nodeMap[ip].ping = true

	return &protos.PingResponse{}, nil
}
