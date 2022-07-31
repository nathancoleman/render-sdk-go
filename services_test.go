package render

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServices_List(t *testing.T) {
	c, err := NewClient(DefaultConfig())
	require.NoError(t, err)

	services, err := c.Services().List(context.TODO())
	require.NoError(t, err)
	fmt.Println(services)
}
