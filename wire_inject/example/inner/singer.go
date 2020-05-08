package inner

type Singer interface {
	Sing()
}

//@inject(inner.Singer,set=helper)
type ZXY struct {

}

func (z ZXY) Sing() {

}

func InitZXY() *ZXY {
	return &ZXY{}
}
