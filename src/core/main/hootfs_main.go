package main

import (
	"flag"
	"log"

	"github.com/google/uuid"
	"github.com/hootfs/hootfs/src/core"
	hootfs "github.com/hootfs/hootfs/src/core/file_storage"
)

func main() {
	var discover_hostname string
	var hootfs_root string

	flag.StringVar(&discover_hostname, "dip", "127.0.0.1", "Address of the discovery server to connect to.")
	flag.Parse()

	flag.StringVar(&hootfs_root, "root", "../../../resources/root", "Root directory for the file system.")
	server := core.NewHootFsServer(
		discover_hostname,
		hootfs.NewFileSystemManager(hootfs_root),
		&hootfs.VirtualFileManager{
			Directories: make(map[uuid.UUID]hootfs.VirtualDirectory),
			Files:       make(map[uuid.UUID]hootfs.VirtualFile),
		})
	if err := server.StartServer(); err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}
}
