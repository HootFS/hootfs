package hootfs

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

var ErrFileNotFound = errors.New("File not found")
var ErrNeedFileNotDir = errors.New("Expected to get a file; got a directory")
var ErrNeedDirNotFile = errors.New("Expected to get a directory; got a file")

func ErrMismatchedNamespace(expected string, found string) error {
	return fmt.Errorf("Namespace %s did not match expected namespace %s", found, expected)
}

type FileSystemManager interface {
	// File Operations
	CreateFile(filename string) error
	WriteFile(filename string, content []byte) error
	ReadFile(filename string) ([]byte, error)
	DeleteFile(filename string) error
	MoveFile(old_filename string, new_filename string) error

	// Directory Operations
	CreateDirectory(directory_name string) error
	DeleteDriectory(directory_name string) error
	GetDirectoryContents(direcotry_name string) ([]FileObject, error)
	MoveDirectory(old_direcotry_string string, new_directory_string string)
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
	vfm  VirtualFileMapper
}

var ErrUnimplemented = errors.New("Method unimplemented")

func CreateNewSystemDirectory(namespace string, name string) *FileObject {
	return &FileObject{namespace: namespace, parentDir: "", relativeFilename: name, filetype: DIRECTORY}
}

func CreateNewSystemFile(namespace string, name string) *FileObject {
	return &FileObject{namespace: namespace, parentDir: "", relativeFilename: name, filetype: FILE}
}

// Lazy file creation
// File will only be written to system once contents are available.
func (m *FileManager) CreateFile(filename string, fileInfo FileInfo) {
	m.vfm.files[fileInfo.objectId] = *CreateNewSystemFile(fileInfo.namespaceId, filename)
}

func (m *FileManager) WriteFile(fileInfo FileInfo, contents []byte) error {
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

	return os.WriteFile(path.Join(m.root, fileObj.namespace, fileObj.parentDir, fileObj.relativeFilename), contents, 666)
}

func (m FileManager) ReadFile(fileInfo FileInfo) ([]byte, error) {
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

	return os.ReadFile(path.Join(m.root, fileInfo.namespaceId, fileObj.parentDir, fileObj.relativeFilename))
}

//
func (m FileManager) DeleteFile(fileInfo FileInfo) error {
	fileObj, exists := m.vfm.files[fileInfo.objectId]
	if !exists {
		return ErrFileNotFound
	}

	err := os.Remove(path.Join(m.root, fileObj.parentDir, fileObj.relativeFilename))
	if err != nil {
		return err
	}
	delete(m.vfm.files, fileInfo.objectId)
	return nil
}

func (m FileManager) CreateDirectory(directory_name string, fileInfo FileInfo) error {
	err := os.Mkdir(path.Join(m.root, fileInfo.namespaceId, directory_name), 775)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}
	m.vfm.files[fileInfo.objectId] = *CreateNewSystemDirectory(fileInfo.namespaceId, directory_name)
	return ErrUnimplemented
}

func (m FileManager) DeleteDirectory(fileInfo FileInfo) error {
	fileObj, exists := m.vfm.files[fileInfo.objectId]
	if !exists {
		return ErrFileNotFound
	}

	err := os.RemoveAll(path.Join(m.root, fileObj.parentDir, fileObj.relativeFilename))
	if err != nil {
		return err
	}
	delete(m.vfm.files, fileInfo.objectId)
	return nil
}