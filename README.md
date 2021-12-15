# HootFS

__HootFS__ is a simple distribute file system meant to run on small clusters.

## Usage
```bash
# Start the discovery server.
go run src/discovery/discovery_main.go

# Start a cluster node.
go run src/core/hootfs_main.go -dip <discovery server IP>
```

__Note :__ A discovery server manages the cluster's nodes. 
Thus, the discovery server must be running before nodes can be added to the cluster.
