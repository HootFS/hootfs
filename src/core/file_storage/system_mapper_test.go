package hootfs

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestAddSystemWorks(t *testing.T) {
	mapper := NewSystemMapper()
	mapper.AddSystem(1)

	if present, exists := mapper.systems[1]; !exists || !present {
		t.Fatalf("System 1 should exist in the mapper, but it does not.")
	}
}

func TestRemoveSystemWorks(t *testing.T) {
	mapper := NewSystemMapper()
	mapper.AddSystem(1)
	mapper.RemoveSystem(1)

	if present, exists := mapper.systems[1]; exists && present {
		t.Fatalf("System 1 should not exists in the system mapper, but it does.")
	}
}

func TestAddToSystemWorks(t *testing.T) {
	mapper := NewSystemMapper()
	mapper.AddSystem(1)
	mapper.AddSystem(2)

	file_id := uuid.MustParse(strings.Repeat("1", 32))
	systems, err := mapper.MapNewFile(file_id)
	if err != nil {
		t.Fatalf("Failure mapping new system: %v", err)
	}

	if len(systems) == 0 {
		t.Errorf("Did not map file to any systems.")
	}
}

func TestAddToSystemFailsIfFileAlreadyExistsInSystem(t *testing.T) {
	mapper := NewSystemMapper()
	mapper.AddSystem(1)

	file_id := uuid.MustParse(strings.Repeat("1", 32))
	mapper.MapNewFile(file_id)
	_, err := mapper.MapNewFile(file_id)
	if err != errFileAlreadyExistsInMapping {
		t.Errorf("Expected to get an error saying that file already exists.")
	}
}

func TestAddToSystemFailsIfNoSystemsInMapper(t *testing.T) {
	mapper := NewSystemMapper()

	file_id := uuid.MustParse(strings.Repeat("1", 32))
	_, err := mapper.MapNewFile(file_id)
	if err != errNoSystemsToMapTo {
		t.Errorf("Should not be able to map to system if no systems exist")
	}
}
