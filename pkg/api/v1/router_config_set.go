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
	errConfigSetRoute  = errors.New("error fullfilling config set request")
	errConfigListRoute = errors.New("error fullfilling config set list request")
	errNullRelation    = errors.New("sqlboiler relation was unexpectedly null")
)

func (r *Router) serverConfigSetCreate(c *gin.Context) {
	var payload ConfigSet

	// Unmarshal JSON payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequestResponse(c, "invalid payload: ConfigSetCreate{}; failed to unmarshal config set", err)
		return
	}

	// Insert DBModel into DB
	id, err := r.insertConfigSet(c.Request.Context(), &payload)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, id)
}

func (r *Router) serverConfigSetGet(c *gin.Context) {
	// Get Config Set
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", errConfigSetRoute)
		return
	}

	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ConfigSetTableColumns.ID), id),
	}

	dbConfigSet, err := r.eagerLoadConfigSet(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// Convert to Marshallable struct
	var set ConfigSet
	err = set.fromDBModelConfigSet(dbConfigSet)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	itemResponse(c, set)
}

func (r *Router) serverConfigSetDelete(c *gin.Context) {
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", errConfigSetRoute)
	}

	set := &models.ConfigSet{}
	set.ID = id

	// Delete Config Set
	count, err := set.Delete(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse2(c, count)
}

func (r *Router) serverConfigSetList(c *gin.Context) {
	params, err := parseConfigSetListParams(c)
	if err != nil {
		badRequestResponse(c, "invalid query params", errConfigListRoute)
		return
	}

	mods := params.queryMods()

	count, err := models.ConfigSets().Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	dbSets, err := r.eagerLoadAllConfigSets(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	sets := make([]ConfigSet, len(dbSets))

	for i, dbSet := range dbSets {
		err = sets[i].fromDBModelConfigSet(dbSet)
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

func (r *Router) serverConfigSetUpdate(c *gin.Context) {
	var payload ConfigSet

	// Get ID
	id := c.Param("uuid")
	if id == "" || id == uuid.Nil.String() {
		badRequestResponse(c, "no UUID query param", errConfigSetRoute)
	}

	// Unmarshal JSON payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		badRequestResponse(c, "invalid payload: ConfigSetUpdate{}; failed to unmarshal config set", err)
		return
	}

	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ConfigSetTableColumns.ID), id),
	}

	oldSet, err := r.eagerLoadConfigSet(c.Request.Context(), mods)
	if err != nil {
		dbErrorResponse2(c, "failed to get config set that we want to update", err)
		return
	}

	// Insert DBModel into DB
	id, err = r.updateConfigSet(c.Request.Context(), &payload, oldSet)
	if err != nil {
		dbErrorResponse2(c, "failed to update config set", err)
		return
	}

	updatedResponse(c, id)
}

func (r *Router) updateConfigSet(ctx context.Context, set *ConfigSet, oldDBSet *models.ConfigSet) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", errors.Wrap(err, "0")
	}

	// If tx committed, rollback wont do anything
	// nolint:errcheck
	defer tx.Rollback()

	dbSet := set.toDBModelConfigSet()
	dbSet.ID = set.ID

	_, err = dbSet.Update(ctx, tx, boil.Infer())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("IDs: %s", dbSet.ID))
	}

	var oldComponents []*models.ConfigComponent
	var oldSettings []*models.ConfigComponentSetting
	var components []*models.ConfigComponent
	var settings []*models.ConfigComponentSetting
	var settingsToDelete []*models.ConfigComponentSetting
	var componentsToDelete []*models.ConfigComponent
	var componentsToUpdate []bool
	var settingsToUpdate [][]bool

	if oldDBSet.R != nil {
		oldComponents = oldDBSet.R.FKConfigSetConfigComponents
	}

	if dbSet.R != nil {
		components = dbSet.R.FKConfigSetConfigComponents
	}

	componentsToUpdate = make([]bool, len(components))
	settingsToUpdate = make([][]bool, len(components))

	// Gather information about what to delete, update, or insert
	for _, oldComponent := range oldComponents {
		componentFound := false
		for c, component := range components {
			if oldComponent.Name == component.Name {
				component.ID = oldComponent.ID
				component.FKConfigSetID = dbSet.ID
				componentFound = true

				componentsToUpdate[c] = true

				if component.R != nil {
					settings = component.R.FKComponentConfigComponentSettings
				} else {
					settings = []*models.ConfigComponentSetting{}
				}

				if oldComponent.R != nil {
					oldSettings = oldComponent.R.FKComponentConfigComponentSettings
				} else {
					oldSettings = []*models.ConfigComponentSetting{}
				}

				settingsToUpdate[c] = make([]bool, len(settings))

				for _, oldSetting := range oldSettings {
					settingFound := false
					for s, setting := range settings {
						if oldSetting.SettingsKey == setting.SettingsKey {
							settingFound = true
							setting.ID = oldSetting.ID
							setting.FKComponentID = component.ID

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

		err := component.R.FKConfigSet.AddFKConfigSetConfigComponents(ctx, tx, !componentsToUpdate[c], component)
		if err != nil {
			return "", err
		}

		for s, setting := range components[c].R.FKComponentConfigComponentSettings {
			err = component.AddFKComponentConfigComponentSettings(ctx, tx, !settingsToUpdate[c][s], setting)
			if err != nil {
				return "", err
			}
		}
	}

	return dbSet.ID, tx.Commit()
}

func (r *Router) insertConfigSet(ctx context.Context, set *ConfigSet) (string, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	// If tx committed, rollback wont do anything
	// nolint:errcheck
	defer tx.Rollback()

	dbSet := set.toDBModelConfigSet()

	err = dbSet.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return "", err
	}

	for _, component := range set.Components {
		dbComponent := component.toDBModelConfigComponent()

		err = dbSet.AddFKConfigSetConfigComponents(ctx, tx, true, dbComponent)
		if err != nil {
			return "", err
		}

		for _, setting := range component.Settings {
			dbSetting := setting.toDBModelConfigComponentSetting()
			err = dbComponent.AddFKComponentConfigComponentSettings(ctx, tx, true, dbSetting)
			if err != nil {
				return "", err
			}
		}
	}

	return dbSet.ID, tx.Commit()
}

func (r *Router) eagerLoadConfigSet(ctx context.Context, mods []qm.QueryMod) (*models.ConfigSet, error) {
	// Eager load relations
	mods = append(mods, qm.Load(models.ConfigSetRels.FKConfigSetConfigComponents))

	dbSet, err := models.ConfigSets(mods...).One(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	if dbSet.R != nil {
		for i := range dbSet.R.FKConfigSetConfigComponents {
			err := dbSet.R.FKConfigSetConfigComponents[i].L.LoadFKComponentConfigComponentSettings(ctx, r.DB, true, dbSet.R.FKConfigSetConfigComponents[i], nil)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, errNullRelation
	}

	return dbSet, nil
}

func (r *Router) eagerLoadAllConfigSets(ctx context.Context, mods []qm.QueryMod) ([]*models.ConfigSet, error) {
	// Eager load relations
	mods = append(mods, qm.Load(models.ConfigSetRels.FKConfigSetConfigComponents))

	dbSets, err := models.ConfigSets(mods...).All(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	for _, dbSet := range dbSets {
		for i := range dbSet.R.FKConfigSetConfigComponents {
			err := dbSet.R.FKConfigSetConfigComponents[i].L.LoadFKComponentConfigComponentSettings(ctx, r.DB, true, dbSet.R.FKConfigSetConfigComponents[i], nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return dbSets, nil
}
