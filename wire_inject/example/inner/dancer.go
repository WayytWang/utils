package inner

import "fmt"

type Dancer interface {
	Dance()
}

//@inject(inner.Dancer,set=helper)
type ZhaoSi struct {

}

func (z *ZhaoSi) Dance() {
	fmt.Println("你四哥在气质这块把握的很好")
}