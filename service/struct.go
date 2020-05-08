package service

type StudentRes struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type StudentParam struct {
	Name string `json:"name,omitempty" form:"name"`
	Age  int    `json:"age,omitempty" form:"age"`
}
