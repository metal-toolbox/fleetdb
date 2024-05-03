package fleetdbapi

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

var (
	errBiosConfigSetRoute  = errors.New("error fullfilling config set request")
	errConfigListRoute = errors.New("error fullfilling config set list request")
	errNullRelation    = errors.New("sqlboiler relation was unexpectedly null")
)

func (r *Router) serverBiosConfigSetCreate(c *gin.Context) {
	var payload BiosConfigSet

	// Unmarshal JSON payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequestResponse(c, "invalid payload: BiosConfigSetCreate{}; failed to unmarshal config set", err)
		return
	}

	// Insert DBModel into DB
	id, err := r.insertBiosConfigSet(c.Request.Context(), &payload)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, id)
}

func (r *Router) serverBiosConfigSetGet(c *gin.Context) {
	// Get Config Set
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", errBiosConfigSetRoute)
		return
	}

	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.BiosConfigSetTableColumns.ID), id),
	}

	dbBiosConfigSet, err := r.eagerLoadBiosConfigSet(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// Convert to Marshallable struct
	var set BiosConfigSet
	err = set.fromDBModelBiosConfigSet(dbBiosConfigSet)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	itemResponse(c, set)
}

func (r *Router) serverBiosConfigSetDelete(c *gin.Context) {
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", errBiosConfigSetRoute)
	}

	set := &models.BiosConfigSet{}
	set.ID = id

	// Delete Config Set
	count, err := set.Delete(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse2(c, count)
}

func (r *Router) serverBiosConfigSetList(c *gin.Context) {
	params, err := parseBiosConfigSetListParams(c)
	if err != nil {
		badRequestResponse(c, "invalid query params", errConfigListRoute)
		return
	}

	mods := params.queryMods()

	count, err := models.BiosConfigSets().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	dbSets, err := r.eagerLoadAllBiosConfigSets(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	sets := make([]BiosConfigSet, len(dbSets))

	for i, dbSet := range dbSets {
		err = sets[i].fromDBModelBiosConfigSet(dbSet)
		if err != nil {
			dbErrorResponse(c, err)
			return
		}
	}

	pd := paginationData{
		pageCount:  len(sets),
		totalCount: count,
		pager:      params.Pagination,
	}

	listResponse(c, sets, pd)
}

func (r *Router) serverBiosConfigSetUpdate(c *gin.Context) {
	var payload BiosConfigSet

	// Get ID
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", errBiosConfigSetRoute)
	}

	// Unmarshal JSON payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequestResponse(c, "invalid payload: BiosConfigSetUpdate{}; failed to unmarshal config set", err)
		return
	}

	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.BiosConfigSetTableColumns.ID), id),
	}

	oldSet, err := r.eagerLoadBiosConfigSet(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse2(c, "failed to get config set that we want to update", err)
		return
	}

	// Insert DBModel into DB
	id, err = r.updateBiosConfigSet(c.Request.Context(), &payload, oldSet)
	if err != nil {
		dbErrorResponse2(c, "failed to update config set", err)
		return
	}

	updatedResponse(c, id)
}

