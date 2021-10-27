package hootfs

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

// var fs fileSystem = osFS{}

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

var ErrFileNotFound = errors.New("File not found")
var ErrNeedFileNotDir = errors.New("Expected to get a file; got a directory")
var ErrNeedDirNotFile = errors.New("Expected to get a directory; got a file")

func ErrMismatchedNamespace(expected string, found string) error {
	return fmt.Errorf("Namespace %s did not match expected namespace %s", found, expected)
}

type FileType int

const (
	FILE FileType = iota + 1
	DIRECTORY
)

type FileInfo struct {
	namespaceId string
	objectId    uuid.UUID
}

type FileObject struct {
	namespace        string   // What namespace does this file belong to?
	parentDir        string   // Parent directory relative to namespace root
	relativeFilename string   // Filename relative to parentDir
	filetype         FileType // File or directory?
}

type VirtualFileMapper struct {
	files map[uuid.UUID]FileObject
}

type FileManager struct {
	root string
	fs   fileSystem
	vfm  VirtualFileMapper
}

var ErrUnimplemented = errors.New("Method unimplemented")

func NewFileSystemManager(root string) *FileManager {
	manager := new(FileManager)
	manager.root = root
	manager.fs = osFS{}
	manager.vfm = VirtualFileMapper{files: make(map[uuid.UUID]FileObject)}

	return manager
}

func CreateNewSystemDirectory(namespace string, name string) *FileObject {
	return &FileObject{namespace: namespace, parentDir: "", relativeFilename: name, filetype: DIRECTORY}
}

func CreateNewSystemFile(namespace string, name string) *FileObject {
	return &FileObject{namespace: namespace, parentDir: "", relativeFilename: name, filetype: FILE}
}

// Lazy file creation
// File will only be written to system once contents are available.
func (m *FileManager) CreateFile(filename string, fileInfo *FileInfo) {
	m.vfm.files[fileInfo.objectId] = *CreateNewSystemFile(fileInfo.namespaceId, filename)
}

func (m *FileManager) WriteFile(fileInfo *FileInfo, contents []byte) error {
	fileObj, exists := m.vfm.files[fileInfo.objectId]
	if !exists {
		return ErrFileNotFound
	}

	if fileObj.namespace != fileInfo.namespaceId {
		return ErrMismatchedNamespace(fileInfo.namespaceId, fileObj.namespace)
	}

	if fileObj.filetype == DIRECTORY {
		return ErrNeedFileNotDir
	}

	return m.fs.WriteFile(path.Join(m.root, fileObj.namespace, fileObj.parentDir, fileObj.relativeFilename), contents, 666)
}

func (m *FileManager) ReadFile(fileInfo *FileInfo) ([]byte, error) {
	fileObj, exists := m.vfm.files[fileInfo.objectId]
	if !exists {
		return nil, ErrFileNotFound
	}

	if fileObj.namespace != fileInfo.namespaceId {
		return nil, ErrMismatchedNamespace(fileInfo.namespaceId, fileObj.namespace)
	}

	if fileObj.filetype == DIRECTORY {
		return nil, ErrNeedFileNotDir
	}

	data, err := m.fs.ReadFile(path.Join(m.root, fileInfo.namespaceId, fileObj.parentDir, fileObj.relativeFilename))
	if err != nil {
		return make([]byte, 0), nil
	}

	return data, nil
}

//
func (m *FileManager) DeleteFile(fileInfo *FileInfo) error {
	fileObj, exists := m.vfm.files[fileInfo.objectId]
	if !exists {
		return ErrFileNotFound
	}

	err := m.fs.Remove(path.Join(m.root, fileObj.parentDir, fileObj.relativeFilename))
	if err != nil {
		return err
	}
	delete(m.vfm.files, fileInfo.objectId)
	return nil
}

func (m *FileManager) createDirectory(directory_name string, fileInfo *FileInfo) error {
	err := m.fs.Mkdir(path.Join(m.root, fileInfo.namespaceId, directory_name), 775)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}
	m.vfm.files[fileInfo.objectId] = *CreateNewSystemDirectory(fileInfo.namespaceId, directory_name)
	return nil
}

func (m *FileManager) deleteDirectory(fileInfo *FileInfo) error {
	fileObj, exists := m.vfm.files[fileInfo.objectId]
	if !exists {
		return ErrFileNotFound
	}

	err := m.fs.RemoveAll(path.Join(m.root, fileObj.parentDir, fileObj.relativeFilename))
	if err != nil {
		return err
	}
	delete(m.vfm.files, fileInfo.objectId)
	return nil
}
