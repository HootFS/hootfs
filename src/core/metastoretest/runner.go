package main

import (
	"log"

	"github.com/hootfs/hootfs/src/core/vfm"
)

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

	err = ms.CreateUser("Chatham")

	if err != nil {
		log.Fatal(err)
	}

	nsid, err := ms.CreateNamespace("NS 3", "Mark")

	if err != nil {
		log.Fatal(err)
	}

	err = ms.DeleteNamespace(nsid, "Chatham")

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Success!")
}
