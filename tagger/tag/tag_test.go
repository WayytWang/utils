package tag

import (
	"testing"
)

func TestParse(t *testing.T) {
	tags,err := Parse(`json:"wang,omitempty" form:"wang"`)
	if err != nil {
		t.Error(err)
		return
	}
	for _, tag := range tags.tags {
		t.Logf("%+v",*tag)
	}
}