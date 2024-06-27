package inventory

import (
	"encoding/json"
	"testing"

	"github.com/bmc-toolbox/common"
	"github.com/stretchr/testify/require"
)

type clutter struct {
	Payload string `json:"payload,omitempty"`
}

type testContainer struct {
	Clutter *clutter       `json:"clutter,omitempty"`
	Status  *common.Status `json:"status,omitempty"`
}

func TestStatusFromJSON(t *testing.T) {
	t.Run("cluttered array", func(t *testing.T) {
		ary := []*testContainer{
			{
				Clutter: &clutter{
					Payload: "stuff",
				},
			},
			{
				Status: &common.Status{
					Health: "great",
					State:  "awesome",
				},
			},
			{
				Clutter: &clutter{},
			},
		}
		data, err := json.Marshal(ary)
		require.NoError(t, err, "pre-requisite")
		st, err := statusFromJSON(data)
		require.NoError(t, err, "function call")
		require.NotNil(t, st)
		require.Equal(t, "great", st.Health)
	})
	t.Run("serialized empty object returns empty object", func(t *testing.T) {
		ary := []*testContainer{
			{
				Status: &common.Status{},
			},
		}
		data, err := json.Marshal(ary)
		require.NoError(t, err, "pre-requisite")
		st, err := statusFromJSON(data)
		require.NoError(t, err, "function call")
		require.NotNil(t, st)
	})
	t.Run("empty non-nil array payload returns nil", func(t *testing.T) {
		ary := []*testContainer{}
		data, err := json.Marshal(ary)
		require.NoError(t, err, "pre-requisite")
		st, err := statusFromJSON(data)
		require.NoError(t, err, "function call")
		require.Nil(t, st)
	})
}
