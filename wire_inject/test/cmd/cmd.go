package main

import "utils/wire_inject"

func main() {
	wire_inject.RunWire("cmd/server/inject")
}
