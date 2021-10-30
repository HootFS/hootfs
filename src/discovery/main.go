package main

import "github.com/hootfs/hootfs/src/discoverServer/discover"

func main() {
	discoverServer := discover.NewDiscoverServer()

	discoverServer.StartServer()
}
