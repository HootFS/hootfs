package hootfs

import (
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type fakeFS struct{}

func (fakeFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return nil
}

func (fakeFS) ReadFile(name string) ([]byte, error) {
	return nil, nil
}

func (fakeFS) Mkdir(name string, perm os.FileMode) error {
	return nil
}

func (fakeFS) Remove(name string) error {
	return nil
}

func (fakeFS) RemoveAll(name string) error {
	return nil
}

func TestCreateFileWorks(t *testing.T) {
	manager := NewFileSystemManager("root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"
	test_file := "test_file"
	manager.CreateFile(test_file, &FileInfo{namespaceId: test_ns, objectId: file_id})

	file_info := manager.vfm.files[file_id]
	if file_info.namespace != test_ns {
		t.Fatalf("Namespaces do not match: %v, %v", file_info.namespace, test_ns)
	}
	if file_info.filetype != FILE {
		t.Fatalf("Expected file, got directory")
	}
	if file_info.relativeFilename != test_file {
		t.Fatalf("Filenames do not match")
	}
}

func TestWriteFileWorksIfFileExists(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"
	test_file := "test_file"
	file_info := FileInfo{namespaceId: test_ns, objectId: file_id}
	manager.CreateFile(test_file, &file_info)

	manager.fs = fakeFS{}
	if err := manager.WriteFile(&file_info, make([]byte, 1)); err != nil {
		t.Fatalf("Failed to write to supposedly existing file; %v", err)
	}
}

func TestWriteFileFailsIfFileDNE(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{namespaceId: test_ns, objectId: file_id}
	manager.fs = fakeFS{}

	if err := manager.WriteFile(&file_info, make([]byte, 1)); err == nil {
		t.Fatalf("Managed to write to non-existing file")
	}
}

func TestReadFileWorksIfFileExists(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"
	test_file := "test_file"
	file_info := FileInfo{namespaceId: test_ns, objectId: file_id}
	manager.CreateFile(test_file, &file_info)

	manager.fs = fakeFS{}
	if _, err := manager.ReadFile(&file_info); err != nil {
		t.Fatalf("Failed to write to supposedly existing file; %v", err)
	}
}

func TestReadFileFailsIfFileDNE(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{namespaceId: test_ns, objectId: file_id}
	manager.fs = fakeFS{}

	_, err := manager.ReadFile(&file_info)
	if err == nil {
		t.Fatalf("Managed to read non-existing file")
	}

	if err != ErrFileNotFound {
		t.Fatalf("Managed to find non-existing file")
	}
}
