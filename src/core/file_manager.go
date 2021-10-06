package hootfs

import (
	"errors"
	"os"
)

type FileManager interface {
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

// Enum for the the type of file. Either a file or a directory.
type FileType int

const (
	FILE FileType = iota + 1
	DIRECTORY
)

// A FileObject struct holds a file's name and its type.
type FileObject struct {
	relative_filename string
	filetype          FileType
}

// LocalFileManger is a implementation of the file manager interface for storing files locally
type LocalFileManager struct {
	source_directory string
}

var ErrUnimplemented = errors.New("Method unimplemented")

// Write contents to a file
func (manager LocalFileManager) WriteFile(filename string, contents []byte) error {
	return os.WriteFile(filename, contents, 0666)
}

func (manager LocalFileManager) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

//
func (manager LocalFileManager) DeleteFile(filename string) error {
	return os.Remove(filename)
}

func (manager LocalFileManager) CreateDirectory(directory_name string) error {
	return os.Mkdir(directory_name, 0755)
}

func (manager LocalFileManager) DeleteDirectory(directory_name string) error {
	return os.Remove(directory_name)
}

func (manager LocalFileManager) GetDirectoryContents(directory_name string) ([]FileObject, error) {
	contents, err := os.ReadDir(directory_name)

	files := make([]FileObject, len(contents))

	for i, object := range contents {
		files[i] = FileObject{relative_filename: object.Name()}
		if object.IsDir() {
			files[i].filetype = DIRECTORY
		} else {
			files[i].filetype = FILE
		}
	}
	return files, err
}
