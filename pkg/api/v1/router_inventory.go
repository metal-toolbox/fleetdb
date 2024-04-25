//nolint:unused
package fleetdbapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	rivets "github.com/metal-toolbox/rivets/types"
	"go.uber.org/zap"

	"github.com/metal-toolbox/fleetdb/internal/inventory"
)

func (r *Router) getInventory(c *gin.Context) {
	srvID, err := r.parseUUID(c)
	if err != nil {
		badRequestResponse(c, "invalid server id", err)
		return
	}

	var doInband bool
	switch c.Query("mode") {
	case "inband":
		doInband = true
	case "outofband":
	default:
		badRequestResponse(c, "invalid inventory mode", nil)
		return
	}

	dv := inventory.DeviceView{
		DeviceID: srvID,
		Inband:   doInband,
	}

	err = dv.FromDatastore(c.Request.Context(), r.DB)
	switch err {
	case nil:
	case inventory.ErrNoInventory:
		msg := fmt.Sprintf("no inventory for %s", srvID)
		notFoundResponse(c, msg)
		return
	default:
		dbErrorResponse(c, err)
		return
	}
	itemResponse(c, dv.Inv)
}

func (r *Router) setInventory(c *gin.Context) {
	srvID, err := r.parseUUID(c)
	if err != nil {
		badRequestResponse(c, "invalid server id", err)
		return
	}

	var doInband bool
	param := c.Query("mode")

	switch param {
	case "inband":
		doInband = true
	case "outofband":
	case "":
		badRequestResponse(c, "missing inventory specifier", nil)
		return
	default:
		badRequestResponse(c, fmt.Sprintf("invalid inventory mode: %s", param), nil)
		return
	}

	srv := &rivets.Server{}
	if err := c.ShouldBindJSON(srv); err != nil {
		badRequestResponse(c, "invalid inventory payload", err)
		return
	}

	if err := inventory.ServerSanityCheck(srv); err != nil {
		badRequestResponse(c, "invalid inventory payload", err)
		return
	}

	view := &inventory.DeviceView{
		DeviceID: srvID,
		Inband:   doInband,
		Inv:      srv,
	}

	txn := r.DB.MustBegin()

	if err := view.UpsertInventory(c.Request.Context(), txn); err != nil {
		if err := txn.Rollback(); err != nil {
			r.Logger.With(
				zap.Error(err),
				zap.String("device_id", srvID.String()),
				zap.Bool("inband", doInband),
			).Warn("rollback error")
			// increment error metrics
		}
		dbErrorResponse(c, err)
	}
	if err := txn.Commit(); err != nil {
		r.Logger.With(
			zap.Error(err),
			zap.String("device_id", srvID.String()),
			zap.Bool("inband", doInband),
		).Warn("commit error")
		// increment error metrics
		dbErrorResponse(c, err)
	}
	updatedResponse(c, "")
}
