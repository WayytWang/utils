package inner

type Dancer interface {
	Dance()
}

//@inject(inner.Dancer,set=helper)
type ZhaoSi struct {

}

func (z *ZhaoSi) Dance() {

}