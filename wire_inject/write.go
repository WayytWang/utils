package wire_inject

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

// 生成文件
func (sc *searcher) write() (err error) {
	sc.sets = make([]string, 0)
	err = os.MkdirAll(sc.genPath, 0775)
	if err != nil {
		panic(err)
	}
	_ = sc.clearOldWireGen()

	for set, m := range sc.injectMap {
		_ = sc.writeSet(set, m)
	}
	return sc.writeSets()
}

func (sc *searcher) writeSets() (err error) {
	sort.Strings(sc.sets)
	fileName := filepath.Join(sc.genPath, "inject_sets.go")
	src := fmt.Sprintf(injectSetNameTemplate, sc.pkg, "sets", strings.Join(sc.sets, ",\n\t"))
	res, err := format.Source([]byte(src))
	if err != nil {
		fmt.Print(src)
		return err
	}
	err = ioutil.WriteFile(fileName, res, 0664)
	return
}

func (sc *searcher) writeSet(set string, m map[string]element) (err error) {
	var providers []string
	for _, e := range m {
		if e.implement == "" {
			providers = append(providers, fmt.Sprintf("%s.%s", e.pkg, e.provider))
		} else {
			itfName := e.implement
			if !strings.Contains(itfName, ".") {
				itfName = e.pkg + "." + itfName
			}
			injectName := e.pkg + "." + e.injectName
			provider := e.provider
			if provider == "" {
				provider = fmt.Sprintf(`wire.Struct(new(%s), "*")`,injectName)
			}else{
				provider = e.pkg + "." + e.provider
			}
			providers = append(providers, fmt.Sprintf("%s,\n\t wire.Bind(new(%s), new(*%s))", provider, itfName, injectName))
		}
		// todo import优化
	}
	fileName := filepath.Join(sc.genPath, "inject_"+strcase.ToSnake(set)+".go")
	setName := strings.Title(strcase.ToCamel(set)) + "Set"
	sc.sets = append(sc.sets, setName)
	src := fmt.Sprintf(injectSetNameTemplate, sc.pkg, setName, strings.Join(providers, ",\n\n\n\t"))
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fs, f)
	if err != nil {
		return
	}
	formatSrc, err := format.Source(buffer.Bytes())
	if err != nil {
		return
	}
	formatSrc, err = imports.Process("", formatSrc, nil)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fileName, formatSrc, 0664)
	if err != nil {
		return
	}
	return
}

// 移除旧的wire_gen
func (sc *searcher) clearOldWireGen() (err error) {
	dirs, err := ioutil.ReadDir(sc.genPath)
	if err != nil {
		return
	}
	if len(dirs) == 0 {
		return
	}
	_ = os.Remove("wire_gen.go")
	for _, f := range dirs {
		if strings.Contains(f.Name(), "inject_") && strings.Contains(f.Name(), ".go") {
			_ = os.Remove(filepath.Join(sc.genPath, f.Name()))
		}
	}
	return
}
