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
		t.Logf("Actual Error : %s\n", actual.Error())
		t.Fatal(ErrIncorrectError)
	}
}

var ErrConditionNotMet = errors.New("Condition not met!")

func ExpectTrue(t *testing.T, res bool) {
	if !res {
		t.Log("Condition Not Met!")
		t.Fatal(ErrConditionNotMet)
	}
}

// We have this test all function to ensure clean up occurs at
// the end of the test
func TestAll(t *testing.T) {
	ms, err := NewMetaStore("TESTING-DB")
	FatalIfErr(t, err)

	// Drop in the beginning is better.
	FatalIfErr(t, ms.db.Drop(context.TODO()))

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

	t.Run("File Locations", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Dany"))
		nsid, err := ms.CreateNamespace("NS 1", "Dany")
		FatalIfErr(t, err)

		void, err := ms.CreateFreeObjectInNamespace(nsid, "Dany", "File 1",
			VFM_File_Type)
		FatalIfErr(t, err)

		ExpectErr(t, ms.SetFileLocations(void, []Machine_ID{1, 2}),
			ErrMachineDoesNotExist)

		FatalIfErr(t, ms.CreateMachine(3))
		FatalIfErr(t, ms.CreateMachine(4))
		FatalIfErr(t, ms.SetFileLocations(void, []Machine_ID{3, 4}))
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
		FatalIfErr(t, ms.CreateUser("Lola"))

		nsid1, err := ms.CreateNamespace("NS 1", "Bob")
		FatalIfErr(t, err)

		FatalIfErr(t, ms.AddUserToNamespace(nsid1, "Bob", "Lola"))

		nsid2, err := ms.CreateNamespace("NS 2", "Lola")
		FatalIfErr(t, err)

		void1, err := ms.CreateFreeObjectInNamespace(nsid1, "Bob", "File 1",
			VFM_File_Type)
		FatalIfErr(t, err)

		FatalIfErr(t, ms.AddObjectToNamespace(nsid2, "Lola", void1))

		ExpectErr(t, ms.DeleteNamespace(nsid1, "Mark"), ErrNoAccess)

		FatalIfErr(t, ms.DeleteNamespace(nsid1, "Bob"))
		FatalIfErr(t, ms.CheckNamespace(nsid1, "Bob", ErrAccess, nil))

		ExpectErr(t, ms.DeleteNamespace(nsid1, "Bob"), ErrNoAccess)

		FatalIfErr(t, ms.CheckVObjectAccess(void1, "Bob", ErrAccess, nil))
		FatalIfErr(t, ms.CheckVObjectAccess(void1, "Lola", nil, ErrNoAccess))
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

	t.Run("Root Object Creation", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Iago"))
		nsid, err := ms.CreateNamespace("NS 1", "Iago")
		FatalIfErr(t, err)

		void, err := ms.CreateFreeObjectInNamespace(nsid, "Iago", "Folder 1",
			VFM_Dir_Type)

		FatalIfErr(t, err)

		FatalIfErr(t, ms.CheckVObject(void, nil, ErrVObjectNotFound))
		FatalIfErr(t, ms.CheckRoot(nsid, void, nil, ErrVObjectNotFound))
		FatalIfErr(t, ms.CheckVObjectAccess(void, "Iago",
			nil, ErrVObjectNotFound))
	})

	t.Run("Object Creation", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Dez"))
		FatalIfErr(t, ms.CreateUser("Othello"))

		nsid, err := ms.CreateNamespace("NS 1", "Dez")
		FatalIfErr(t, err)

		void1, err := ms.CreateFreeObjectInNamespace(nsid, "Dez", "Folder 1",
			VFM_Dir_Type)
		FatalIfErr(t, err)

		void2, err := ms.CreateObject(void1, "Dez", "File 1", VFM_File_Type)
		FatalIfErr(t, err)

		FatalIfErr(t, ms.CheckVObject(void2, nil, ErrVObjectNotFound))
		FatalIfErr(t, ms.CheckVObjectAccess(void2, "Dez", nil, ErrNoAccess))

		_, err = ms.CreateObject(void1, "Othello", "File 2", VFM_File_Type)
		ExpectErr(t, err, ErrNoAccess)

		_, err = ms.CreateObject(void2, "Dez", "File 2", VFM_File_Type)
		ExpectErr(t, err, ErrNotADirectory)
	})

	t.Run("Simple Object Deletion", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Marco"))

		nsid, err := ms.CreateNamespace("NS 1", "Marco")
		FatalIfErr(t, err)

		void, err := ms.CreateFreeObjectInNamespace(nsid, "Marco", "Folder 1",
			VFM_Dir_Type)
		FatalIfErr(t, err)

		ExpectErr(t, ms.DeleteObject(void, "Billy"), ErrNoAccess)
		FatalIfErr(t, ms.DeleteObject(void, "Marco"))

		FatalIfErr(t, ms.CheckVObject(void, ErrVObjectFound, nil))
		FatalIfErr(t, ms.CheckRoot(nsid, void, ErrVObjectFound, nil))
	})

	t.Run("Addition To Namespace", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("David"))
		FatalIfErr(t, ms.CreateUser("Paul"))

		nsid1, err := ms.CreateNamespace("NS 1", "David")
		FatalIfErr(t, err)

		void1, err := ms.CreateFreeObjectInNamespace(nsid1, "David", "Folder 1",
			VFM_Dir_Type)
		FatalIfErr(t, err)

		ExpectErr(t, ms.AddObjectToNamespace(nsid1, "David", void1),
			ErrObjectInNamespace)

		void2, err := ms.CreateObject(void1, "David", "Folder 2", VFM_Dir_Type)
		FatalIfErr(t, err)

		void3, err := ms.CreateObject(void2, "David", "File 1", VFM_File_Type)
		FatalIfErr(t, err)

		void4, err := ms.CreateObject(void2, "David", "File 2", VFM_File_Type)
		FatalIfErr(t, err)

		nsid2, err := ms.CreateNamespace("NS 2", "David")
		FatalIfErr(t, err)

		ExpectErr(t, ms.AddObjectToNamespace(nsid2, "Paul", void2), ErrNoAccess)
		FatalIfErr(t, ms.AddObjectToNamespace(nsid2, "David", void2))

		FatalIfErr(t, ms.CheckVObjectNamespace(void2, nsid2, nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckVObjectNamespace(void3, nsid2, nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckVObjectNamespace(void4, nsid2, nil, ErrNoAccess))

		FatalIfErr(t, ms.AddObjectToNamespace(nsid2, "David", void1))

		FatalIfErr(t, ms.CheckVObjectNamespace(void2, nsid2, nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckVObjectNamespace(void3, nsid2, nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckVObjectNamespace(void4, nsid2, nil, ErrNoAccess))
	})

	t.Run("Removal From Namespace", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Rico"))

		nsid1, err := ms.CreateNamespace("NS 1", "Rico")
		FatalIfErr(t, err)

		void1, err := ms.CreateFreeObjectInNamespace(nsid1, "Rico",
			"Folder 1", VFM_Dir_Type)
		FatalIfErr(t, err)

		void2, err := ms.CreateObject(void1, "Rico", "Folder 2", VFM_Dir_Type)
		FatalIfErr(t, err)

		void3, err := ms.CreateObject(void2, "Rico", "File 1", VFM_File_Type)
		FatalIfErr(t, err)

		FatalIfErr(t, ms.CreateUser("Brad"))

		nsid2, err := ms.CreateNamespace("NS 2", "Brad")
		FatalIfErr(t, err)

		FatalIfErr(t, ms.AddUserToNamespace(nsid2, "Brad", "Rico"))
		FatalIfErr(t, ms.AddObjectToNamespace(nsid2, "Rico", void2))

		ExpectErr(t, ms.RemoveObjectFromNamespace(nsid2, "Rico", void3),
			ErrNotRoot)
		ExpectErr(t, ms.RemoveObjectFromNamespace(nsid1, "Brad", void1),
			ErrNoAccess)

		err = ms.RemoveObjectFromNamespace(nsid2, "Rico", void2)
		FatalIfErr(t, err)

		FatalIfErr(t, ms.CheckVObjectAccess(void2, "Brad", ErrAccess, nil))
		FatalIfErr(t, ms.CheckVObjectAccess(void3, "Brad", ErrAccess, nil))
		FatalIfErr(t, ms.CheckVObjectAccess(void2, "Rico", nil, ErrNoAccess))
		FatalIfErr(t, ms.CheckVObjectAccess(void3, "Rico", nil, ErrNoAccess))
	})

	t.Run("Namespace Details", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Joey"))

		nsid, err := ms.CreateNamespace("NS 1", "Joey")
		FatalIfErr(t, err)

		void, err := ms.CreateFreeObjectInNamespace(nsid, "Joey",
			"Folder 1", VFM_Dir_Type)
		FatalIfErr(t, err)

		namespace, err := ms.GetNamespaceDetails(nsid, "Joey")
		ExpectTrue(t, namespace.NSID == nsid)
		ExpectTrue(t, namespace.Name == "NS 1")
		ExpectTrue(t, namespace.RootObjects[0] == void)
		ExpectTrue(t, namespace.Users[0] == "Joey")
	})

	t.Run("Object Details", func(t *testing.T) {
		FatalIfErr(t, ms.CreateUser("Houston"))
		FatalIfErr(t, ms.CreateUser("Dallas"))

		nsid1, err := ms.CreateNamespace("NS 1", "Houston")
		FatalIfErr(t, err)
		FatalIfErr(t, ms.AddUserToNamespace(nsid1, "Houston", "Dallas"))

		nsid2, err := ms.CreateNamespace("NS 2", "Houston")
		FatalIfErr(t, err)

		nsid3, err := ms.CreateNamespace("NS 3", "Dallas")
		FatalIfErr(t, err)

		nsid4, err := ms.CreateNamespace("NS 4", "Houston")
		FatalIfErr(t, err)

		void1, err := ms.CreateFreeObjectInNamespace(nsid1, "Houston",
			"Folder 1", VFM_Dir_Type)
		FatalIfErr(t, err)

		void2, err := ms.CreateObject(void1, "Houston",
			"Folder 2", VFM_Dir_Type)
		FatalIfErr(t, err)

		void3, err := ms.CreateObject(void2, "Houston", "File 1",
			VFM_File_Type)

		FatalIfErr(t, ms.AddObjectToNamespace(nsid2, "Houston", void2))
		FatalIfErr(t, ms.AddObjectToNamespace(nsid3, "Dallas", void2))
		FatalIfErr(t, ms.AddObjectToNamespace(nsid4, "Houston", void3))

		obj1_h, err := ms.GetObjectDetails(void1, "Houston")
		FatalIfErr(t, err)

		ExpectTrue(t, obj1_h.GetID() == void1)
		sobjects1_h, _ := obj1_h.GetSubObjects()
		ExpectTrue(t, sobjects1_h[0].Namespaces[0].NSID == nsid2)
		ExpectTrue(t, len(sobjects1_h[0].Namespaces) == 1)

		obj2_d, err := ms.GetObjectDetails(void2, "Dallas")
		FatalIfErr(t, err)

		ExpectTrue(t, len(obj2_d.GetNamespaces()) == 2)
		sobjects2_d, _ := obj2_d.GetSubObjects()
		ExpectTrue(t, len(sobjects2_d[0].Namespaces) == 0)
	})

	FatalIfErr(t, ms.Disconnect())
}
