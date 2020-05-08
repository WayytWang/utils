package tag

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Tags struct {
	tags []*Tag
}

func NewTags(tags ...*Tag) *Tags {
	return &Tags{tags: tags}
}

func (t *Tags) Append(tags ...*Tag) {
	t.tags = append(t.tags, tags...)
}

var (
	errTagSyntax      = errors.New("bad syntax for param tag pair")
	errTagKeySyntax   = errors.New("bad syntax for param tag key")
	errTagValueSyntax = errors.New("bad syntax for param tag value")

	errKeyNotSet      = errors.New("tag key does not exist")
	errTagNotExist    = errors.New("tag does not exist")
	errTagKeyMismatch = errors.New("mismatch between key and tag.key")
)

type Tag struct {
	// eg: `json:"wang,omitempty"`
	// key : json
	Key string
	// value : wang
	Name string
	// options : ["omitempty"]
	Options []string
}

func Parse(tag string) (*Tags, error) {
	var tags []*Tag

	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		// eg: `json:"wang,omitempty" form:"wang"`
		// 拿到key的下一个字符索引
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}

		if i == 0 {
			return nil, errTagKeySyntax
		}
		if i+1 > len(tag) || tag[i] != ':' {
			return nil, errTagSyntax
		}
		if tag[i+1] != '"' {
			return nil, errTagValueSyntax
		}
		key := tag[:i]
		tag = tag[i+1:]

		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			return nil, errTagValueSyntax
		}
		qValue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qValue)
		if err != nil {
			return nil, errTagValueSyntax
		}
		res := strings.Split(value, ",")
		name := res[0]
		options := res[1:]
		if len(options) == 0 {
			options = nil
		}
		tags = append(tags, &Tag{
			Key:     key,
			Name:    name,
			Options: options,
		})
	}

	return &Tags{tags: tags}, nil
}

func (t *Tags) HasKey(key string) (*Tag,bool) {
	for _, tag := range t.tags {
		if tag.Key == key {
			return tag,true
		}
	}
	return nil,false
}

func (t *Tags) Get(key string) (*Tag, error) {
	for _, tag := range t.tags {
		if tag.Key == key {
			return tag, nil
		}
	}
	return nil, errTagNotExist
}

func (t *Tags) Set(tag *Tag) error {
	if tag.Key == "" {
		return errKeyNotSet
	}

	added := false
	for i, tg := range t.tags {
		if tg.Key == tag.Key {
			added = true
			t.tags[i] = tag
		}
	}

	if !added {
		t.tags = append(t.tags, tag)
	}
	return nil
}

// Options 相关
func (t *Tag) HasOption(opt string) bool {
	for _, tagOpt := range t.Options {
		if tagOpt == opt {
			return true
		}
	}
	return false
}

func (t *Tags) AddOptions(key string, options ...string) {
	for i, tg := range t.tags {
		if tg.Key != key {
			continue
		}

		for _, opt := range options {
			if !tg.HasOption(opt) {
				tg.Options = append(tg.Options, opt)
			}
		}

		t.tags[i] = tg
	}
}

func (t *Tags) DeleteOptions(key string, options ...string) {
	// 原本的option是否在需要被删除的列表中
	hasOption := func(option string) bool {
		for _, opt := range options {
			if option == opt {
				return true
			}
		}
		return false
	}
	for i, tg := range t.tags {
		if tg.Key != key {
			continue
		}

		var updated []string
		for _, opt := range options {
			if !hasOption(opt) {
				updated = append(updated, opt)
			}
		}
		tg.Options = updated
		t.tags[i] = tg
	}
}

func (t *Tags) Delete(keys ...string) {
	hasKey := func(key string) bool {
		for _, k := range keys {
			if key == k {
				return true
			}
		}
		return false
	}

	var updated []*Tag
	for _, tg := range t.tags {
		if !hasKey(tg.Key) {
			updated = append(updated, tg)
		}
	}
	t.tags = updated
}

func (t *Tag) Value() string {
	options := strings.Join(t.Options, ",")
	if options != "" {
		return fmt.Sprintf(`%s,%s`, t.Name, options)
	}
	return t.Name
}

func (t *Tag) String() string {
	return fmt.Sprintf(`%s:%q`, t.Key, t.Value())
}

func (t *Tags) String() string {
	tags := t.tags
	if len(tags) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("`")
	for i, tag := range t.tags {
		buf.WriteString(tag.String())
		if i != len(t.tags)-1 {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("`")
	return buf.String()
}
