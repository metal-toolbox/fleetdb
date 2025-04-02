package fleetdbapi

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/metal-toolbox/rivets/v2/ginauth"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"gocloud.dev/secrets"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// Router provides a router for the v1 API
type Router struct {
	AuthMW        *ginauth.MultiTokenMiddleware
	DB            *sqlx.DB
	SecretsKeeper *secrets.Keeper
	Logger        *zap.Logger
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	amw := r.AuthMW

	// /servers
	srvs := rg.Group("/servers")
	{
		srvs.GET("", amw.AuthRequired(readScopes("server")), r.serverList)
		srvs.POST("", amw.AuthRequired(createScopes("server")), r.serverCreate)

		srvs.GET("/components", amw.AuthRequired(readScopes("server:component")), r.serverComponentList)

		// /servers/:uuid
		srv := srvs.Group("/:uuid")
		{
			srv.GET("", amw.AuthRequired(readScopes("server")), r.serverGet)
			srv.PUT("", amw.AuthRequired(updateScopes("server")), r.serverUpdate)
			srv.DELETE("", amw.AuthRequired(deleteScopes("server")), r.serverDelete)

			// /servers/:uuid/attributes
			srvAttrs := srv.Group("/attributes")
			{
				srvAttrs.GET("", amw.AuthRequired(readScopes("server", "server:attributes")), r.serverAttributesList)
				srvAttrs.POST("", amw.AuthRequired(createScopes("server", "server:attributes")), r.serverAttributesCreate)
				srvAttrs.GET("/:namespace", amw.AuthRequired(readScopes("server", "server:attributes")), r.serverAttributesGet)
				srvAttrs.PUT("/:namespace", amw.AuthRequired(updateScopes("server", "server:attributes")), r.serverAttributesUpdate)
				srvAttrs.DELETE("/:namespace", amw.AuthRequired(deleteScopes("server", "server:attributes")), r.serverAttributesDelete)
			}

			// /servers/:uuid/components
			srvComponents := srv.Group("/components")
			{
				srvComponents.POST("", amw.AuthRequired(createScopes("server", "server:component")), r.serverComponentsCreate)
				srvComponents.GET("", amw.AuthRequired(readScopes("server", "server:component")), r.serverComponentGet)
				srvComponents.PUT("", amw.AuthRequired(updateScopes("server", "server:component")), r.serverComponentUpdate)
				srvComponents.DELETE("", amw.AuthRequired(deleteScopes("server", "server:component")), r.serverComponentDelete)
			}

			// /servers/:uuid/credentials/:slug
			svrCreds := srv.Group("credentials/:slug")
			{
				svrCreds.GET("", amw.AuthRequired([]string{"read:server:credentials"}), r.serverCredentialGet)
				svrCreds.PUT("", amw.AuthRequired([]string{"write:server:credentials"}), r.serverCredentialUpsert)
				svrCreds.DELETE("", amw.AuthRequired([]string{"write:server:credentials"}), r.serverCredentialDelete)
			}

			// /servers/:uuid/versioned-attributes
			srvVerAttrs := srv.Group("/versioned-attributes")
			{
				srvVerAttrs.GET("", amw.AuthRequired(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesList)
				srvVerAttrs.POST("", amw.AuthRequired(createScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesCreate)
				srvVerAttrs.GET("/:namespace", amw.AuthRequired(readScopes("server", "server:versioned-attributes")), r.serverVersionedAttributesGet)
			}
		}
	}

	// /server-component-types
	srvCmpntType := rg.Group("/server-component-types")
	{
		srvCmpntType.GET("", amw.AuthRequired(readScopes("server-component-types")), r.serverComponentTypeList)
		srvCmpntType.POST("", amw.AuthRequired(updateScopes("server-component-types")), r.serverComponentTypeCreate)
	}

	// /server-component-firmwares
	srvCmpntFw := rg.Group("/server-component-firmwares")
	{
		srvCmpntFw.GET("", amw.AuthRequired(readScopes("server-component-firmwares")), r.serverComponentFirmwareList)
		srvCmpntFw.POST("", amw.AuthRequired(createScopes("server-component-firmwares")), r.serverComponentFirmwareCreate)
		srvCmpntFw.GET("/:uuid", amw.AuthRequired(readScopes("server-component-firmwares")), r.serverComponentFirmwareGet)
		srvCmpntFw.PUT("/:uuid", amw.AuthRequired(updateScopes("server-component-firmwares")), r.serverComponentFirmwareUpdate)
		srvCmpntFw.DELETE("/:uuid", amw.AuthRequired(deleteScopes("server-component-firmwares")), r.serverComponentFirmwareDelete)
	}

	// /server-credential-types
	srvCredentialTypes := rg.Group("/server-credential-types")
	{
		srvCredentialTypes.GET("", amw.AuthRequired(readScopes("server-credential-types")), r.serverCredentialTypesList)
		srvCredentialTypes.POST("", amw.AuthRequired(createScopes("server-credential-types")), r.serverCredentialTypesCreate)
	}

	// /server-component-firmware-sets
	srvCmpntFwSets := rg.Group("/server-component-firmware-sets")
	{
		createScopeMiddleware := amw.AuthRequired(createScopes("server-component-firmware-sets"))
		readScopeMiddleware := amw.AuthRequired(readScopes("server-component-firmware-sets"))
		updateScopeMiddleware := amw.AuthRequired(updateScopes("server-component-firmware-sets"))
		deleteScopeMiddleware := amw.AuthRequired(deleteScopes("server-component-firmware-sets"))

		// list all sets
		srvCmpntFwSets.GET("", readScopeMiddleware, r.serverComponentFirmwareSetList)

		// create/read/update/delete individual firmware sets
		srvCmpntFwSets.POST("", createScopeMiddleware, r.serverComponentFirmwareSetCreate)
		srvCmpntFwSets.GET("/:uuid", readScopeMiddleware, r.serverComponentFirmwareSetGet)
		srvCmpntFwSets.PUT("/:uuid", updateScopeMiddleware, r.serverComponentFirmwareSetUpdate)
		srvCmpntFwSets.DELETE("/:uuid", deleteScopeMiddleware, r.serverComponentFirmwareSetDelete)

		// remove a component firmware from the set
		srvCmpntFwSets.POST("/:uuid/remove-firmware", deleteScopeMiddleware, r.serverComponentFirmwareSetRemoveFirmware)

		// mark the set as validated
		srvCmpntFwSets.POST("/validate-firmware-set", updateScopeMiddleware, r.validateFirmwareSet)
	}

	// /bill-of-materials
	srvBoms := rg.Group("/bill-of-materials")
	{
		// /bill-of-materials/batch-boms-upload
		uploadFile := srvBoms.Group("/batch-upload")
		{
			uploadFile.POST("", amw.AuthRequired(createScopes("batch-upload")), r.bomsUpload)
		}

		// /bill-of-materials/aoc-mac-address
		srvBomByAocMacAddress := srvBoms.Group("/aoc-mac-address")
		{
			srvBomByAocMacAddress.GET("/:aoc_mac_address", amw.AuthRequired(readScopes("aoc-mac-address")), r.getBomFromAocMacAddress)
		}

		// /bill-of-materials/bmc-mac-address
		srvBomByBmcMacAddress := srvBoms.Group("/bmc-mac-address")
		{
			srvBomByBmcMacAddress.GET("/:bmc_mac_address", amw.AuthRequired(readScopes("bmc-mac-address")), r.getBomFromBmcMacAddress)
		}
	}

	// inventory endpoints
	srvInventory := rg.Group("/inventory")
	{
		// uuid is the server id
		srvInventory.GET("/:uuid", amw.AuthRequired(readScopes("server")), r.getInventory)
		srvInventory.PUT("/:uuid", amw.AuthRequired(updateScopes("server")), r.setInventory)
	}

	srvEvents := rg.Group("/events")
	{
		srvEvents.GET("/:evtID", amw.AuthRequired(readScopes("server")), r.getHistoryByConditionID)
		srvEvents.GET("/by-server/:srvID", amw.AuthRequired(readScopes("server")), r.getServerEvents)
		srvEvents.PUT("/:evtID", amw.AuthRequired(updateScopes("server")), r.updateEvent)
	}

	// /server-bios-config-sets
	srvCfgSets := rg.Group("/server-bios-config-sets")
	{
		srvCfgSets.GET("", amw.AuthRequired(readScopes("server-bios-config-sets")), r.serverBiosConfigSetList)
		srvCfgSets.POST("", amw.AuthRequired(createScopes("server-bios-config-sets")), r.serverBiosConfigSetCreate)
		srvCfgSets.GET("/:uuid", amw.AuthRequired(readScopes("server-bios-config-sets")), r.serverBiosConfigSetGet)
		srvCfgSets.PUT("/:uuid", amw.AuthRequired(updateScopes("server-bios-config-sets")), r.serverBiosConfigSetUpdate)
		srvCfgSets.DELETE("/:uuid", amw.AuthRequired(deleteScopes("server-bios-config-sets")), r.serverBiosConfigSetDelete)
	}
}

func createScopes(items ...string) []string {
	s := []string{"write", "create"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("create:%s", i))
	}

	return s
}

func readScopes(items ...string) []string {
	s := []string{"read"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("read:%s", i))
	}

	return s
}

