package hootfs

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewVirtualFileManger(t *testing.T) {
	vfm := NewVirtualFileManager()

	if vfm.Directories == nil || vfm.Files == nil {
		t.Fatal("Failed in creating a virtual file manager!")
	}
}

func TestCreateNewFile(t *testing.T) {
	vfm := NewVirtualFileManager()

	fid1 := uuid.MustParse(strings.Repeat("1", 32))
	fn1 := "test file 1"

	fid2 := uuid.MustParse(strings.Repeat("2", 32))
	fn2 := "test file 2"
	vfm.CreateNewFile("testFile1")
}
