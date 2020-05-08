package param

type Student struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type StudentParam struct {
	Name string `json:"name,omitempty" form:"name"`
	Age  int    `json:"age,omitempty" form:"age"`
}

type Teacher struct {
	Name    string `json:"name,omitempty"`
	Age     string `json:"age,omitempty"`
	Subject string `json:"subject,omitempty"`
}

type TeacherParam struct {
	Name    string `json:"name,omitempty" form:"name"`
	Age     string `json:"age,omitempty" form:"age"`
	Subject string `json:"subject,omitempty" form:"subject"`
}
