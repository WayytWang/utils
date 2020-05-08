//+build wireinject

package inject

import (
	"github.com/google/wire"
	"utils/wire_inject/example"
)

func InitCXK()(*example.CXK) {
	panic(wire.Build(sets))
}