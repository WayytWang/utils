package wire_inject

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

var tmp = map[string]*searcher{}

// todo:日志记录

// run wire命令
func RunWire(genPath string) {
	o := newGenOpt(genPath)
	file := o.searchPath
	pkg := o.pkg

	err := SearchAllPath(file, genPath, pkg)
	if err != nil {
		panic(err)
	}

	log.Printf("write inject files success")
	log.Printf("start runnning wire")
	p, e := exec.LookPath("wire")
	if e != nil {
		panic(fmt.Errorf("wire not found: %v \n%s\n", e,
			"please install wire by [ go get github.com/google/wire/cmd/wire ]"))
	}
	cmd := exec.Command(p)
	var s bytes.Buffer
	cmd.Dir = genPath
	cmd.Stderr = &s
	err = cmd.Run()
	if err != nil {
		log.Printf("[gen failed] %s", s.String())
		panic(err)
	}
	log.Printf("[gen success] %s", s.String())
}

// 目前只考虑结构体注入 provider -> set
//@inject(set=xxx)

func SearchAllPath(file string, genPath string, pkg string) (err error) {
	dir, err := ioutil.ReadDir(file)
	if err != nil {
		return
	}

	modBase, err := GetModBase()
	if err != nil {
		return
	}
	sc := &searcher{
		files:     make([]*fileSearcher, 10),
		genPath:   genPath, // 生成文件地址
		pkg:       pkg,     // 生成文件包名
		modBase:   modBase, // 项目的gomod名
		injectMap: make(map[string]map[string]element),
	}

	for _, f := range dir {
		fn := f.Name()
		if f.IsDir() {
			// todo:有些文件不需要解析
			err = sc.SearchAllPath(filepath.Join(file, fn))
			if err != nil {
				return
			}
		} else if ln := len(fn); ln > 3 && fn[ln-3:] == ".go" {
			fileSc := &fileSearcher{
				srcSearcher: sc,
				file:        filepath.Join(file, fn),
			}
			sc.files = append(sc.files, fileSc)
			err = fileSc.SearchInject()
			if err != nil {
				return
			}
		}
	}
	return writeGen(sc)
}

func writeGen(sc *searcher) (err error) {
	return sc.write()
}
