//go:build testtools && integration

package dbtools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerComponentTypes(t *testing.T) {
	db := DatabaseTest(t)
	ctx := context.TODO()

	err := SetupComponentTypes(ctx, db)
	require.NoError(t, err)

	for _, typ := range knownComponentTypes {
		_, err := ComponentTypeIDFromName(typ)
		require.NoError(t, err, "couldn't find %s", typ)
	}

	_, err = ComponentTypeIDFromName("bogus")
	require.Error(t, err, "no error on bogus")
}
