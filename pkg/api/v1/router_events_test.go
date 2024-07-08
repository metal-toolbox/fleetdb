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

func TestGetHistoryByID(t *testing.T) {
	s := serverTest(t)
	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		evts, _, err := s.Client.GetHistoryByID(ctx, dbtools.FixtureEventHistoryRelatedID)
		if !expectError {
			// we should get back the fixture data
			require.NoError(t, err)
			require.Len(t, evts, 3)
		}
		return err
	})

	// event id not found
	bogus := uuid.New()
	_, _, err := s.Client.GetHistoryByID(context.Background(), bogus)
	require.Error(t, err)
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
			require.Len(t, evts, 6)
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

	// explicitly test adding 2 events with the same UUID and target with different types
	// this is exactly the "composite condition" case from ConditionOrc
	relatedEvt := &fleetdbapi.Event{
		EventID:    evt.EventID,
		Type:       "related type",
		Start:      evt.Start,
		End:        evt.End,
		Target:     evt.Target,
		FinalState: "succeeded",
	}

	_, err = s.Client.UpdateEvent(context.Background(), relatedEvt)
	require.NoError(t, err)
}
