package main

import "utils/wire_inject/example/cmd/server/inject"

func main() {
	cxk := inject.InitCXK()
	cxk.Sing()
	cxk.Dance()
	cxk.Rap()
	cxk.Ball()
}
