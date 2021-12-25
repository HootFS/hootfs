package main

import (
	"flag"
	"log"

	hootfs "github.com/hootfs/hootfs/src/core/file_storage"
	core "github.com/hootfs/hootfs/src/core/hootfs"
)

func main() {
	var discover_hostname string
	var hootfs_root string

	flag.StringVar(&discover_hostname, "dip", "127.0.0.1", "Address of the discovery server to connect to.")
	flag.Parse()

	flag.StringVar(&hootfs_root, "root", "./resources/root", "Root directory for the file system.")
	server := core.NewHootFsServer(
		discover_hostname,
		hootfs.NewFileSystemManager(hootfs_root),
		hootfs.NewVirtualFileManager())

	if err := server.StartServer(); err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}
}
