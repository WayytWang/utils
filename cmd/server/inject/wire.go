//+build wireinject

package inject

import (
	"github.com/google/wire"
	"utils/wire_inject/test"
)

func InitApplication()(*test.Application) {
	panic(wire.Build(sets))
}