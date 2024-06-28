package fleetdbapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestGetEventByID(t *testing.T) {
	s := serverTest(t)
	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		evtID := uuid.MustParse(dbtools.FixtureEventHistories[0].EventID)
		evt, _, err := s.Client.GetEventByID(ctx, evtID)
		if !expectError {
			// we should get back the fixture data for the inventory server
			// we don't care about the name of the server
			require.NoError(t, err)
			require.NotNil(t, evt)
			require.Equal(t, dbtools.FixtureEventHistoryServer.ID, evt.Target.String())
			require.Equal(t, "succeeded", evt.FinalState)
		}
		return err
	})

	// event id not found
	bogus := uuid.New()
	_, _, err := s.Client.GetEventByID(context.Background(), bogus)
	se := &fleetdbapi.ServerError{}
	require.ErrorAs(t, err, se)
	expStr := fmt.Sprintf("no event for id %s", bogus.String())
	require.Equal(t, 404, se.StatusCode)
	require.Equal(t, expStr, se.Message)
}

func TestGetServerEvents(t *testing.T) {
	s := serverTest(t)
	realClientTests(t, func(ctx context.Context, auth string, code int, expErr bool) error {
		s.Client.SetToken(auth)
		srvID := uuid.MustParse(dbtools.FixtureEventHistoryServer.ID)
		evts, _, err := s.Client.GetServerEvents(ctx, srvID, nil)
		if !expErr {
			require.NoError(t, err)
			require.Equal(t, 3, len(evts))
		}
		return err
	})
	bogus := uuid.New()
	_, _, err := s.Client.GetServerEvents(context.Background(), bogus, &fleetdbapi.PaginationParams{})
	se := &fleetdbapi.ServerError{}
	require.ErrorAs(t, err, se)
	expStr := fmt.Sprintf("no events for server %s", bogus.String())
	require.Equal(t, 404, se.StatusCode)
	require.Equal(t, expStr, se.Message)
}

func TestUpdateEvent(t *testing.T) {
	s := serverTest(t)
	evt := &fleetdbapi.Event{
		EventID:    uuid.New(),
		Type:       "test event 2",
		Start:      time.Now().Add(-1 * time.Minute),
		End:        time.Now().Add(-30 * time.Second),
		Target:     uuid.MustParse(dbtools.FixtureEventHistoryServer.ID),
		FinalState: "succeeded",
	}
	realClientTests(t, func(ctx context.Context, auth string, code int, expErr bool) error {
		s.Client.SetToken(auth)
		r, err := s.Client.UpdateEvent(ctx, evt)
		if !expErr {
			require.NoError(t, err)
			require.Equal(t, evt.EventID.String(), r.Slug)
		}
		return err
	})
	// repeat with the same payload, get the same result
	r, err := s.Client.UpdateEvent(context.Background(), evt)
	require.NoError(t, err)
	require.Equal(t, evt.EventID.String(), r.Slug)
	// change the payload, get an error
	evt.FinalState = "failed"
	r, err = s.Client.UpdateEvent(context.Background(), evt)
	se := &fleetdbapi.ServerError{}
	require.ErrorAs(t, err, se)
	expStr := fmt.Sprintf("id in use: %s", evt.EventID.String())
	require.Equal(t, 400, se.StatusCode)
	require.Equal(t, expStr, se.Message)
}
