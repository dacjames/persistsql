package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type Placeholders interface {
	NextPlaceholder() string
	NextValue(n int) string
}

type placeholders struct {
	i int
}

func (p *placeholders) NextPlaceholder() string {
	next := fmt.Sprintf("$%d", p.i)
	p.i = p.i + 1
	return next
}

func (p *placeholders) NextValue(n int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = p.NextPlaceholder()
	}
	return `(` + strings.Join(parts, ", ") + `)`
}

func NewPlaceholders() Placeholders {
	return NewPlaceholdersAt(1)
}

func NewPlaceholdersAt(start int) Placeholders {
	return &placeholders{i: start}
}

func NewRandSource() *rand.Rand {
	t := time.Now().UTC()
	return rand.New(rand.NewSource(t.UnixNano()))
}

func NestedTags(t reflect.Type, key string) []string {
	if t.Kind() == reflect.Ptr {
		return NestedTags(t.Elem(), key)
	}

	values := []string{}

	if t.Kind() != reflect.Struct {
		return values
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Ptr {
			values = append(values, NestedTags(field.Type.Elem(), key)...)
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			values = append(values, NestedTags(field.Type, key)...)
			continue
		}

		// if field.Type.Kind() == reflect.Ptr {
		// 	field.Type.Elem()
		// }

		value, ok := field.Tag.Lookup(key)
		if ok {
			values = append(values, value)
		}
	}

	return values
}
