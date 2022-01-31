package builder

import (
	"blockwatch.cc/tzgo/micheline"
	"fmt"
	"math/big"
)

func Build(prims ...interface{}) micheline.Prim {
	b := New()
	for _, p := range prims {
		b = b.Any(p)
	}
	return b.Finalize()
}

type Builder struct {
	prim *micheline.Prim
	ptr  *micheline.Prim
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Int(i *big.Int) *Builder {
	return b.addPrim(intPrim(i))
}

func (b *Builder) Bytes(bytes []byte) *Builder {
	return b.addPrim(bytesPrim(bytes))
}

func (b *Builder) String(s string) *Builder {
	return b.addPrim(stringPrim(s))
}

func (b *Builder) Bool(v bool) *Builder {
	return b.addPrim(boolPrim(v))
}

func (b *Builder) Primer(primer Primer) *Builder {
	return b.addPrim(primer.ToPrim())
}

func (b *Builder) Seq(seq interface{}) *Builder {
	return b.addPrim(seqPrim(seq))
}

func (b *Builder) Option(o interface{}) *Builder {
	return b.addPrim(optionPrim(o))
}

func (b *Builder) Any(v interface{}) *Builder {
	prim, err := ToPrim(v)
	if err != nil {
		panic(fmt.Sprintf("invalid primitive: %v", err))
	}
	return b.addPrim(prim)
}

func (b *Builder) Finalize() micheline.Prim {
	if b.prim == nil {
		return micheline.NewPrim(micheline.D_UNIT)
	}
	return *b.prim
}

func (b *Builder) addPrim(prim micheline.Prim) *Builder {
	if b.prim == nil {
		b.prim = &prim
		b.ptr = b.prim
		return b
	}
	*b.ptr = micheline.NewPairValue(*b.ptr, prim)
	b.ptr = &b.ptr.Args[1]
	return b
}
