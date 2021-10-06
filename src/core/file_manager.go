package hootfs

import "errors"

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

type FileType int

const (
	FILE FileType = iota + 1
	DIRECTORY
)

type FileObject struct {
	relative_filename string
	filetype          FileType
}

type LocalFileManager struct {
}

var ErrUnimplemented = errors.New("Method unimplemented")

func (manager LocalFileManager) WriteFile(filename string, contents []byte) error {
	return ErrUnimplemented

}

func (manager LocalFileManager) ReadFile(filename string) ([]byte, error) {
	return nil, ErrUnimplemented
}

//
func (manager LocalFileManager) DeleteFile(filename string) error {
	return ErrUnimplemented
}

func (manager LocalFileManager) CreateDirectory(directory_name string) error {
	return ErrUnimplemented
}

func (manager LocalFileManager) DeleteDirectory(directory_name string) error {
	return ErrUnimplemented
}

func (manager LocalFileManager) GetDirectoryContents(directory_name string) ([]FileObject, error) {
	return nil, ErrUnimplemented
}
