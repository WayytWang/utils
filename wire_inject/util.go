package wire_inject

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

var modTmp string

// 获取gomod的name
func GetModBase() (modBase string,err error) {
	modPath := GetGoModFilePath()
	mb,_ := ioutil.ReadFile(modPath)
	f,err := modfile.Parse("",mb, func(path, version string) (s string, e error) {
		return version,nil
	})
	if err != nil {
		return
	}
	if f.Module == nil {
		err = errors.New("parse mod error,please check your go env")
		return
	}
	modBase = f.Module.Mod.Path
	return
}

// 获取项目的本地路径
func GetGoModDir() (modPath string) {
	mod := GetGoModFilePath()
	modPath,_ = filepath.Split(mod)
	return
}

func GetGoModFilePath() (modPath string) {
	if len(modTmp) > 0 {
		return modTmp
	}
	cmd := exec.Command("go","env","GOMOD")
	buffer := &bytes.Buffer{}
	cmd.Stdout = buffer
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	mod := buffer.String()
	modPath = strings.Trim(mod,"\n")
	modTmp = modPath
	return
}
