package vfm

import (
	"context"
	"testing"
)

func FatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// We have this test all function to ensure clean up occurs at
// the end of the test
func TestAll(t *testing.T) {
	ms, err := NewMetaStore("TESTING-DB")

	if err != nil {
		t.Fatal(err)
	}

	t.Run("Machine Addition", func(t *testing.T) {
		FatalIfErr(t, ms.CreateMachine(0))
		FatalIfErr(t, ms.CheckMachine(0, nil, ErrMachineDoesNotExist))
		FatalIfErr(t, ms.CheckMachine(1, ErrMachineExists, nil))
	})

	t.Run("Machine Removal", func(t *testing.T) {
		FatalIfErr(t, ms.CreateMachine(2))
		FatalIfErr(t, ms.DeleteMachine(2))
		FatalIfErr(t, ms.CheckMachine(2, ErrMachineExists, nil))
	})

	t.Run("User Addition", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Chatham"))
		FatalIfErr(t, ms.CheckUser("Chatham", nil, ErrUserDoesNotExist))
		FatalIfErr(t, ms.CheckUser("Mark", ErrUserExists, nil))
	})

	t.Run("User Removal", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Matt"))
		FatalIfErr(t, ms.DeleteUser("Matt"))
		FatalIfErr(t, ms.CheckUser("Matt", ErrUserExists, nil))
	})

	t.Run("Namespace Addition", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Paula"))

		nsid, err := ms.CreateNamespace("NS 1", "Paula")
		FatalIfErr(t, err)

		FatalIfErr(t, ms.CheckNamespace(nsid, "Paula", nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Paul", ErrAccess, nil))
	})

	err = ms.db.Drop(context.TODO())

	if err != nil {
		t.Fatal(err)
	}

	err = ms.Disconnect()

	if err != nil {
		t.Fatal(err)
	}
}
