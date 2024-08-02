package fleetdbapi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

func (r *Router) serverSkuCreate(c *gin.Context) {
	var payload ServerSku

	// Unmarshal JSON payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequestResponse(c, "invalid payload; failed to unmarshal sku", err)
	}

	// Insert DBModel into DB
	id, err := r.insertServerSku(c.Request.Context(), &payload)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, id)
}

func (r *Router) serverSkuGet(c *gin.Context) {
	// Get ID
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", ErrRouteServerSku)
		return
	}

	// Setup query
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ServerSkuColumns.ID), id),
	}

	// Get Server Sku
	dbServerSku, err := r.eagerLoadServerSku(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// Convert to Marshallable struct
	var sku ServerSku
	sku.fromDBModelServerSku(dbServerSku)

	itemResponse(c, sku)
}

func (r *Router) serverSkuUpdate(c *gin.Context) {
	var payload ServerSku

	// Get ID
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", ErrRouteServerSku)
	}

	// Unmarshal JSON payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequestResponse(c, "invalid payload; failed to unmarshal sku", err)
		return
	}

	// Setup query
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ServerSkuColumns.ID), id),
	}

	// Get Current Server Sku
	oldDBServerSku, err := r.eagerLoadServerSku(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	newDBServerSku := payload.toDBModelServerSkuDeep()

	// Insert DBModel into DB
	id, err = r.updateServerSkuTransaction(c.Request.Context(), newDBServerSku, oldDBServerSku)
	if err != nil {
		dbErrorResponse2(c, "failed to update server sku", err)
		return
	}

	updatedResponse(c, id)
}

func (r *Router) serverSkuDelete(c *gin.Context) {
	// Get ID
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", ErrRouteServerSku)
	}

	set := &models.ServerSku{}
	set.ID = id

	// Delete Config Set
	count, err := set.Delete(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse2(c, count)
}

func (r *Router) serverSkuList(c *gin.Context) {
	params, err := parseServerSkuListParams(c)
	if err != nil {
		badRequestResponse(c, "invalid query params", ErrRouteBiosConfigSet)
		return
	}

	mods := params.queryMods()

	count, err := models.BiosConfigSets().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	dbSkus, err := r.eagerLoadAllServerSku(c.Request.Context(), mods, params.Pagination.Preload)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	skus := make([]ServerSku, len(dbSkus))

	for i, dbSku := range dbSkus {
		skus[i].fromDBModelServerSku(dbSku)
		if err != nil {
			dbErrorResponse(c, err)
			return
		}
	}

	pd := paginationData{
		pageCount:  len(skus),
		totalCount: count,
		pager:      params.Pagination,
	}

	listResponse(c, skus, pd)
}

func (r *Router) eagerLoadServerSku(ctx context.Context, mods []qm.QueryMod) (*models.ServerSku, error) {
	// Include all relations
	mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuAuxDevices))
	mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuDisks))
	mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuMemories))
	mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuNics))

	// Execute query
	dbSku, err := models.ServerSkus(mods...).One(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	return dbSku, nil
}

func (r *Router) eagerLoadAllServerSku(ctx context.Context, mods []qm.QueryMod, preload bool) ([]*models.ServerSku, error) {
	// Eager load relations
	if preload {
		mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuAuxDevices))
		mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuDisks))
		mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuMemories))
		mods = append(mods, qm.Load(models.ServerSkuRels.SkuServerSkuNics))
	}

	// Execute query
	dbSku, err := models.ServerSkus(mods...).All(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	return dbSku, nil
}

func (r *Router) insertServerSku(ctx context.Context, sku *ServerSku) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	defer loggedRollback(r, tx)

	dbSku := sku.toDBModelServerSku()

	err = dbSku.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return "", err
	}

	// Aux Devices
	dbAuxDevice := make([]*models.ServerSkuAuxDevice, len(sku.AuxDevices))
	for i, auxDevice := range sku.AuxDevices {
		dbAuxDevice[i] = auxDevice.toDBModelServerSkuAuxDevice()
	}
	err = dbSku.AddSkuServerSkuAuxDevices(ctx, tx, true, dbAuxDevice...)
	if err != nil {
		return "", err
	}

	// Disks
	dbDisks := make([]*models.ServerSkuDisk, len(sku.Disks))
	for i, disk := range sku.Disks {
		dbDisks[i] = disk.toDBModelServerSkuDisk()
	}
	err = dbSku.AddSkuServerSkuDisks(ctx, tx, true, dbDisks...)
	if err != nil {
		return "", err
	}

	// Memory
	dbMemory := make([]*models.ServerSkuMemory, len(sku.Memory))
	for i, memory := range sku.Memory {
		dbMemory[i] = memory.toDBModelServerSkuMemory()
	}
	err = dbSku.AddSkuServerSkuMemories(ctx, tx, true, dbMemory...)
	if err != nil {
		return "", err
	}

	// Nics
	dbNics := make([]*models.ServerSkuNic, len(sku.Nics))
	for i, nic := range sku.Nics {
		dbNics[i] = nic.toDBModelServerSkuNic()
	}
	err = dbSku.AddSkuServerSkuNics(ctx, tx, true, dbNics...)
	if err != nil {
		return "", err
	}

	return dbSku.ID, tx.Commit()
}

