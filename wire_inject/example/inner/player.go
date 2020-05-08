package inner

type Player interface {
	Ball()
}

//@inject(inner.Player,set=helper)
type Kobe struct {

}

func (k *Kobe) Ball() {

}