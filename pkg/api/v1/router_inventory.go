//nolint:unused
package fleetdbapi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	rivets "github.com/metal-toolbox/rivets/v2/types"
	"go.uber.org/zap"

	"github.com/metal-toolbox/fleetdb/internal/inventory"
)

func (r *Router) readInventoryFromDB(c *gin.Context) (*inventory.DeviceView, error) {
	srvID, err := r.parseUUID(c)
	if err != nil {
		badRequestResponse(c, "invalid server id", err)
		return nil, err
	}

	var doInband bool
	switch c.Query("mode") {
	case "inband":
		doInband = true
	case "outofband":
	default:
		badRequestResponse(c, "invalid inventory mode", nil)
		return nil, inventory.ErrBadInventoryMode
	}

	dv := &inventory.DeviceView{
		DeviceID: srvID,
		Inband:   doInband,
	}

	err = dv.FromDatastore(c.Request.Context(), r.DB)
	switch err {
	case nil:
	case inventory.ErrNoInventory:
		msg := fmt.Sprintf("no inventory for %s", srvID)
		notFoundResponse(c, msg)
		return nil, err
	default:
		dbErrorResponse(c, err)
		return nil, err
	}
	return dv, dv.FromDatastore(c.Request.Context(), r.DB)
}

func (r *Router) getInventory(c *gin.Context) {
	dv, err := r.readInventoryFromDB(c)
	if err != nil {
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

func compareComponent(srcComp, tgtComp interface{}) string {
	var builder strings.Builder
	valSrc := reflect.ValueOf(srcComp)
	valTgt := reflect.ValueOf(tgtComp)
	valSrc = valSrc.Elem()
	valTgt = valTgt.Elem()
	typeOfSrc := valSrc.Type()
	typeOfTgt := valTgt.Type()

	for i := 0; i < valSrc.NumField(); i++ {
		fieldSrc := valSrc.Field(i)
		fieldTgt := valTgt.Field(i)
		fieldName := typeOfSrc.Field(i).Name
		if !fieldTgt.IsValid() {
			builder.WriteString(fmt.Sprintf("- %s: %v\n", fieldName, fieldSrc.Interface()))
		} else if !reflect.DeepEqual(fieldSrc.Interface(), fieldTgt.Interface()) {
			builder.WriteString(fmt.Sprintf("- %s: %v\n", fieldName, fieldSrc.Interface()))
			builder.WriteString(fmt.Sprintf("+ %s: %v\n", fieldName, fieldTgt.Interface()))
		}
	}

	for i := 0; i < valTgt.NumField(); i++ {
		fieldTgt := valTgt.Field(i)
		fieldName := typeOfTgt.Field(i).Name

		if !valSrc.Field(i).IsValid() {
			builder.WriteString(fmt.Sprintf("+ %s: %v\n", fieldName, fieldTgt.Interface()))
		}
	}
	return builder.String()
}

func componentsToMap(cs []*rivets.Component) map[string][]*rivets.Component {
	theMap := make(map[string][]*rivets.Component)
	for _, c := range cs {
		name := c.Name
		// cSlice can be nil. Appending to nil is OK.
		cSlice := theMap[name]
		cSlice = append(cSlice, c)
		theMap[name] = cSlice
	}
	return theMap
}

func (r *Router) compareInventory(c *gin.Context) {
	dv, err := r.readInventoryFromDB(c)
	if err != nil {
		return
	}
	// compareComponents compares rivets.Components between two rivets.Server.
	srv := &rivets.Server{}
	if err := c.ShouldBindJSON(srv); err != nil {
		badRequestResponse(c, "invalid inventory payload", err)
		return
	}
	srcCompsMap := componentsToMap(srv.Components)
	tgtCompsMap := componentsToMap(dv.Inv.Components)
	var gaps string
	for _, tgtComps := range tgtCompsMap {
		slug := tgtComps[0].Name
		srcComps, ok := srcCompsMap[slug]
		if !ok {
			for _, tgtComp := range tgtComps {
				gaps += fmt.Sprintf("%v: %v\n", slug, compareComponent(&rivets.Component{}, tgtComp))
			}
			continue
		}

		for i, tgtComp := range tgtComps {
			if i >= len(srcComps) {
				gaps += fmt.Sprintf("%v: %v\n", slug, compareComponent(&rivets.Component{}, tgtComp))
				continue
			}
			val := reflect.ValueOf(tgtComp)
			if val.Kind() == reflect.Slice {
				// TODO: suupport comparing nested slices
			} else if val.Kind() == reflect.Ptr {
				gaps += fmt.Sprintf("%v: %v\n", slug, compareComponent(srcComps[i], tgtComp))
			}
		}
	}

	for _, srcComps := range srcCompsMap {
		slug := srcComps[0].Name
		_, ok := tgtCompsMap[slug]
		if !ok {
			for _, srcComp := range srcComps {
				gaps += fmt.Sprintf("%v: %v\n", slug, compareComponent(srcComp, &rivets.Component{}))
			}
			continue
		}
	}
	itemResponse(c, gaps)
}
