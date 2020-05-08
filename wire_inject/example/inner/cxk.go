package inner

//@inject(set=service)
type CXK struct {
	Singer Singer
	Dancer Dancer
	Rapper Rapper
	Player Player
}

func (c *CXK) Sing()  {
	c.Singer.Sing()
}

func (c *CXK) Dance() {
	c.Dancer.Dance()
}

func (c *CXK) Rap()  {
	c.Rapper.Rap()
}

func (c *CXK) Ball()  {
	c.Player.Ball()
}

func NewCXK(s Singer, d Dancer, r Rapper, p Player) *CXK {
	return &CXK{
		Singer: s,
		Dancer: d,
		Rapper: r,
		Player: p,
	}
}
