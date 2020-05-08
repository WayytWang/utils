//+build wireinject

package inject

import (
"github.com/google/wire"
"utils/wire_inject/example/inner"
)

func InitCXK()(*inner.CXK) {
	panic(wire.Build(sets))
}
