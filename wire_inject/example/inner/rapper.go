package inner

type Rapper interface {
	Rapper()
}

//@inject(inner.Rapper,set=helper)
type ChrisWu struct {

}

func (c *ChrisWu) Rapper() {

}
