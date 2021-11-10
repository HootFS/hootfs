package hootfs

import (
	"fmt"

	"sync"

	"github.com/google/uuid"
)

type VirtualDirectory struct {
	Name    string
	Id      uuid.UUID
	Subdirs map[uuid.UUID]bool
	Files   map[uuid.UUID]bool
}

type VirtualFile struct {
	Name string
	Id   uuid.UUID
}

type VirtualFileManager struct {
	// The virtual file manager has little information on the namespace.
	Directories map[uuid.UUID]VirtualDirectory
	Files       map[uuid.UUID]VirtualFile

	RWLock sync.RWMutex
}

func NewVirtualFileManager() *VirtualFileManager {
	return &VirtualFileManager{
		Directories: make(map[uuid.UUID]VirtualDirectory),
		Files:       make(map[uuid.UUID]VirtualFile),
	}
}

func (m *VirtualFileManager) CreateNewFile(filename string, parent uuid.UUID) (uuid.UUID, error) {
	fileUUID, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to get new UUID for virtual file: %v", err)
	}

	m.RWLock.Lock()
	m.Directories[parent].Files[fileUUID] = true
	m.Files[fileUUID] = VirtualFile{Name: filename, Id: fileUUID}
	m.RWLock.Unlock()

	return fileUUID, nil
}

func (m *VirtualFileManager) CreateNewDirectory(dirname string, parent uuid.UUID) (uuid.UUID, error) {
	dirUUID, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to get new UUID for virtual file: %v", err)
	}

	m.RWLock.Lock()
	m.Directories[parent].Files[dirUUID] = true
	m.Directories[dirUUID] = VirtualDirectory{
		Name:    dirname,
		Id:      dirUUID,
		Subdirs: make(map[uuid.UUID]bool),
		Files:   make(map[uuid.UUID]bool)}
	m.RWLock.Unlock()

	return dirUUID, nil
}

func (m *VirtualFileManager) AddNewFile(file *VirtualFile, parent uuid.UUID) error {
	dir, exists := m.Directories[parent]
	if !exists {
		return ErrDirNotFound(parent)
	}

	m.RWLock.Lock()
	dir.Files[file.Id] = true
	m.Files[file.Id] = *file
	m.RWLock.Unlock()

	return nil
}

func (m *VirtualFileManager) AddNewDirectory(dir *VirtualDirectory, parent uuid.UUID) error {
	par_dir, exists := m.Directories[parent]
	if !exists {
		return ErrDirNotFound(parent)
	}

	m.RWLock.Lock()
	par_dir.Subdirs[dir.Id] = true
	m.Directories[dir.Id] = *dir
	m.RWLock.Unlock()

	return nil
}

func (m *VirtualFileManager) MoveObject(file uuid.UUID, oldParent uuid.UUID, newParent uuid.UUID) error {
	m.RWLock.Lock()
	defer m.RWLock.Unlock()

	oldDir, oldExists := m.Directories[oldParent]
	if !oldExists {
		return ErrDirNotFound(oldParent)
	}
	newDir, newExists := m.Directories[newParent]
	if !newExists {
		return ErrDirNotFound(newParent)
	}

	_, fileExists := oldDir.Files[file]
	_, dirExists := oldDir.Subdirs[file]

	// If both a directory and file exist with the same UUID
	// there has been an internal error. Report it.
	if fileExists && dirExists {
		filename := m.Directories[file].Name
		dirname := m.Files[file].Name
		return ErrDuplicateIDFound(filename, dirname)
	}

	if !fileExists && !dirExists {
		return ErrObjectNotFound
	}

	if fileExists {
		delete(oldDir.Files, file)
		newDir.Files[file] = true
	} else {
		delete(oldDir.Subdirs, file)
		newDir.Subdirs[file] = true
	}

	return nil
}

func (m *VirtualFileManager) RemoveObject(file uuid.UUID, parent uuid.UUID) error {
	m.RWLock.Lock()
	defer m.RWLock.Unlock()

	dir, exists := m.Directories[parent]
	if !exists {
		return ErrDirNotFound(parent)
	}

	_, fileExists := dir.Files[file]
	_, dirExists := dir.Subdirs[file]

	if fileExists && dirExists {
		filename := m.Directories[file].Name
		dirname := m.Files[file].Name
		return ErrDuplicateIDFound(filename, dirname)
	}

	if !fileExists && !dirExists {
		return ErrObjectNotFound
	}

	if fileExists {
		delete(dir.Files, file)
	} else {
		delete(dir.Subdirs, file)
	}

	return nil
}
