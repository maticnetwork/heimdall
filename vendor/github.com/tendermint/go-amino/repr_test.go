package amino

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	a string
	b int
	c []*Foo
	D string // exposed
}

type pair struct {
	Key   string
	Value interface{}
}

func (pr pair) get(key string) (value interface{}) {
	if pr.Key != key {
		panic(fmt.Sprintf("wanted %v but is %v", key, pr.Key))
	}
	return pr.Value
}

func (f Foo) MarshalAmino() ([]pair, error) {
	return []pair{
		{"a", f.a},
		{"b", f.b},
		{"c", f.c},
		{"D", f.D},
	}, nil
}

func (f *Foo) UnmarshalAmino(repr []pair) error {
	f.a = repr[0].get("a").(string)
	f.b = repr[1].get("b").(int)
	f.c = repr[2].get("c").([]*Foo)
	f.D = repr[3].get("D").(string)
	return nil
}

func TestMarshalAminoBinary(t *testing.T) {

	cdc := NewCodec()
	cdc.RegisterInterface((*interface{})(nil), nil)
	cdc.RegisterConcrete(string(""), "string", nil)
	cdc.RegisterConcrete(int(0), "int", nil)
	cdc.RegisterConcrete(([]*Foo)(nil), "[]*Foo", nil)

	var f = Foo{
		a: "K",
		b: 2,
		c: []*Foo{nil, nil, nil},
		D: "J",
	}
	bz, err := cdc.MarshalBinary(f)
	assert.Nil(t, err)

	t.Logf("bz %X", bz)

	var f2 Foo
	err = cdc.UnmarshalBinary(bz, &f2)
	assert.Nil(t, err)

	assert.Equal(t, f, f2)
	assert.Equal(t, f.a, f2.a) // In case the above doesn't check private fields?
}

func TestMarshalAminoJSON(t *testing.T) {

	cdc := NewCodec()
	cdc.RegisterInterface((*interface{})(nil), nil)
	cdc.RegisterConcrete(string(""), "string", nil)
	cdc.RegisterConcrete(int(0), "int", nil)
	cdc.RegisterConcrete(([]*Foo)(nil), "[]*Foo", nil)

	var f = Foo{
		a: "K",
		b: 2,
		c: []*Foo{nil, nil, nil},
		D: "J",
	}
	bz, err := cdc.MarshalJSON(f)
	assert.Nil(t, err)

	t.Logf("bz %X", bz)

	var f2 Foo
	err = cdc.UnmarshalJSON(bz, &f2)
	assert.Nil(t, err)

	assert.Equal(t, f, f2)
	assert.Equal(t, f.a, f2.a) // In case the above doesn't check private fields?
}
