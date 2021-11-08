package hootfs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
)

type fileSystem interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadFile(name string) ([]byte, error)
	Mkdir(name string, perm os.FileMode) error
	Remove(name string) error
	RemoveAll(name string) error
}

type osFS struct{}

func (osFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (osFS) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (osFS) Remove(name string) error {
	return os.Remove(name)
}

func (osFS) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

var ErrUnimplemented = errors.New("Method unimplemented")
var ErrObjectNotFound = errors.New("Object not found")
var ErrFileNotFound = errors.New("File not found")
var ErrNeedFileNotDir = errors.New("Expected to get a file; got a directory")
var ErrNeedDirNotFile = errors.New("Expected to get a directory; got a file")

func ErrMismatchedNamespace(expected string, found string) error {
	return fmt.Errorf("Namespace %s did not match expected namespace %s", found, expected)
}

func ErrDirNotFound(directory uuid.UUID) error {
	return fmt.Errorf("Directory with ID %s not found", directory.String())
}

func ErrDuplicateIDFound(filename string, dirname string) error {
	return fmt.Errorf("Objects have same ID (%s) (%s)", filename, dirname)
}

type FileType int

const (
	FILE FileType = iota + 1
	DIRECTORY
)

type FileInfo struct {
	NamespaceId string
	ObjectId    uuid.UUID
}

type FileObject struct {
	Namespace        string   // What namespace does this file belong to?
	ParentDir        string   // Parent directory relative to namespace root
	RelativeFilename string   // Filename relative to parentDir
	Filetype         FileType // File or directory?
}

type VirtualFileMapper struct {
	files  map[uuid.UUID]FileObject
	rwLock sync.RWMutex
}

type FileManager struct {
	Root string
	Fs   fileSystem
	Vfm  VirtualFileMapper
}

func NewFileSystemManager(root string) *FileManager {
	manager := new(FileManager)
	manager.Root = root
	manager.Fs = osFS{}
	manager.Vfm = VirtualFileMapper{files: make(map[uuid.UUID]FileObject)}

	return manager
}

func CreateNewSystemDirectory(namespace string, name string) *FileObject {
	return &FileObject{Namespace: namespace, ParentDir: "", RelativeFilename: name, Filetype: DIRECTORY}
}

func CreateNewSystemFile(namespace string, name string) *FileObject {
	return &FileObject{Namespace: namespace, ParentDir: "", RelativeFilename: name, Filetype: FILE}
}

// Lazy file creation
// File will only be written to system once contents are available.
func (m *FileManager) CreateFile(filename string, fileInfo *FileInfo) {
	m.Vfm.rwLock.Lock()
	m.Vfm.files[fileInfo.ObjectId] = *CreateNewSystemFile(fileInfo.NamespaceId, filename)
	m.Vfm.rwLock.Unlock()
}

func (m *FileManager) WriteFile(fileInfo *FileInfo, contents []byte) error {
	m.Vfm.rwLock.RLock()
	fileObj, exists := m.Vfm.files[fileInfo.ObjectId]
	m.Vfm.rwLock.RUnlock()

	if !exists {
		return ErrFileNotFound
	}

	if fileObj.Namespace != fileInfo.NamespaceId {
		return ErrMismatchedNamespace(fileInfo.NamespaceId, fileObj.Namespace)
	}

	if fileObj.Filetype == DIRECTORY {
		return ErrNeedFileNotDir
	}

	return m.Fs.WriteFile(path.Join(m.Root, fileObj.Namespace, fileObj.ParentDir, fileObj.RelativeFilename), contents, 666)
}

func (m *FileManager) ReadFile(fileInfo *FileInfo) ([]byte, error) {
	m.Vfm.rwLock.RLock()
	fileObj, exists := m.Vfm.files[fileInfo.ObjectId]
	m.Vfm.rwLock.RUnlock()
	if !exists {
		return nil, ErrFileNotFound
	}

	if fileObj.Namespace != fileInfo.NamespaceId {
		return nil, ErrMismatchedNamespace(fileInfo.NamespaceId, fileObj.Namespace)
	}

	if fileObj.Filetype == DIRECTORY {
		return nil, ErrNeedFileNotDir
	}

	data, err := m.Fs.ReadFile(path.Join(m.Root, fileInfo.NamespaceId, fileObj.ParentDir, fileObj.RelativeFilename))
	if err != nil {
		// Here, we attempt to read a file that might not yet exist on disk.
		// However, due to the earlier check, we can be reasonably sure that the
		// file exsits in the virtual file system. So instead of propagating the
		// error, we simply return an empty byte slice.
		return make([]byte, 0), nil
	}

	return data, nil
}

func (m *FileManager) DeleteFile(fileInfo *FileInfo) error {
	m.Vfm.rwLock.Lock()
	defer m.Vfm.rwLock.Unlock()

	fileObj, exists := m.Vfm.files[fileInfo.ObjectId]

	if !exists {
		return ErrFileNotFound
	}

	err := m.Fs.Remove(path.Join(m.Root, fileObj.ParentDir, fileObj.RelativeFilename))
	if err != nil {
		return err
	}
	delete(m.Vfm.files, fileInfo.ObjectId)

	return nil
}

func (m *FileManager) CreateDirectory(directory_name string, fileInfo *FileInfo) error {
	err := m.Fs.Mkdir(path.Join(m.Root, fileInfo.NamespaceId, directory_name), 775)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}

	m.Vfm.rwLock.Lock()
	m.Vfm.files[fileInfo.ObjectId] = *CreateNewSystemDirectory(fileInfo.NamespaceId, directory_name)
	m.Vfm.rwLock.Unlock()
	return nil
}

func (m *FileManager) DeleteDirectory(fileInfo *FileInfo) error {
	m.Vfm.rwLock.Lock()
	defer m.Vfm.rwLock.Unlock()

	fileObj, exists := m.Vfm.files[fileInfo.ObjectId]
	if !exists {
		return ErrFileNotFound
	}

	if fileObj.Filetype == DIRECTORY {
		return ErrNeedDirNotFile
	}

	err := m.Fs.RemoveAll(path.Join(m.Root, fileObj.ParentDir, fileObj.RelativeFilename))
	if err != nil {
		return err
	}
	delete(m.Vfm.files, fileInfo.ObjectId)
	return nil
}
