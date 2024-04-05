//nolint:unused
package fleetdbapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/metal-toolbox/fleetdb/internal/inventory"
)

func unimplemented(c *gin.Context) {
	m := map[string]string{
		"err": "unimplemented",
	}
	c.JSON(http.StatusInternalServerError, m)
}

func (r *Router) getInventory(c *gin.Context) {
	unimplemented(c)
}

func (r *Router) setInventory(c *gin.Context) {
	srvID, err := r.parseUUID(c)
	if err != nil {
		badRequestResponse(c, "invalid server id", err)
		return
	}

	var doInband bool
	switch c.Param("mode") {
	case "inband":
		doInband = true
	case "outofband":
	default:
		badRequestResponse(c, "invalid inventory mode", nil)
	}

	view := inventory.DeviceView{}
	if err := c.ShouldBindJSON(&view); err != nil {
		badRequestResponse(c, "invalid inventory payload", err)
		return
	}

	view.DeviceID = srvID
	view.Inband = doInband

	txn := r.DB.MustBegin()

	// XXX what about BIOS config?
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
}

/*(	if err != nil {
		tx.Rollback() //nolint errcheck
		dbErrorResponse(c, err)

		return
	}
	// compose the attributes from the inventory
	// - server vendor attributes
	// - server metadata attributes

	// compose the components

}

func (r *Router) setOutofbandInventory(c *gin.Context) {
	unimplemented(c)
}*/
