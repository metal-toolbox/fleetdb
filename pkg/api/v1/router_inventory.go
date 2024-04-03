//nolint:unused
package fleetdbapi

import (
	"net/http"

	"github.com/metal-toolbox/fleetdb/internal/inventory"

	"github.com/gin-gonic/gin"
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

	var view inventory.DeviceView
	if err := c.ShouldBindJSON(&view); err != nil {
		badRequestResponse(c, "invalid inventory payload", err)
		return
	}

	if err := view.UpsertInventory(c.Request.Context(), r.DB, srvID, doInband); err != nil {
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
