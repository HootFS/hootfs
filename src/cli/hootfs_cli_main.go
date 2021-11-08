package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	head "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc"
)

const (
	// After merging, maybe have this resolve from
	// head_node.go ... instead of redefining it here.
	connectingPort = ":50060"
)

func executeCommand(args []string, rpcClient head.HootFsServiceClient) {
	if len(args) == 0 { // blank command
		return
	}
	switch args[0] {
	case "read":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: read [filename]")
			return
		}
		fileid, _ := uuid.Parse(args[1])
		req := head.GetFileContentsRequest{FileId: &head.UUID{Value: fileid[:]}}
		rpcClient.GetFileContents(context.Background(), &req)
	case "write":
		if len(args) != 4 {
			fmt.Fprintln(os.Stderr, "usage: write [dir] [name] [contents]")
			return
		}
		dstuuid, _ := uuid.Parse(args[1])
		contents, err := os.ReadFile(args[3])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		req := head.AddNewFileRequest{Contents: contents, FileName: args[2], DirId: &head.UUID{Value: dstuuid[:]}} // TODO add uuid
		rpcClient.AddNewFile(context.Background(), &req)
	case "update":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "usage: update [file] [contents]")
			return
		}
		dstuuid, _ := uuid.Parse(args[1])
		contents, err := os.ReadFile(args[2])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		req := head.UpdateFileContentsRequest{FileId: &head.UUID{Value: dstuuid[:]}, Contents: contents}
		rpcClient.UpdateFileContents(context.Background(), &req)
	case "move":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "usage: move [dir] [newname] [src]")
			return
		}
		srcuuid, _ := uuid.Parse(args[3])
		dstuuid, _ := uuid.Parse(args[1])
		_ = dstuuid
		req := head.MoveObjectRequest{ObjectId: &head.UUID{Value: srcuuid[:]}, DirId: &head.UUID{Value: dstuuid[:]}, NewName: args[2]}
		rpcClient.MoveObject(context.Background(), &req)
	case "delete":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: delete [object]")
			return
		}
		objid, _ := uuid.Parse(args[1])
		req := head.RemoveObjectRequest{ObjectId: &head.UUID{Value: objid[:]}}
		rpcClient.RemoveObject(context.Background(), &req)
	case "mkdir":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: mkdir [rootdir] [dirname]")
			return
		}
		rootdirid, _ := uuid.Parse(args[1])
		req := head.MakeDirectoryRequest{DirId: &head.UUID{Value: rootdirid[:]}, DirName: args[2]}
		rpcClient.MakeDirectory(context.Background(), &req)
	case "ls":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: ls [dir]")
		}
		dirid, _ := uuid.Parse(args[1])
		req := head.GetDirectoryContentsRequest{DirId: &head.UUID{Value: dirid[:]}}
		resp, _ := rpcClient.GetDirectoryContents(context.Background(), &req)
		respobjects := resp.Objects
		for _, obj := range respobjects {
			respobjectuuid, _ := uuid.FromBytes(obj.ObjectId.Value)
			fmt.Printf("%s (id=%s)\n", obj.ObjectName, respobjectuuid.String())
		}
	default:
		fmt.Fprintf(os.Stderr, "no such command %s\n", args[0])
	}
}

func main() {
	var serverAddr string
	flag.StringVar(&serverAddr, "s", "127.0.0.1", "the address of the head node to connect to")
	flag.Parse()

	// connect to the rpc
	var grpcOpts []grpc.DialOption
	conn, err := grpc.Dial(serverAddr+connectingPort, grpcOpts...)
	if err != nil {
		// failed to connect?
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	// initialize rpc from connection
	rpcClient := head.NewHootFsServiceClient(conn)

	// run a "shell" where commands can be typed
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("hootfs> ")
		sc.Scan()
		cmd := sc.Text()
		executeCommand(strings.Fields(cmd), rpcClient)
	}
}
