package main

import "utils/wire_inject"

func main() {
	wire_inject.RunWire("wire_inject/example/cmd/server/inject")
}
