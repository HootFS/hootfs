package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	head "github.com/hootfs/hootfs/protos"
	"google.golang.org/grpc"
)

func executeCommand(args []string, rpcClient head.HootFsServiceClient) {
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "read":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: read [filename]")
			return
		}
		req := head.GetFileContentsRequest{} // TODO do a proper initialization of this
		rpcClient.GetFileContents(context.Background(), &req)
	case "write":
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "usage: write [dst] [src]")
			return
		}
		contents, err := os.ReadFile(args[2])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		req := head.AddNewFileRequest{Contents: contents, FileName: args[2]} // TODO add uuid
		rpcClient.AddNewFile(context.Background(), &req)
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
	conn, err := grpc.Dial(serverAddr, grpcOpts...)
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
