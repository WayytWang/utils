package tagger

const (
	camelCase = 1
	snakeCase = 2
)

type GenTagOption struct {
	Suffix  string   // 结构体名后缀
	Key     string   // tagKey
	TagType int      // 驼峰/下划线
	Cover   bool     // 是否覆盖
	Options []string // 选项
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
