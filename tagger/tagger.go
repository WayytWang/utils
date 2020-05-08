package tagger

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

// 找到目录上所有的go文件
func ParseAllPath(file string, opts ...GenTagOption) (err error) {
	dir, err := ioutil.ReadDir(file)
	if err != nil {
		return
	}
	for _, sub := range dir {
		if sub.IsDir() {
			// 递归解析目录
			err = ParseAllPath(file+"/"+sub.Name(), opts...)
			if err != nil {
				return
			}
			continue
		}
		if sub.Name()[len(sub.Name())-3:] == ".go" {
			err = ParseTag(file+"/"+sub.Name(), opts...)
			if err != nil {
				return
			}
		}
	}
	return
}

func ParseTag(file string, opts ...GenTagOption) (err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	fSet := token.NewFileSet()
	node, err := parser.ParseFile(fSet, file, data, parser.ParseComments)
	if err != nil {
		return
	}
	fn := fSetNode{
		fSet: fSet,
		node: node,
	}
	res, err := parseContent(fn, opts...)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(file, res, os.FileMode(0664))
	return
}

type fSetNode struct {
	fSet *token.FileSet
	node ast.Node
}

type tagVisitor struct {
	opts []GenTagOption
}

func (t *tagVisitor) Visit(node ast.Node) ast.Visitor {
	ts, ok := node.(*ast.TypeSpec)
	if !ok {
		return t
	}
	st, ok := ts.Type.(*ast.StructType)
	if !ok {
		return t
	}
	return t.visit(ts.Name.Name, st)
}

func (t *tagVisitor) visit(name string, st *ast.StructType) ast.Visitor {
	for i, f := range st.Fields.List {
		st.Fields.List[i].Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: genTag(name, f, t.opts),
		}
	}
	return nil
}

func genTag(sName string, node *ast.Field, opts []GenTagOption) (tg string) {
	value := ""
	if node.Tag != nil {
		value = strings.ReplaceAll(node.Tag.Value,"`","")
	}
	tags, err := Parse(value)
	if err != nil {
		return
	}

	for _, opt := range opts {
		if !strings.HasSuffix(sName, opt.Suffix) {
			continue
		}
		// 有key不覆盖
		oldTag,hasKey := tags.HasKey(opt.Key)
		if hasKey && !opt.Cover {
			continue
		}
		if oldTag!= nil && oldTag.Name == "-" {
			continue
		}
		tagName := ""
		switch opt.TagType {
		case camelCase:
			tagName = strcase.ToLowerCamel(node.Names[0].Name)
		case snakeCase:
			tagName = strcase.ToSnake(node.Names[0].Name)
		default:
			panic(errors.New("暂不支持"))
		}
		tgStruct := &Tag{
			Key:     opt.Key,
			Name:    tagName,
			Options: opt.Options,
		}
		err = tags.Set(tgStruct)
		if err != nil {
			panic(err)
		}
	}
	return tags.String()
}

// 解析文件内容
func parseContent(fn fSetNode, opts ...GenTagOption) (res []byte, err error) {
	tv := tagVisitor{
		opts: opts,
	}
	ast.Walk(&tv, fn.node)
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fn.fSet, fn.node)
	if err != nil {
		return
	}
	res, err = format.Source(buffer.Bytes())
	return
}
