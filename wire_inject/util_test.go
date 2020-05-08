package wire_inject

import "testing"

func TestGetGoModDir(t *testing.T) {
	modPath := GetGoModDir()
	t.Log(modPath)
}

func TestGetGoModFilePath(t *testing.T) {
	modPath := GetGoModFilePath()
	t.Log(modPath)
}

func TestGetModBase(t *testing.T) {
	modBase, err := GetModBase()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(modBase)
}