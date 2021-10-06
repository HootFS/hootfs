package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func executeCommand(args []string) {
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "read":
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: read [filename]")
		}
	}
}

func main() {
	var servername string
	flag.StringVar(&servername, "s", "" /* todo set reasonable default */, "the uri of the head node to connect to")
	flag.Parse()

	// todo connect to head node

	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("hootfs> ")
		sc.Scan()
		cmd := sc.Text()
		executeCommand(strings.Fields(cmd))
	}
}
