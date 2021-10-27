package hootfs

import (
	"errors"

	"github.com/google/uuid"
)

type SystemMapper struct {
	systems       map[uint64]bool
	systemMapping map[uuid.UUID][]uint64
}

var errFileAlreadyExistsInMapping = errors.New("File already exists in System mapper")

func NewSystemMapper() *SystemMapper {
	return &SystemMapper{systems: make(map[uint64]bool), systemMapping: make(map[uuid.UUID][]uint64)}
}

func (s *SystemMapper) AddSystem(sysId uint64) {
	s.systems[sysId] = true
}

func (s *SystemMapper) RemoveSystem(sysId uint64) {
	s.systems[sysId] = false
}

func (s *SystemMapper) MapNewFile(fileId uuid.UUID) ([]uint64, error) {
	if _, exists := s.systemMapping[fileId]; exists {
		return nil, errFileAlreadyExistsInMapping
	}

	s.systemMapping[fileId] = make([]uint64, len(s.systems))
	for k := range s.systems {
		s.systemMapping[fileId] = append(s.systemMapping[fileId], k)
	}

	return s.systemMapping[fileId], nil
}

func (s *SystemMapper) RemoveFile(fileId uuid.UUID) error {
	delete(s.systemMapping, fileId)
	return nil
}
