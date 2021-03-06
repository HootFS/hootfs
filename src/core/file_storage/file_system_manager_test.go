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
	manager.CreateFile(test_file, &FileInfo{NamespaceId: test_ns, ObjectId: file_id})

	file_info := manager.Vfm.files[file_id]
	if file_info.Namespace != test_ns {
		t.Fatalf("Namespaces do not match: %v, %v", file_info.Namespace, test_ns)
	}
	if file_info.Filetype != FILE {
		t.Fatalf("Expected file, got directory")
	}
	if file_info.RelativeFilename != test_file {
		t.Fatalf("Filenames do not match")
	}
}

func TestWriteFileWorksIfFileExists(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"
	test_file := "test_file"
	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.CreateFile(test_file, &file_info)

	manager.Fs = fakeFS{}
	if err := manager.WriteFile(&file_info, make([]byte, 1)); err != nil {
		t.Fatalf("Failed to write to supposedly existing file; %v", err)
	}
}

func TestWriteFileFailsIfFileDNE(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}

	if err := manager.WriteFile(&file_info, make([]byte, 1)); err == nil {
		t.Fatalf("Managed to write to non-existing file")
	}
}

func TestReadFileWorksIfFileExists(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"
	test_file := "test_file"
	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.CreateFile(test_file, &file_info)

	manager.Fs = fakeFS{}
	if _, err := manager.ReadFile(&file_info); err != nil {
		t.Fatalf("Failed to write to supposedly existing file; %v", err)
	}
}

func TestReadFileFailsIfFileDNE(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}

	_, err := manager.ReadFile(&file_info)
	if err != ErrFileNotFound {
		t.Fatalf("Managed to find non-existing file")
	}
}

func TestReadFileFailsIfFiletypeIsDirectory(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	manager.CreateDirectory("dir", &file_info)

	_, err := manager.ReadFile(&file_info)
	if err != ErrNeedFileNotDir {
		t.Fatalf("Managed to read file contents of non-file")
	}
}

func TestDeleteFileWorks(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	manager.CreateDirectory("dir", &file_info)

	err := manager.DeleteFile(&file_info)
	if err != ErrNeedFileNotDir {
		t.Fatalf("Managed to read file contents of non-file")
	}
}

func TestDeleteFileFailsIfFileDne(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	manager.CreateFile("file", &file_info)

	err := manager.DeleteFile(&file_info)
	if err != nil {
		t.Fatalf("Failed to delete directory: %v", err)
	}
}

func TestDeleteFileFailsIfCalledOnDirectory(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	manager.CreateDirectory("dir", &file_info)

	err := manager.DeleteFile(&file_info)
	if err != ErrNeedFileNotDir {
		t.Fatalf("Managed to perform file deletion on directory")
	}
}

func TestCreateDirectoryWorks(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	if err := manager.CreateDirectory("dir", &file_info); err != nil {
		t.Fatalf("Error creating directory")
	}
}

func TestDeleteDirectoryWorks(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	manager.CreateDirectory("dir", &file_info)

	err := manager.DeleteDirectory(&file_info)
	if err != nil {
		t.Fatalf("Failed to delete virtual directory: %v", err)
	}
}

func TestDeleteDirectoryFailsIfDirectoryNotFound(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}

	if err := manager.DeleteDirectory(&file_info); err != ErrFileNotFound {
		t.Fatalf("Expected to receive a file-not-found error for non-existent directory.")
	}
}

func TestDeleteDirectoryFailsIfCalledOnFile(t *testing.T) {
	manager := NewFileSystemManager("fs_root")
	file_id := uuid.MustParse(strings.Repeat("1", 32))
	test_ns := "test_namespace"

	file_info := FileInfo{NamespaceId: test_ns, ObjectId: file_id}
	manager.Fs = fakeFS{}
	manager.CreateFile("file", &file_info)

	if err := manager.DeleteDirectory(&file_info); err != ErrNeedDirNotFile {
		t.Fatalf("Performed directory deletion on file")
	}
}
