// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package inject

import (
	"utils/wire_inject/example/inner"
)

// Injectors from wire.go:

func InitCXK() *inner.CXK {
	zxy := inner.InitZXY()
	zhaoSi := &inner.ZhaoSi{}
	chrisWu := &inner.ChrisWu{}
	kobe := &inner.Kobe{}
	cxk := inner.NewCXK(zxy, zhaoSi, chrisWu, kobe)
	return cxk
}