func updateScopes(items ...string) []string {
	s := []string{"write", "update"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("update:%s", i))
	}

	return s
}

func deleteScopes(items ...string) []string {
	s := []string{"write", "delete"}
	for _, i := range items {
		s = append(s, fmt.Sprintf("delete:%s", i))
	}

	return s
}

func (r *Router) parseUUID(c *gin.Context) (uuid.UUID, error) {
	u, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		badRequestResponse(c, "failed to parse uuid", err)
	}

	return u, err
}

func (r *Router) loadServerFromParams(c *gin.Context) (*models.Server, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, errors.Wrap(ErrUUIDParse, err.Error())
	}

	srv, err := models.FindServer(c.Request.Context(), r.DB, u.String())
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (r *Router) loadOrCreateServerFromParams(c *gin.Context) (*models.Server, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	srv, err := models.FindServer(c.Request.Context(), r.DB, u.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv = &models.Server{ID: u.String()}
			if err := srv.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
				dbErrorResponse(c, err)
				return nil, err
			}

			return srv, nil
		}

		dbErrorResponse(c, err)

		return nil, err
	}

	return srv, nil
}

func (r *Router) loadComponentFirmwareVersionFromParams(c *gin.Context) (*models.ComponentFirmwareVersion, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	firmware, err := models.FindComponentFirmwareVersion(c.Request.Context(), r.DB, u.String())
	if err != nil {
		dbErrorResponse(c, err)

		return nil, err
	}

	return firmware, nil
}
