package fleetdbapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"github.com/metal-toolbox/fleetdb/internal/metrics"
	"github.com/metal-toolbox/fleetdb/internal/models"
)

type Event struct {
	EventID     uuid.UUID       `json:"event_id" binding:"required,uuid4_rfc4122"`
	Type        string          `json:"event_type" binding:"required"`
	Start       time.Time       `json:"event_start" binding:"required,ltfield=End"`
	End         time.Time       `json:"event_end" binding:"required,gtfield=Start"`
	Target      uuid.UUID       `json:"target_server" binding:"required,uuid4_rfc4122"`
	Parameters  json.RawMessage `json:"parameters,omitempty" binding:"-"`
	FinalState  string          `json:"final_state" binding:"required"`
	FinalStatus json.RawMessage `json:"final_status,omitempty" binding:"-"`
}

func (r *Router) getHistoryByConditionID(c *gin.Context) {
	evtID, err := uuid.Parse(c.Param("evtID"))
	if err != nil {
		badRequestResponse(c, "failed to parse event id", err)
		return
	}

	// we expect a small number of events with the same id, O(5-10)
	ehs, err := models.EventHistories(
		models.EventHistoryWhere.EventID.EQ(evtID.String()),
	).All(c.Request.Context(), r.DB)

	if err != nil {
		metrics.DBError("fetching event history")
		r.Logger.With(
			zap.Error(err),
			zap.String("event_id", evtID.String()),
		).Warn("event history by id")
		dbErrorResponse(c, err)
		return
	}

	if len(ehs) == 0 {
		msg := fmt.Sprintf("no event for id %s", evtID.String())
		notFoundResponse(c, msg)
		return
	}

	events := make([]*Event, 0, len(ehs))
	for _, eh := range ehs {
		evt := &Event{
			EventID:    uuid.MustParse(eh.EventID),
			Type:       eh.EventType,
			Start:      eh.EventStart,
			End:        eh.EventEnd,
			Target:     uuid.MustParse(eh.TargetServer),
			FinalState: eh.FinalState,
		}

		if eh.Parameters.Valid {
			evt.Parameters = eh.Parameters.JSON
		}

		if eh.FinalStatus.Valid {
			evt.FinalStatus = eh.FinalStatus.JSON
		}
		events = append(events, evt)
	}

	pd := paginationData{
		pageCount: len(events),
	}

	listResponse(c, events, pd)
}

