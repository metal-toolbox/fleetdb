//nolint:unused
package fleetdbapi

import (
	"encoding/json"
	"net/http"

	"github.com/bmc-toolbox/common"
	"github.com/gin-gonic/gin"
)

// A reminder for maintenance: this type needs to be able to contain all the
// relevant fields from Component-Inventory or Alloy.
type serverInventory struct {
	Inv        *common.Device    `json:"inventory"`
	BiosConfig map[string]string `json:"bios_config,omitempty"`
}

func (si *serverInventory) mustJSON() []byte {
	byt, err := json.Marshal(si)
	if err != nil {
		panic("bad inventory")
	}
	return byt
}

func (si *serverInventory) fromJSON(b []byte) error {
	return json.Unmarshal(b, si)
}

func unimplemented(c *gin.Context) {
	m := map[string]string{
		"err": "unimplemented",
	}
	c.JSON(http.StatusInternalServerError, m)
}

func (r *Router) getInbandInventory(c *gin.Context) {
	unimplemented(c)
}

func (r *Router) getOutofbandInventory(c *gin.Context) {
	unimplemented(c)
}

func (r *Router) setInbandInventory(c *gin.Context) {
	unimplemented(c)
}

func (r *Router) setOutofbandInventory(c *gin.Context) {
	unimplemented(c)
}
