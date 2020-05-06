package tagger

// 用来配置
// 1.需要tagger的结构体名后缀
// 2.需要tagger的tagName和方式(下划线/驼峰) map["tagName"]"tagType"
// 3.是否覆盖现有tag

const (
	camelCase = 1
	snakeCase = 2
)

type GenTagOption struct {
	Suffix  string
	Key     string
	TagType int
	Cover   bool
	Options []string
}

func CamelCase(suffix, tag string, cover bool, opt ...string) GenTagOption {
	return GenTagOption{
		Suffix:  suffix,
		Key:     tag,
		TagType: 1,
		Cover:   cover,
		Options: opt,
	}
}

func SnakeCase(suffix, tag string, cover bool, opt ...string) GenTagOption {
	return GenTagOption{
		Suffix:  suffix,
		Key:     tag,
		TagType: 2,
		Cover:   cover,
		Options: opt,
	}
}
