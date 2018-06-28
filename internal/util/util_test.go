package util_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dacjames/persistsql/internal/util"
)

func TestPlaceholders(t *testing.T) {
	require.Equal(t, true, true)

	p := util.NewPlaceholders()

	require.Equal(t, "$1", p.NextPlaceholder())
	require.Equal(t, "$2", p.NextPlaceholder())
	require.Equal(t, "$3", p.NextPlaceholder())
}

func TestPlaceholdersAt(t *testing.T) {
	require.Equal(t, true, true)

	p := util.NewPlaceholdersAt(3)

	require.Equal(t, "$3", p.NextPlaceholder())
	require.Equal(t, "$4", p.NextPlaceholder())
	require.Equal(t, "$5", p.NextPlaceholder())
}

func TestPlaceholdersValue(t *testing.T) {
	p := util.NewPlaceholders()

	require.Equal(t, "($1)", p.NextValue(1))
	require.Equal(t, "()", p.NextValue(0))
	require.Equal(t, "($2, $3, $4)", p.NextValue(3))
}

type InnerA struct {
	X int `db:"x"`
	Y int
}

type InnerB struct {
	Y int
	Z int `db:"z"`
}

type Outer struct {
	IA InnerA
	IB *InnerB
}

func TestNestedTags(t *testing.T) {
	s := Outer{
		IA: InnerA{},
		IB: &InnerB{},
	}

	i := interface{}(s)

	expected := []string{"x", "z"}

	require.Equal(t, expected, util.NestedTags(reflect.TypeOf(i), "db"))
}
