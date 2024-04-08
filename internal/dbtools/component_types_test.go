//go:build testtools && integration

package dbtools

import (
	"context"
	"testing"

	"github.com/bmc-toolbox/common"
	"github.com/stretchr/testify/require"
)

func TestServerComponentTypes(t *testing.T) {
	db := DatabaseTest(t)
	ctx := context.TODO()

	err := SetupComponentTypes(ctx, db)
	require.NoError(t, err)

	for _, typ := range knownComponentTypes {
		_, err := ComponentTypeIDFromName(ctx, db, typ)
		require.NoError(t, err, "couldn't find %s", typ)
	}

	require.NotPanics(t, func() { _ = MustComponentTypeID(ctx, db, common.SlugBackplaneExpander) })

	require.Panics(t, func() { _ = MustComponentTypeID(ctx, db, "bogus") })

	_, err = ComponentTypeIDFromName(ctx, db, "bogus")
	require.Error(t, err, "no error on bogus")
}