func (r *Router) updateBiosConfigSet(ctx context.Context, set *BiosConfigSet, oldDBSet *models.BiosConfigSet) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", errors.Wrap(err, "0")
	}

	// If tx committed, rollback wont do anything
	// nolint:errcheck
	defer tx.Rollback()

	dbSet := set.toDBModelBiosConfigSet()
	dbSet.ID = set.ID

	_, err = dbSet.Update(ctx, tx, boil.Infer())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("IDs: %s", dbSet.ID))
	}

	var oldComponents []*models.BiosConfigComponent
	var oldSettings []*models.BiosConfigSetting
	var components []*models.BiosConfigComponent
	var settings []*models.BiosConfigSetting
	var settingsToDelete []*models.BiosConfigSetting
	var componentsToDelete []*models.BiosConfigComponent
	var componentsToUpdate []bool
	var settingsToUpdate [][]bool

	if oldDBSet.R != nil {
		oldComponents = oldDBSet.R.FKBiosConfigSetBiosConfigComponents
	}

	if dbSet.R != nil {
		components = dbSet.R.FKBiosConfigSetBiosConfigComponents
	}

	componentsToUpdate = make([]bool, len(components))
	settingsToUpdate = make([][]bool, len(components))

	// Gather information about what to delete, update, or insert
	for _, oldComponent := range oldComponents {
		componentFound := false
		for c, component := range components {
			if oldComponent.Name == component.Name {
				component.ID = oldComponent.ID
				component.FKBiosConfigSetID = dbSet.ID
				componentFound = true

				componentsToUpdate[c] = true

				if component.R != nil {
					settings = component.R.FKBiosConfigComponentBiosConfigSettings
				} else {
					settings = []*models.BiosConfigSetting{}
				}

				if oldComponent.R != nil {
					oldSettings = oldComponent.R.FKBiosConfigComponentBiosConfigSettings
				} else {
					oldSettings = []*models.BiosConfigSetting{}
				}

				settingsToUpdate[c] = make([]bool, len(settings))

				for _, oldSetting := range oldSettings {
					settingFound := false
					for s, setting := range settings {
						if oldSetting.SettingsKey == setting.SettingsKey {
							settingFound = true
							setting.ID = oldSetting.ID
							setting.FKBiosConfigComponentID = component.ID

							settingsToUpdate[c][s] = true
						}
					}

					if !settingFound {
						settingsToDelete = append(settingsToDelete, oldSetting)
					}
				}
			}
		}

		if !componentFound {
			componentsToDelete = append(componentsToDelete, oldComponent)
		}
	}

	// Delete components not found in new set
	for _, component := range componentsToDelete {
		_, err := component.Delete(ctx, tx) // Dont need to delete settings. CASCADE will handle that
		if err != nil {
			return "", err
		}
	}

	// Delete settings not found in updated components
	for _, setting := range settingsToDelete {
		_, err := setting.Delete(ctx, tx)
		if err != nil {
			return "", err
		}
	}

	// Insert/Update components
	for c, component := range components {
		if component.R == nil {
			return "", errNullRelation
		}

		err := component.R.FKBiosConfigSet.AddFKBiosConfigSetBiosConfigComponents(ctx, tx, !componentsToUpdate[c], component)
		if err != nil {
			return "", err
		}

		for s, setting := range components[c].R.FKBiosConfigComponentBiosConfigSettings {
			err = component.AddFKBiosConfigComponentBiosConfigSettings(ctx, tx, !settingsToUpdate[c][s], setting)
			if err != nil {
				return "", err
			}
		}
	}

	return dbSet.ID, tx.Commit()
}

func (r *Router) insertBiosConfigSet(ctx context.Context, set *BiosConfigSet) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	// If tx committed, rollback wont do anything
	// nolint:errcheck
	defer tx.Rollback()

	dbSet := set.toDBModelBiosConfigSet()

	err = dbSet.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return "", err
	}

	for _, component := range set.Components {
		dbComponent := component.toDBModelBiosConfigComponent()

		err = dbSet.AddFKBiosConfigSetBiosConfigComponents(ctx, tx, true, dbComponent)
		if err != nil {
			return "", err
		}

		for _, setting := range component.Settings {
			dbSetting := setting.toDBModelBiosConfigSetting()
			err = dbComponent.AddFKBiosConfigComponentBiosConfigSettings(ctx, tx, true, dbSetting)
			if err != nil {
				return "", err
			}
		}
	}

	return dbSet.ID, tx.Commit()
}

func (r *Router) eagerLoadBiosConfigSet(ctx context.Context, mods []qm.QueryMod) (*models.BiosConfigSet, error) {
	// Eager load relations
	mods = append(mods, qm.Load(models.BiosConfigSetRels.FKBiosConfigSetBiosConfigComponents))

	dbSet, err := models.BiosConfigSets(mods...).One(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	if dbSet.R != nil {
		for i := range dbSet.R.FKBiosConfigSetBiosConfigComponents {
			err := dbSet.R.FKBiosConfigSetBiosConfigComponents[i].L.LoadFKBiosConfigComponentBiosConfigSettings(ctx, r.DB, true, dbSet.R.FKBiosConfigSetBiosConfigComponents[i], nil)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, errNullRelation
	}

	return dbSet, nil
}

func (r *Router) eagerLoadAllBiosConfigSets(ctx context.Context, mods []qm.QueryMod) ([]*models.BiosConfigSet, error) {
	// Eager load relations
	mods = append(mods, qm.Load(models.BiosConfigSetRels.FKBiosConfigSetBiosConfigComponents))

	dbSets, err := models.BiosConfigSets(mods...).All(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	for _, dbSet := range dbSets {
		for i := range dbSet.R.FKBiosConfigSetBiosConfigComponents {
			err := dbSet.R.FKBiosConfigSetBiosConfigComponents[i].L.LoadFKBiosConfigComponentBiosConfigSettings(ctx, r.DB, true, dbSet.R.FKBiosConfigSetBiosConfigComponents[i], nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return dbSets, nil
}
