package example

import "utils/wire_inject/example/inner"

//@inject(set=service)
type CXK struct {
	Singer inner.Singer
	Dancer inner.Dancer
	Rapper inner.Rapper
	Player inner.Player
}

func (c *CXK) Sing()  {}
func (c *CXK) Dance() {}
func (c *CXK) Rap()   {}
func (c *CXK) Ball()  {}

func NewCXK(s inner.Singer, d inner.Dancer, r inner.Rapper, p inner.Player) *CXK {
	return &CXK{
		Singer: s,
		Dancer: d,
		Rapper: r,
		Player: p,
	}
}
