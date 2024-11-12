package dbtools

import (
	"context"
	"database/sql"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	common "github.com/metal-toolbox/bmc-common"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

var (
	errAddTypes = errors.New("unable to add component types")
	errBadName  = errors.New("unknown component name")
)

// XXX: if bmc-toolbox exported this as a list, we could just import it.
var knownComponentTypes = []string{
	common.SlugBackplaneExpander,
	common.SlugChassis,
	common.SlugTPM,
	common.SlugGPU,
	common.SlugCPU,
	common.SlugPhysicalMem,
	common.SlugStorageController,
	common.SlugBMC,
	common.SlugBIOS,
	common.SlugDrive,
	common.SlugDriveTypePCIeNVMEeSSD,
	common.SlugDriveTypeSATASSD,
	common.SlugDriveTypeSATAHDD,
	common.SlugNIC,
	common.SlugPSU,
	common.SlugCPLD,
	common.SlugEnclosure,
	common.SlugUnknown,
	common.SlugMainboard,
}

var componentTypeIDCache = map[string]string{}

// SetupComponentTypes upserts all known component types to the database.
// Despite the descriptor, in the database the Name field of the component type
// is the verbatim value of the string, and the Slug is computed as a lower-case
// english-localized variant.
func SetupComponentTypes(ctx context.Context, db *sqlx.DB) error {
	txn := db.MustBeginTx(ctx, &sql.TxOptions{})
	for _, typ := range knownComponentTypes {
		existing, err := models.ServerComponentTypes(
			models.ServerComponentTypeWhere.Name.EQ(typ),
		).One(ctx, txn)

		switch err {
		case nil:
			componentTypeIDCache[typ] = existing.ID
		case sql.ErrNoRows:
			sct := &models.ServerComponentType{
				Name: typ,
				Slug: slug.Make(typ),
			}
			if err := sct.Insert(ctx, txn, boil.Infer()); err != nil {
				_ = txn.Rollback()
				return errors.Wrap(errAddTypes, err.Error())
			}
			componentTypeIDCache[typ] = sct.ID
		default:
			return errors.Wrap(errAddTypes, err.Error())
		}
	}
	return txn.Commit()
}

// ComponentTypeIDFromName expects the name of the component (as defined in
// bmc-toolbox) and will return the internal database ID for that name.
func ComponentTypeIDFromName(name string) (string, error) {
	id, ok := componentTypeIDCache[name]
	if !ok {
		return "", errors.Wrap(errBadName, name)
	}
	return id, nil
}
