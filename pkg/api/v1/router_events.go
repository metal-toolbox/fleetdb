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
	EventID     uuid.UUID       `json:"event_id"`
	Type        string          `json:"event_type"`
	Start       time.Time       `json:"event_start"`
	End         time.Time       `json:"event_end"`
	Target      uuid.UUID       `json:"target_server"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	FinalState  string          `json:"final_state"`
	FinalStatus json.RawMessage `json:"final_status,omitempty"`
}

func (r *Router) getEventByID(c *gin.Context) {
	evtID, err := uuid.Parse(c.Param("evtID"))
	if err != nil {
		badRequestResponse(c, "failed to parse event id", err)
		return
	}

	eh, err := models.EventHistories(
		models.EventHistoryWhere.EventID.EQ(evtID.String()),
	).One(c.Request.Context(), r.DB)

	switch err {
	case nil:
	case sql.ErrNoRows:
		msg := fmt.Sprintf("no event for id %s", evtID.String())
		notFoundResponse(c, msg)
		return
	default:
		metrics.DBError("fetching event history")
		r.Logger.With(
			zap.Error(err),
			zap.String("event_id", evtID.String()),
		).Warn("event history by id")
		dbErrorResponse(c, err)
		return
	}

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

	itemResponse(c, evt)
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
		qm.OrderBy("event_end DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(c.Request.Context(), r.DB)

	r.Logger.With(
		zap.Error(err),
		zap.Int("record_count", len(ehs)),
		zap.String("server_id", srvID.String()),
	).Debug("retrieved event history")

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

	r.Logger.With(
		zap.Int("event_count", len(evts)),
		zap.String("server_id", srvID.String()),
	).Debug("returning events")

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
	).One(ctx, r.DB)

	switch err {
	case sql.ErrNoRows:
	case nil:
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
