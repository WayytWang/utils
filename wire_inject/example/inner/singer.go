package inner

import "fmt"

type Singer interface {
	Sing()
}

//@inject(inner.Singer,set=helper)
type ZXY struct {
}

func (z ZXY) Sing() {
	fmt.Println("我要送你红色玫瑰")
}

func InitZXY() *ZXY {
	return &ZXY{}
}
