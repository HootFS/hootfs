package hootfs

import "errors"

type FileManager interface {
	CreateFile(filename string) error
	WriteFile(filename string, content []byte) error
	ReadFile(filename string) ([]byte, error)
	DeleteFile(filename string) error
	CreateDirectory(directory_name string) error
	DeleteDriectory(directory_name string) error
	GetDirectoryContents(direcotry_name string) ([]FileObject, error)
}

type FileType int

const (
	File FileType = iota + 1
	Directory
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
