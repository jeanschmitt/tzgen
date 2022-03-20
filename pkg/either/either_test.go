package either

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLeft(t *testing.T) {
	got := Left[string, int]("hello")
	require.Equal(t, left[string, int]{"hello"}, got)
}

func TestRight(t *testing.T) {
	got := Right[string, int](42)
	require.Equal(t, right[string, int]{42}, got)
}
