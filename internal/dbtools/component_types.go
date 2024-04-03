package dbtools

import (
	"context"
	"database/sql"

	"github.com/bmc-toolbox/common"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

var errAddTypes = errors.New("unable to add component types")

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

// SetupComponentTypes upserts all known component types to the database.
// Despite the descriptor, in the database the Name field of the component type
// is the verbatim value of the string, and the Slug is computed as a lower-case
// english-localized variant.
func SetupComponentTypes(ctx context.Context, db *sqlx.DB) error {
	txn := db.MustBeginTx(ctx, &sql.TxOptions{})
	for _, typ := range knownComponentTypes {
		sct := &models.ServerComponentType{
			Name: typ,
			Slug: slug.Make(typ),
		}
		if err := sct.Upsert(ctx, txn, false, []string{"slug"}, boil.None(), boil.Infer()); err != nil {
			_ = txn.Rollback()
			return errors.Wrap(errAddTypes, err.Error())
		}
	}
	return txn.Commit()
}

// ComponentTypeIDFromName expects the name of the component (as defined in
// bmc-toolbox) and will return the internal database ID for that name.
func ComponentTypeIDFromName(ctx context.Context, exec boil.ContextExecutor, name string) (string, error) {
	sct, err := models.ServerComponentTypes(
		models.ServerComponentTypeWhere.Name.EQ(name),
	).One(ctx, exec)
	if err != nil {
		return "", err
	}
	return sct.ID, nil
}

// MustComponentTypeID returns the component type id for the given component name or panics
func MustComponentTypeID(ctx context.Context, exec boil.ContextExecutor, name string) string {
	id, err := ComponentTypeIDFromName(ctx, exec, name)
	if err != nil {
		panic(err)
	}
	return id
}
