package inner

import "fmt"

type Player interface {
	Ball()
}

//@inject(inner.Player,set=helper)
type Kobe struct {

}

func (k *Kobe) Ball() {
	fmt.Println("这球就要给科比")
}