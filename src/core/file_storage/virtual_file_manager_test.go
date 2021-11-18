package hootfs

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreateNewNamespaceWorks(t *testing.T) {
	vfm := NewVirtualFileManager()
	ns_id, err := vfm.CreateNewNamespace("namespace")
	if err != nil {
		t.Errorf("Failed to create new namespace: %v", err)
	}

	_, exists := vfm.Directories[ns_id]
	if !exists {
		t.Errorf("No namespace root directory added to virtual file system")
	}
}

func TestCreateNewFileWorks(t *testing.T) {
	vfm := NewVirtualFileManager()
	ns_id, err := vfm.CreateNewNamespace("ns")
	if err != nil {
		t.Fatalf("Failed to create new namespace: %v", err)
	}

	file_id, err := vfm.CreateNewFile("filename", ns_id)
	if err != nil {
		t.Errorf("Failed to create new file: %v", err)
	}
	virfile, exists := vfm.Files[file_id]
	if !exists {
		t.Errorf("Created file does not exist in the virtual file manager.")
	}

	if virfile.Id != file_id || virfile.Name != "filename" {
		t.Errorf("File fields do not match.")
	}
}

func TestCreateNewFileFailsIfParentDirectoryD(t *testing.T) {
	vfm := NewVirtualFileManager()

	if _, err := vfm.CreateNewFile("filename", uuid.Nil); err != ErrParentDirDNE {
		t.Errorf("Expected to get '%v', got '%v'", ErrParentDirDNE, err)
	}

}

func TestCreateNewDirectoryWorks(t *testing.T) {
	vfm := NewVirtualFileManager()
	ns_id, err := vfm.CreateNewNamespace("ns")
	if err != nil {
		t.Fatalf("Failed to create new namespace: %v", err)
	}

	dir_id, err := vfm.CreateNewDirectory("dirname", ns_id)
	if err != nil {
		t.Errorf("Failed to create new directory: %v", err)
	}
	virdir, exists := vfm.Directories[dir_id]
	if !exists {
		t.Errorf("Created file does not exist in the virtual file manager.")
	}

	if virdir.Id != dir_id || virdir.Name != "dirname" {
		t.Errorf("File fields do not match.")
	}
}

func TestCreateNewDirectoryFailsIfParentDirectoryDne(t *testing.T) {
	vfm := NewVirtualFileManager()
	if _, err := vfm.CreateNewDirectory("dirname", uuid.Nil); err != ErrParentDirDNE {
		t.Errorf("Expected to get '%v', got '%v'", ErrParentDirDNE, err)
	}
}

func TestAddNewFileWorks(t *testing.T) {
	vfm := NewVirtualFileManager()
	ns_id, err := vfm.CreateNewNamespace("ns")
	if err != nil {
		t.Fatalf("Failed to create new namespace: %v", err)
	}

	file_id := uuid.MustParse(strings.Repeat("1", 32))

	vf := VirtualFile{Name: "filename", Id: file_id}
	if err := vfm.AddNewFile(&vf, ns_id); err != nil {
		t.Errorf("Error adding new file to virtual filesystem: %v", err)
	}

	if vfm.Files[file_id] != vf {
		t.Errorf("File in Virtual File Manager does not match added file.")
	}

	if val, exists := vfm.Directories[ns_id].Files[file_id]; exists == false || val != true {
		t.Errorf("File does not exist in parent directory.")
	}
}

func TestAddNewFileFailsIfParentDirDNE(t *testing.T) {
	vfm := NewVirtualFileManager()

	file_id := uuid.MustParse(strings.Repeat("1", 32))

	vf := VirtualFile{Name: "filename", Id: file_id}
	if err := vfm.AddNewFile(&vf, uuid.Nil); errors.Is(err, ErrDirNotFound(uuid.Nil)) {
		t.Errorf("Expected error %v, got error %v instead", ErrDirNotFound(uuid.Nil), err)
	}
}

func TestAddNewDirectoryWorks(t *testing.T) {
	vfm := NewVirtualFileManager()
	ns_id, err := vfm.CreateNewNamespace("ns")
	if err != nil {
		t.Fatalf("Failed to create new namespace: %v", err)
	}

	dir_id := uuid.MustParse(strings.Repeat("1", 32))

	vd := makeVirtualDirectory("dirname", dir_id)
	if err := vfm.AddNewDirectory(vd, ns_id); err != nil {
		t.Errorf("Error adding new file to virtual filesystem: %v", err)
	}

	if !reflect.DeepEqual(vfm.Directories[dir_id], *vd) {
		t.Errorf("File in Virtual File Manager does not match added file.")
	}

	if val, exists := vfm.Directories[ns_id].Subdirs[dir_id]; exists == false || val != true {
		t.Errorf("File does not exist in parent directory.")
	}
}

func TestAddNewDirectoryFailsIfParentDirDNE(t *testing.T) {
	vfm := NewVirtualFileManager()

	dir_id := uuid.MustParse(strings.Repeat("1", 32))

	vd := makeVirtualDirectory("dir", dir_id)
	if err := vfm.AddNewDirectory(vd, uuid.Nil); errors.Is(err, ErrDirNotFound(uuid.Nil)) {
		t.Errorf("Expected error %v, got error %v instead", ErrDirNotFound(uuid.Nil), err)
	}
}
