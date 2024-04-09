//nolint:all  // XXX remove this!
package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bmc-toolbox/common"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	inbandComponentNamespace    = "sh.hollow.alloy.inband.metadata"
	outofbandComponentNamespace = "sh.hollow.alloy.outofband.metadata"

	errComponent     = errors.New("component error")
	errAttribute     = errors.New("attribute error")
	errVersionedAttr = errors.New("versioned attribute error")
)

func createOrUpdateComponent(ctx context.Context, exec boil.ContextExecutor, sc *models.ServerComponent) error {
	existing, err := models.ServerComponents(
		models.ServerComponentWhere.Name.EQ(sc.Name),
		models.ServerComponentWhere.ServerID.EQ(sc.ServerID),
		models.ServerComponentWhere.Serial.EQ(sc.Serial),
		models.ServerComponentWhere.ServerComponentTypeID.EQ(sc.ServerComponentTypeID),
	).One(ctx, exec)

	switch err {
	case nil:
		sc.ID = existing.ID
		_, updErr := sc.Update(ctx, exec, boil.Infer())
		return updErr
	case sql.ErrNoRows:
		return sc.Insert(ctx, exec, boil.Infer())
	default:
		return err
	}
}

// This encapsulates much of the repetitive work of getting a component to the database layer.
// The caller needs to compose the correct attributes for its given component.
func composeRecords(ctx context.Context, exec boil.ContextExecutor, cmn *common.Common,
	deviceID, namespace, slug string, attr *attributes) error {
	typeID := dbtools.MustComponentTypeID(ctx, exec, slug)

	sc := &models.ServerComponent{
		Name:                  null.StringFrom(slug),
		Vendor:                null.NewString(cmn.Vendor, cmn.Vendor != ""),
		Model:                 null.NewString(cmn.Model, cmn.Model != ""),
		Serial:                null.NewString(cmn.Serial, cmn.Serial != ""),
		ServerID:              deviceID,
		ServerComponentTypeID: typeID,
	}

	prodName := strings.TrimSpace(cmn.ProductName)
	if sc.Model.IsZero() && prodName != "" {
		sc.Model.SetValid(prodName)
	}

	if sc.Serial.IsZero() {
		sc.Serial = null.StringFrom("0")
	}

	if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
		return errors.Wrap(errComponent, slug+": "+err.Error())
	}

	// avoid computing this twice
	attr.ProductName = prodName

	attrData := attr.MustJSON()

	// update the component attribute
	if err := updateAnyAttribute(ctx, exec, false, sc.ID, namespace, attrData); err != nil {
		return errors.Wrap(errAttribute, slug+": "+err.Error())
	}

	// compose the versioned attributes
	vattr := &versionedAttributes{
		Firmware: cmn.Firmware,
		Status:   cmn.Status,
	}

	if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, namespace, vattr.MustJSON()); err != nil {
		return errors.Wrap(errVersionedAttr, slug+": "+err.Error())
	}

	return nil
}

func (dv *DeviceView) ComposeComponents(ctx context.Context, exec boil.ContextExecutor) error {
	if err := dv.writeBios(ctx, exec); err != nil {
		return err
	}
	if err := dv.writeBMC(ctx, exec); err != nil {
		return err
	}
	if err := dv.writeMainboard(ctx, exec); err != nil {
		return err
	}
	if err := dv.writeDimms(ctx, exec); err != nil {
		return err
	}
	return nil
}

func (dv *DeviceView) writeBios(ctx context.Context, exec boil.ContextExecutor) error {
	bios := dv.Inv.BIOS

	namespace := inbandComponentNamespace
	if !dv.Inband {
		namespace = outofbandComponentNamespace
	}

	attr := &attributes{
		Capabilities:  bios.Capabilities,
		CapacityBytes: bios.CapacityBytes,
		Description:   bios.Description,
		Metadata:      bios.Metadata,
		Oem:           bios.Oem,
		SizeBytes:     bios.SizeBytes,
	}

	return composeRecords(ctx, exec, &bios.Common, dv.DeviceID.String(), namespace, common.SlugBIOS, attr)
}

func (dv *DeviceView) writeBMC(ctx context.Context, exec boil.ContextExecutor) error {
	bmc := dv.Inv.BMC

	namespace := inbandComponentNamespace
	if !dv.Inband {
		namespace = outofbandComponentNamespace
	}

	attr := &attributes{
		Capabilities: bmc.Capabilities,
		Description:  bmc.Description,
		Metadata:     bmc.Metadata,
		Oem:          bmc.Oem,
	}

	return composeRecords(ctx, exec, &bmc.Common, dv.DeviceID.String(), namespace, common.SlugBMC, attr)
}

func (dv *DeviceView) writeMainboard(ctx context.Context, exec boil.ContextExecutor) error {
	mb := dv.Inv.Mainboard

	namespace := inbandComponentNamespace
	if !dv.Inband {
		namespace = outofbandComponentNamespace
	}

	attr := &attributes{
		Capabilities: mb.Capabilities,
		Description:  mb.Description,
		Metadata:     mb.Metadata,
		Oem:          mb.Oem,
		PhysicalID:   mb.PhysicalID,
	}

	return composeRecords(ctx, exec, &mb.Common, dv.DeviceID.String(), namespace, common.SlugMainboard, attr)
}

func (dv *DeviceView) writeDimms(ctx context.Context, exec boil.ContextExecutor) error {
	for idx, dimm := range dv.Inv.Memory {
		// skip bogus dimms
		if dimm.Vendor == "" &&
			dimm.ProductName == "" &&
			dimm.SizeBytes == 0 &&
			dimm.ClockSpeedHz == 0 {
			continue
		}

		if strings.TrimSpace(dimm.Serial) == "" {
			dimm.Serial = fmt.Sprintf("%d", idx)
		}

		namespace := inbandComponentNamespace
		if !dv.Inband {
			namespace = outofbandComponentNamespace
		}

		attr := &attributes{
			Capabilities: dimm.Capabilities,
			ClockSpeedHz: dimm.ClockSpeedHz,
			Description:  dimm.Description,
			FormFactor:   dimm.FormFactor,
			Metadata:     dimm.Metadata, // maybe this should be versioned?
			PartNumber:   dimm.PartNumber,
			SizeBytes:    dimm.SizeBytes,
			Slot:         strings.TrimPrefix(dimm.Slot, "DIMM.Socket."),
		}

		if err := composeRecords(ctx, exec, &dimm.Common, dv.DeviceID.String(), namespace, common.SlugPhysicalMem, attr); err != nil {
			return err
		}
	}
	return nil
}
