package hootfs

import (
	"fmt"

	"sync"

	"github.com/google/uuid"
)

func ErrDirNotFound(directory uuid.UUID) error {
	return fmt.Errorf("Directory with ID %s not found", directory.String())
}

type VirtualDirectory struct {
	name    string
	id      uuid.UUID
	subdirs map[uuid.UUID]bool
	files   map[uuid.UUID]bool
}

type VirtualFile struct {
	Name string
	Id   uuid.UUID
}

type VirtualFileManager struct {
	directories map[uuid.UUID]VirtualDirectory
	files       map[uuid.UUID]VirtualFile

	rwLock sync.RWMutex
}

func (m *VirtualFileManager) CreateNewFile(filename string, parent uuid.UUID) (uuid.UUID, error) {
	fileUUID, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to get new UUID for virtual file: %v", err)
	}

	m.rwLock.Lock()
	m.directories[parent].files[fileUUID] = true
	m.files[fileUUID] = VirtualFile{Name: filename, Id: fileUUID}
	m.rwLock.Unlock()

	return fileUUID, nil
}

func (m *VirtualFileManager) CreateNewDirectory(dirname string, parent uuid.UUID) (uuid.UUID, error) {
	dirUUID, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to get new UUID for virtual file: %v", err)
	}

	m.rwLock.Lock()
	m.directories[parent].files[dirUUID] = true
	m.directories[dirUUID] = VirtualDirectory{
		name:    dirname,
		id:      dirUUID,
		subdirs: make(map[uuid.UUID]bool),
		files:   make(map[uuid.UUID]bool)}
	m.rwLock.Unlock()

	return dirUUID, nil
}

func (m *VirtualFileManager) AddNewFile(file VirtualFile, parent uuid.UUID) error {
	dir, exists := m.directories[parent]
	if !exists {
		return ErrDirNotFound(parent)
	}

	m.rwLock.Lock()
	dir.files[file.Id] = true
	m.files[file.Id] = file
	m.rwLock.Unlock()

	return nil
}

func (m *VirtualFileManager) AddNewDirectory(dir VirtualDirectory, parent uuid.UUID) error {
	par_dir, exists := m.directories[parent]
	if !exists {
		return ErrDirNotFound(parent)
	}

	m.rwLock.Lock()
	par_dir.subdirs[dir.id] = true
	m.directories[dir.id] = dir
	m.rwLock.Unlock()

	return nil
}

func (m *VirtualFileManager) MoveObject(file uuid.UUID, oldParent uuid.UUID, newParent uuid.UUID) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	oldDir, oldExists := m.directories[oldParent]
	if !oldExists {
		return ErrDirNotFound(oldParent)
	}
	newDir, newExists := m.directories[newParent]
	if !newExists {
		return ErrDirNotFound(newParent)
	}

	_, fileExists := oldDir.files[file]
	_, dirExists := oldDir.subdirs[file]
	if !fileExists && !dirExists {
		return ErrFileNotFound
	} else if fileExists {
		delete(oldDir.files, file)
		newDir.files[file] = true
	} else {
		delete(oldDir.subdirs, file)
		newDir.subdirs[file] = true
	}

	return nil
}

func (m *VirtualFileManager) RemoveObject(file uuid.UUID, parent uuid.UUID) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	dir, exists := m.directories[parent]
	if !exists {
		return ErrDirNotFound(parent)
	}

	_, fileExists := dir.files[file]
	_, dirExists := dir.subdirs[file]
	if !fileExists && !dirExists {
		// If the file doesn't exist, we ignore
	} else if fileExists {
		delete(dir.files, file)
	} else {
		delete(dir.subdirs, file)
	}

	return ErrUnimplemented
}
