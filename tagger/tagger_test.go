package tagger

import (
	"testing"
)

func TestParseAllPath(t *testing.T) {
	err := ParseAllPath("../service",
		CamelCase("","json",true,"omitempty"),
		CamelCase("Param","form",true),
		)
	if err != nil {
		t.Error(err)
		return
	}
}