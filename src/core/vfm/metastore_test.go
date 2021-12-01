package vfm

import (
	"context"
	"errors"
	"testing"
)

func FatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

var ErrIncorrectError = errors.New("Expected error was not thrown!")

func ExpectErr(t *testing.T, actual error, expected error) {
	if actual != expected {
		t.Fatal(ErrIncorrectError)
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

		ExpectErr(t, ms.CreateMachine(0), ErrMachineExists)
	})

	t.Run("Machine Removal", func(t *testing.T) {
		FatalIfErr(t, ms.CreateMachine(2))
		FatalIfErr(t, ms.DeleteMachine(2))
		FatalIfErr(t, ms.CheckMachine(2, ErrMachineExists, nil))

		ExpectErr(t, ms.DeleteMachine(4), ErrMachineDoesNotExist)
	})

	t.Run("User Addition", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Chatham"))
		FatalIfErr(t, ms.CheckUser("Chatham", nil, ErrUserDoesNotExist))
		FatalIfErr(t, ms.CheckUser("Mark", ErrUserExists, nil))

		ExpectErr(t, ms.CreateUser("Chatham"), ErrUserExists)
	})

	t.Run("User Removal", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Matt"))
		FatalIfErr(t, ms.DeleteUser("Matt"))
		FatalIfErr(t, ms.CheckUser("Matt", ErrUserExists, nil))

		ExpectErr(t, ms.DeleteUser("Matt"), ErrUserDoesNotExist)
	})

	t.Run("Namespace Addition", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Paula"))

		nsid, err := ms.CreateNamespace("NS 1", "Paula")
		FatalIfErr(t, err)

		FatalIfErr(t, ms.CheckNamespace(nsid, "Paula", nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Paul", ErrAccess, nil))
	})

	t.Run("Namespace Removal", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Bob"))

		nsid, err := ms.CreateNamespace("NS 1", "Bob")
		FatalIfErr(t, err)

		ExpectErr(t, ms.DeleteNamespace(nsid, "Mark"), ErrNoAccess)

		FatalIfErr(t, ms.DeleteNamespace(nsid, "Bob"))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Bob", ErrAccess, nil))

		ExpectErr(t, ms.DeleteNamespace(nsid, "Bob"), ErrNoAccess)
	})

	t.Run("Namespace Permission Addition", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Mike"))
		FatalIfErr(t, ms.CreateUser("Karl"))

		nsid, err := ms.CreateNamespace("NS 1", "Mike")
		FatalIfErr(t, err)

		ExpectErr(t, ms.AddUserToNamespace(nsid, "Karl", "Karl"), ErrNoAccess)
		ExpectErr(t, ms.AddUserToNamespace(nsid, "Mike", "Trent"),
			ErrUserDoesNotExist)

		FatalIfErr(t, ms.AddUserToNamespace(nsid, "Mike", "Karl"))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Mike", nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Karl", nil, ErrNoAccess))

		FatalIfErr(t, ms.CreateUser("Jack"))
		FatalIfErr(t, ms.AddUserToNamespace(nsid, "Karl", "Jack"))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Jack", nil, ErrNoAccess))
	})

	t.Run("Namespace Permission Removal", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Josh"))
		FatalIfErr(t, ms.CreateUser("Reid"))

		nsid, err := ms.CreateNamespace("NS 1", "Josh")
		FatalIfErr(t, err)

		ExpectErr(t, ms.RemoveUserFromNamespace(nsid, "Josh", "Reid"),
			ErrNoUserInNamespace)

		FatalIfErr(t, ms.AddUserToNamespace(nsid, "Josh", "Reid"))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Reid", nil,
			ErrNoUserInNamespace))

		FatalIfErr(t, ms.RemoveUserFromNamespace(nsid, "Reid", "Josh"))
		FatalIfErr(t, ms.CheckNamespace(nsid, "Josh", ErrAccess, nil))
	})

	t.Run("User Removal With Namespaces", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Eddy"))
		FatalIfErr(t, ms.CreateUser("Joe"))
		FatalIfErr(t, ms.CreateUser("Bobby"))

		nsid1, err := ms.CreateNamespace("NS 1", "Eddy")
		FatalIfErr(t, err)
		FatalIfErr(t, ms.AddUserToNamespace(nsid1, "Eddy", "Joe"))

		nsid2, err := ms.CreateNamespace("NS 2", "Joe")
		FatalIfErr(t, err)
		FatalIfErr(t, ms.AddUserToNamespace(nsid2, "Joe", "Bobby"))

		FatalIfErr(t, ms.DeleteUser("Joe"))

		FatalIfErr(t, ms.CheckNamespace(nsid1, "Joe", ErrAccess, nil))
		FatalIfErr(t, ms.CheckNamespace(nsid2, "Joe", ErrAccess, nil))

		FatalIfErr(t, ms.CheckNamespace(nsid1, "Eddy", nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckNamespace(nsid2, "Bobby", nil, ErrNoAccess))
	})

	FatalIfErr(t, ms.db.Drop(context.TODO()))
	FatalIfErr(t, ms.Disconnect())
}
