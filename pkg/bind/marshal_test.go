package bind

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalParams(t *testing.T) {
	got, err := MarshalParams("1", []string{"2", "3"}, "4", "5", "6")
	require.NoError(t, err)

	j, _ := json.MarshalIndent(got, "", "  ")
	fmt.Println(string(j))
}
