package utils

import (
	"testing"

	"utils/tagger"
)

func TestParseAllPath(t *testing.T) {
	// 在../service目录下所有的.go文件中修改option对应的结构体字段tag
	err := tagger.ParseAllPath("../param",
		tagger.CamelCase("", "json", true, "omitempty"),
		tagger.CamelCase("Param", "form", true),
	)
	if err != nil {
		t.Error(err)
		return
	}
}
