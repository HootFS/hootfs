package hootfs

import "testing"

func TestAddToSystemWorks(t *testing.T) {
	mapper := NewSystemMapper()
	mapper.AddSystem(1)

	if present, ex := mapper.systems[1]; !ex || !present {
		t.Fatalf("System 1 should exist in the mapper, but it does not.")
	}
}

// TODO (josh8551021): Finish testing
