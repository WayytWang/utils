package wire_inject

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

const injectTag = "@inject"
const set = "set"

const typeInjectFunc = 1
const typeInjectType = 2

var providerPrefix = []string{"Init", "New"}

type Option func(*opt)

type opt struct {
	genPath    string // 文件生成的目标路径
	pkg        string // 生成文件的包名
	searchPath string // gomod项目的绝对路径(入口)
}

// 每一个.go文件有一个fileSearcher
type fileSearcher struct {
	file        string // 文件名
	pkg         string // 文件包名 package ***
	pkgImport   string // 别的文件导入本包时的导入路径
	srcSearcher *searcher
}

type searcher struct {
	sets      []string // set集合
	files     []*fileSearcher
	genPath   string
	pkg       string                        // 生成文件的包名
	modBase   string                        // gomod名 modBase + 文件名 = import
	injectMap map[string]map[string]element // 每一个setName -> 包别名.函数 -> 实体
}

type inject struct {
	doc  string // 注释文本
	name string // 函数/type name
	typ  int    // provider类型
}

// write原料
type element struct {
	setName    string
	pkg        string // 包名
	pkgImport  string // 外部引入路径
	provider   string // 构造函数
	injectName string
	implement  string // 实现的
}

// setName ->
var injectMap map[string]map[string]element

func newGenOpt(genPath string, opts ...Option) *opt {
	o := &opt{genPath: genPath}
	for _, opt := range opts {
		opt(o)
	}
	o.fix()
	return o
}

func (o *opt) fix() {
	if len(o.pkg) == 0 {
		var err error
		// 取o.genPath的绝对路径
		o.pkg, err = filepath.Abs(o.genPath)
		if err != nil {
			o.pkg = "inject"
		} else {
			o.pkg = path.Base(filepath.ToSlash(o.pkg))
		}
	}
	if len(o.searchPath) == 0 {
		modPath := GetGoModDir()
		if len(modPath) > 0 {
			o.searchPath = modPath
		}
	}
}

func (sc *searcher) SearchAllPath(file string) (err error) {
	dir, err := ioutil.ReadDir(file)
	if err != nil {
		return
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
	return
}

func (f *fileSearcher) SearchInject() (err error) {
	data, err := ioutil.ReadFile(f.file)
	if err != nil {
		return err
	}
	if !bytes.Contains(data, []byte(injectTag)) {
		return
	}
	fSet := token.NewFileSet()
	fNode, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		return
	}
	var injects []inject
	ast.Inspect(fNode, func(node ast.Node) bool {
		f.pkg = fNode.Name.Name // 本文件的package name
		// type
		if typeFunc, ok := node.(*ast.GenDecl); ok && strings.Contains(typeFunc.Doc.Text(), injectTag) && typeFunc.Tok == token.TYPE {
			for _, spec := range typeFunc.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				injects = append(injects, inject{
					doc:  typeFunc.Doc.Text(),
					name: typeSpec.Name.Name,
					typ:  typeInjectType,
				})
			}
		}
		// func
		if funcDecl, ok := node.(*ast.FuncDecl); ok && strings.Contains(funcDecl.Doc.Text(), injectTag) {
			funcDecl.Doc.Text()
			injects = append(injects, inject{
				doc:  funcDecl.Doc.Text(),
				name: funcDecl.Name.Name,
				typ:  typeInjectFunc,
			})
		}
		return true
	})
	for _, inj := range injects {
		inj.analysis(fNode, f)
	}
	return
}

// 计算其他文件导入本文件的导入路径 modbase + 文件名
func (f *fileSearcher) getPkgImport() (pkgImport string) {
	abs, err := filepath.Abs(f.file)
	if err != nil {
		panic(err)
	}
	dirName := GetGoModDir()
	if len(abs) < len(dirName) {
		panic("文件名错误")
	}
	pkgImport = filepath.ToSlash(filepath.Dir(filepath.Join(f.srcSearcher.modBase, abs[len(dirName):])))
	return
}

func (inj *inject) analysis(fNode *ast.File, f *fileSearcher) {
	// 解析语法
	lines := strings.Split(inj.doc, "\n")
	for _, comment := range lines {
		//@inject(interface,set=setName) -> set=setName
		comment = strings.TrimSpace(comment)
		if !strings.HasPrefix(comment, injectTag) {
			continue
		}
		comment = comment[len(injectTag):]
		if !(len(comment) >= 2 && comment[0] == '(' && comment[len(comment)-1] == ')') {
			continue
		}
		comment = comment[1 : len(comment)-1]
		implement := ""
		sp := strings.Split(comment, ",")
		if len(sp) > 1 {
			implement = strings.TrimSpace(sp[0])
			if len(implement) == 0 {
				panic("interface/setName can not empty")
			}
			comment = strings.TrimSpace(sp[1])
		}

		if !(strings.HasPrefix(comment, set) && strings.Contains(comment, "=")) {
			return
		}
		spo := strings.Split(comment, "=")
		var setName string // 注入的set
		if len(spo) > 1 {
			setName = spo[1]
		}

		setName = strcase.ToLowerCamel(setName)

		e := element{
			provider:   "",
			setName:    setName,
			implement:  implement,
			injectName: inj.name,
		}

		if inj.typ == typeInjectFunc {
			e.provider = inj.name
		}
		if inj.typ == typeInjectType {
			// 找构造函数
			for _, prefix := range providerPrefix {
				if ct, ok := fNode.Scope.Objects[prefix+inj.name]; ok && ct.Kind == ast.Fun {
					e.provider = prefix + inj.name
					break
				}
			}
		}
		f.analysisInject(e)
	}
}

func (f *fileSearcher) analysisInject(e element) {
	// 将set对应的函数写入f.src.injectMap
	if _, ok := f.srcSearcher.injectMap[e.setName]; !ok {
		f.srcSearcher.injectMap[e.setName] = make(map[string]element)
	}

	pkgImport := f.getPkgImport()
	e.pkgImport = pkgImport
	e.pkg = f.pkg
	f.srcSearcher.injectMap[e.setName][fmt.Sprintf("%s.%s", f.pkg, e.injectName)] = e
}
