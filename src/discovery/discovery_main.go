package main

import "github.com/hootfs/hootfs/src/discovery/discover"

func main() {
	discoverServer := discover.NewDiscoverServer()

	discoverServer.StartServer()
}
