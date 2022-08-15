package bind

import (
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"context"
	"github.com/pkg/errors"
	"strconv"
)

type Bigmap[K, V any] struct {
	id      int64
	keyType *micheline.Type
	rpc     RPC
}

func NewBigmap[K, V any](id int64) *Bigmap[K, V] {
	return &Bigmap[K, V]{id: id}
}

func (b *Bigmap[K, B]) ID() int64 {
	return b.id
}

func (b *Bigmap[K, B]) SetRPC(client RPC) *Bigmap[K, B] {
	b.rpc = client
	return b
}

func (b *Bigmap[K, B]) SetKeyType(keyType micheline.Type) *Bigmap[K, B] {
	b.keyType = &keyType
	return b
}

func (b *Bigmap[K, V]) Get(ctx context.Context, key K) (v V, err error) {
	if b.rpc == nil {
		return v, errors.New("rpc not set in bigmap")
	}

	keyVal, err := MarshalPrim(key, true)
	if err != nil {
		return v, err
	}

	if b.keyType == nil {
		b.SetKeyType(keyVal.BuildType())
	}

	k, err := micheline.NewKey(*b.keyType, keyVal)
	if err != nil {
		return v, err
	}

	keyHash := k.Hash()

	prim, err := b.rpc.GetBigmapValue(ctx, b.id, keyHash, rpc.Head)
	if err != nil {
		var httpError rpc.HTTPError
		if errors.As(err, &httpError) && httpError.StatusCode() == 404 {
			return v, &ErrKeyNotFound{key: keyHash.String()}
		}
		return v, err
	}

	if err = UnmarshalPrim(prim, &v); err != nil {
		return v, err
	}

	return v, nil
}

func (b *Bigmap[K, V]) String() string {
	return "Bigmap#" + strconv.Itoa(int(b.id))
}

func (b *Bigmap[K, V]) UnmarshalPrim(prim micheline.Prim) error {
	*b = *NewBigmap[K, V](prim.Int.Int64())
	return nil
}

type ErrKeyNotFound struct {
	key string
}

func (e *ErrKeyNotFound) Error() string {
	return "bigmap key not found: " + e.key
}

func (e *ErrKeyNotFound) Is(target error) bool {
	other, ok := target.(*ErrKeyNotFound)
	if !ok {
		return false
	}
	return e.key == other.key
}
