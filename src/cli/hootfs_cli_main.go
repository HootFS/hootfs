package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func executeCommand(args []string) {
	// TODO add
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("hootfs> ")
		sc.Scan()
		cmd := sc.Text()
		executeCommand(strings.Fields(cmd))
	}
}