func (r *Router) updateServerSkuTransaction(ctx context.Context, sku *models.ServerSku, oldSku *models.ServerSku) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	defer loggedRollback(r, tx)

	_, err = sku.Update(ctx, tx, boil.Infer())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("ID: %s", sku.ID))
	}

	err = r.updateServerSkuAuxDevices(ctx, tx, sku, oldSku)
	if err != nil {
		return "", err
	}

	err = r.updateServerSkuDisks(ctx, tx, sku, oldSku)
	if err != nil {
		return "", err
	}

	err = r.updateServerSkuMemory(ctx, tx, sku, oldSku)
	if err != nil {
		return "", err
	}

	err = r.updateServerSkuNics(ctx, tx, sku, oldSku)
	if err != nil {
		return "", err
	}

	return sku.ID, tx.Commit()
}

// TODO; ADD COMMENT TO PR
// These functions are why I believe C Macros are a net benefit. Using generics wouldnt work without also using reflection to find the relation structs.
// I could also be doing to much work, and maybe there is a simpler way of doing this?

func (r *Router) updateServerSkuAuxDevices(ctx context.Context, tx *sql.Tx, sku *models.ServerSku, oldSku *models.ServerSku) error {
	var oldAuxDevices []*models.ServerSkuAuxDevice
	var auxDevices []*models.ServerSkuAuxDevice

	if oldSku.R != nil {
		oldAuxDevices = oldSku.R.SkuServerSkuAuxDevices
	}

	if sku.R != nil {
		auxDevices = sku.R.SkuServerSkuAuxDevices
	}

	// Find aux devices no longer present and remove them
	for _, oldAuxDevice := range oldAuxDevices {
		auxDeviceFound := false
		for _, auxDevice := range auxDevices {
			if auxDevice.ID == oldAuxDevice.ID {
				auxDeviceFound = true
				break
			}
		}

		if !auxDeviceFound {
			_, err := oldAuxDevice.Delete(ctx, tx)
			if err != nil {
				return err
			}
		}
	}

	// Upsert aux devices
	for _, auxDevice := range auxDevices {
		err := auxDevice.Upsert(ctx, tx, true,
			[]string{models.ServerSkuAuxDeviceColumns.ID},
			boil.Whitelist(
				models.ServerSkuAuxDeviceColumns.Vendor,
				models.ServerSkuAuxDeviceColumns.Model,
				models.ServerSkuAuxDeviceColumns.DeviceType,
				models.ServerSkuAuxDeviceColumns.Details,
				models.ServerSkuAuxDeviceColumns.UpdatedAt,
			),
			boil.Whitelist(
				models.ServerSkuAuxDeviceColumns.ID,
				models.ServerSkuAuxDeviceColumns.SkuID,
				models.ServerSkuAuxDeviceColumns.Vendor,
				models.ServerSkuAuxDeviceColumns.Model,
				models.ServerSkuAuxDeviceColumns.DeviceType,
				models.ServerSkuAuxDeviceColumns.Details,
				models.ServerSkuAuxDeviceColumns.CreatedAt,
				models.ServerSkuAuxDeviceColumns.UpdatedAt,
			))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Router) updateServerSkuDisks(ctx context.Context, tx *sql.Tx, sku *models.ServerSku, oldSku *models.ServerSku) error {
	var oldDisks []*models.ServerSkuDisk
	var disks []*models.ServerSkuDisk

	if oldSku.R != nil {
		oldDisks = oldSku.R.SkuServerSkuDisks
	}

	if sku.R != nil {
		disks = sku.R.SkuServerSkuDisks
	}

	// Find disks no longer present and remove them
	for _, oldDisk := range oldDisks {
		diskFound := false
		for _, disk := range disks {
			if disk.ID == oldDisk.ID {
				diskFound = true
				break
			}
		}

		if !diskFound {
			_, err := oldDisk.Delete(ctx, tx)
			if err != nil {
				return err
			}
		}
	}

	// Upsert disks
	for _, disk := range disks {
		err := disk.Upsert(ctx, tx, true,
			[]string{models.ServerSkuDiskColumns.ID},
			boil.Whitelist(
				models.ServerSkuDiskColumns.Bytes,
				models.ServerSkuDiskColumns.Protocol,
				models.ServerSkuDiskColumns.Count,
				models.ServerSkuDiskColumns.UpdatedAt,
			),
			boil.Whitelist(
				models.ServerSkuDiskColumns.ID,
				models.ServerSkuDiskColumns.SkuID,
				models.ServerSkuDiskColumns.Bytes,
				models.ServerSkuDiskColumns.Protocol,
				models.ServerSkuDiskColumns.Count,
				models.ServerSkuDiskColumns.CreatedAt,
				models.ServerSkuDiskColumns.UpdatedAt,
			))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Router) updateServerSkuMemory(ctx context.Context, tx *sql.Tx, sku *models.ServerSku, oldSku *models.ServerSku) error {
	var oldMemory []*models.ServerSkuMemory
	var memory []*models.ServerSkuMemory

	if oldSku.R != nil {
		oldMemory = oldSku.R.SkuServerSkuMemories
	}

	if sku.R != nil {
		memory = sku.R.SkuServerSkuMemories
	}

	// Find memory no longer present and remove them
	for _, oldMemoryItem := range oldMemory {
		memoryFound := false
		for _, memoryItem := range memory {
			if memoryItem.ID == oldMemoryItem.ID {
				memoryFound = true
				break
			}
		}

		if !memoryFound {
			_, err := oldMemoryItem.Delete(ctx, tx)
			if err != nil {
				return err
			}
		}
	}

	// Upsert memory
	for _, memoryItem := range memory {
		err := memoryItem.Upsert(ctx, tx, true,
			[]string{models.ServerSkuMemoryColumns.ID},
			boil.Whitelist(
				models.ServerSkuMemoryColumns.Bytes,
				models.ServerSkuMemoryColumns.Count,
				models.ServerSkuMemoryColumns.UpdatedAt,
			),
			boil.Whitelist(
				models.ServerSkuMemoryColumns.ID,
				models.ServerSkuMemoryColumns.SkuID,
				models.ServerSkuMemoryColumns.Bytes,
				models.ServerSkuMemoryColumns.Count,
				models.ServerSkuMemoryColumns.CreatedAt,
				models.ServerSkuMemoryColumns.UpdatedAt,
			))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Router) updateServerSkuNics(ctx context.Context, tx *sql.Tx, sku *models.ServerSku, oldSku *models.ServerSku) error {
	var oldNics []*models.ServerSkuNic
	var nics []*models.ServerSkuNic

	if oldSku.R != nil {
		oldNics = oldSku.R.SkuServerSkuNics
	}

	if sku.R != nil {
		nics = sku.R.SkuServerSkuNics
	}

	// Find nics no longer present and remove them
	for _, oldNic := range oldNics {
		nicFound := false
		for _, nic := range nics {
			if nic.ID == oldNic.ID {
				nicFound = true
				break
			}
		}

		if !nicFound {
			_, err := oldNic.Delete(ctx, tx)
			if err != nil {
				return err
			}
		}
	}

	// Upsert nics
	for _, nic := range nics {
		err := nic.Upsert(ctx, tx, true,
			[]string{models.ServerSkuNicColumns.ID},
			boil.Whitelist(
				models.ServerSkuNicColumns.PortBandwidth,
				models.ServerSkuNicColumns.PortCount,
				models.ServerSkuNicColumns.Count,
				models.ServerSkuNicColumns.UpdatedAt,
			),
			boil.Whitelist(
				models.ServerSkuNicColumns.ID,
				models.ServerSkuNicColumns.SkuID,
				models.ServerSkuNicColumns.PortBandwidth,
				models.ServerSkuNicColumns.PortCount,
				models.ServerSkuNicColumns.Count,
				models.ServerSkuNicColumns.CreatedAt,
				models.ServerSkuNicColumns.UpdatedAt,
			))
		if err != nil {
			return err
		}
	}

	return nil
}
