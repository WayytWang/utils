package inner

import "fmt"

type Rapper interface {
	Rap()
}

//@inject(inner.Rapper,set=helper)
type ChrisWu struct {

}

func (c *ChrisWu) Rap() {
	fmt.Println("你看这个碗，又大又圆")
}
