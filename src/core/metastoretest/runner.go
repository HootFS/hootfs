package main

import (
	"fmt"
	"log"

	"github.com/hootfs/hootfs/src/core/vfm"
)

type Person struct {
	Name     string
	Siblings []string
}

func main() {
	ms, err := vfm.NewMetaStore("Test1")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := ms.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	err = ms.DeleteMachine(130)

	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("Success!")
	}
}