func (r *Router) getServerEvents(c *gin.Context) {
	srvID, err := uuid.Parse(c.Param("srvID"))
	if err != nil {
		badRequestResponse(c, "failed to parse target server id", err)
		return
	}

	// parse pagination -- only honoring limit and page here
	pageParams, err := parsePagination(c)
	if err != nil {
		badRequestResponse(c, "parsing pagination", err)
		return
	}

	limit := pageParams.Limit // default is 100
	if limit > 1000 {         //nolint:gomnd // it's fine
		limit = 1000 // more than 1000 event records in a single shot is sus
	}

	var offset int
	if pageParams.Page > 1 { // default page is 1
		offset = limit * (pageParams.Page - 1)
	}

	// N.B. count returns 0, not sql.ErrNoRows
	historyTotal, err := models.EventHistories(
		models.EventHistoryWhere.TargetServer.EQ(srvID.String()),
	).Count(c.Request.Context(), r.DB)

	if err != nil {
		r.Logger.With(
			zap.Error(err),
			zap.String("server_id", srvID.String()),
		).Warn("counting event history")
		metrics.DBError("counting event history")
		dbErrorResponse(c, err)
		return
	}

	if historyTotal == 0 {
		msg := fmt.Sprintf("no events for server %s", srvID.String())
		notFoundResponse(c, msg)
		return
	}

	ehs, err := models.EventHistories(
		models.EventHistoryWhere.TargetServer.EQ(srvID.String()),
		qm.OrderBy("event_id, event_end DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(c.Request.Context(), r.DB)

	if err != nil {
		metrics.DBError("event history by server")
		r.Logger.With(
			zap.Error(err),
			zap.String("target_id", srvID.String()),
		).Warn("fetching server event history")
		dbErrorResponse(c, err)
		return
	}

	var evts []*Event
	for _, eh := range ehs {
		evt := &Event{
			EventID:    uuid.MustParse(eh.EventID),
			Type:       eh.EventType,
			Start:      eh.EventStart,
			End:        eh.EventEnd,
			Target:     uuid.MustParse(eh.TargetServer),
			FinalState: eh.FinalState,
		}

		if eh.Parameters.Valid {
			evt.Parameters = eh.Parameters.JSON
		}

		if eh.FinalStatus.Valid {
			evt.FinalStatus = eh.FinalStatus.JSON
		}

		evts = append(evts, evt)
	}

	pd := paginationData{
		pageCount:  len(evts),
		totalCount: historyTotal,
		pager:      pageParams,
	}
	listResponse(c, evts, pd)
}

func equivalentEvents(evt *Event, eh *models.EventHistory) bool {
	// test everything but the time
	return evt.EventID.String() == eh.EventID &&
		evt.Target.String() == eh.TargetServer &&
		evt.Type == eh.EventType &&
		evt.FinalState == eh.FinalState &&
		bytes.Equal(evt.Parameters, eh.Parameters.JSON) &&
		bytes.Equal(evt.FinalStatus, eh.FinalStatus.JSON)
}

func (r *Router) updateEvent(c *gin.Context) {
	evtID, err := uuid.Parse(c.Param("evtID"))
	if err != nil {
		badRequestResponse(c, "failed to parse event id", err)
		return
	}

	evt := &Event{}
	if err := c.ShouldBindJSON(evt); err != nil {
		badRequestResponse(c, "invalid event payload", err)
		return
	}

	// initial sanity check
	if evtID != evt.EventID {
		badRequestResponse(c, "payload does not match presented id", nil)
	}

	ctx := c.Request.Context()
	// shortcut if we've seen this event before already
	existing, err := models.EventHistories(
		models.EventHistoryWhere.EventID.EQ(evt.EventID.String()),
		models.EventHistoryWhere.EventType.EQ(evt.Type),
		models.EventHistoryWhere.TargetServer.EQ(evt.Target.String()),
	).One(ctx, r.DB)

	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err == nil:
		if equivalentEvents(evt, existing) {
			createdResponse(c, existing.EventID)
			return
		}
		badRequestResponse(c, fmt.Sprintf("id in use: %s", existing.EventID), errors.New("existing event"))
		return
	default:
		metrics.DBError("fetching event history")
		r.Logger.With(
			zap.Error(err),
			zap.String("event_id", evt.EventID.String()),
		).Warn("event history by id")
		dbErrorResponse(c, err)
		return
	}

	eh := &models.EventHistory{
		EventID:      evt.EventID.String(),
		EventType:    evt.Type,
		EventStart:   evt.Start,
		EventEnd:     evt.End,
		TargetServer: evt.Target.String(),
		Parameters:   null.JSONFrom([]byte(evt.Parameters)),
		FinalState:   evt.FinalState,
		FinalStatus:  null.JSONFrom([]byte(evt.FinalStatus)),
	}

	txn, err := r.DB.Begin()
	if err != nil {
		r.Logger.With(
			zap.Error(err),
			zap.String("event_id", evt.EventID.String()),
		).Warn("unable to create transaction")
		metrics.DBError("creating transaction")
		dbErrorResponse(c, err)
		return
	}

	doRollback := false
	rollbackFn := func() {
		if doRollback {
			if rbErr := txn.Rollback(); rbErr != nil {
				r.Logger.With(
					zap.Error(rbErr),
				).Warn("rollback error on event insertion")
				metrics.DBError("rollback event insertion")
			}
		}
	}
	defer rollbackFn()

	if err := eh.Insert(c.Request.Context(), txn, boil.Infer()); err != nil {
		doRollback = true
		r.Logger.With(
			zap.Error(err),
			zap.String("event_id", evt.EventID.String()),
		).Warn("failed inserting event history")
		badRequestResponse(c, "failed inserting event", errors.Wrap(err,
			fmt.Sprintf("event %s", evt.EventID.String())))
		return
	}

	if err := txn.Commit(); err != nil {
		doRollback = true
		r.Logger.With(
			zap.Error(err),
			zap.String("event_id", evt.EventID.String()),
		).Warn("unable to commit transaction")
		metrics.DBError("commit event transaction")
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, evt.EventID.String())
}
