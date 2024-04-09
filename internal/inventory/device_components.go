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
	typeID := dbtools.MustComponentTypeID(ctx, exec, common.SlugBIOS)

	bios := dv.Inv.BIOS
	sc := &models.ServerComponent{
		Name:                  null.StringFrom(common.SlugBIOS),
		Vendor:                null.NewString(bios.Vendor, bios.Vendor != ""),
		Model:                 null.NewString(bios.Model, bios.Model != ""),
		Serial:                null.NewString(bios.Serial, bios.Serial != ""),
		ServerID:              dv.DeviceID.String(),
		ServerComponentTypeID: typeID,
	}

	prodName := strings.TrimSpace(bios.ProductName)
	if sc.Model.IsZero() && prodName != "" {
		sc.Model.SetValid(prodName)
	}
	if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
		return errors.Wrap(errComponent, "bios: "+err.Error())
	}

	namespace := inbandComponentNamespace
	if !dv.Inband {
		namespace = outofbandComponentNamespace
	}

	attrData := (&attributes{
		Capabilities:  bios.Capabilities,
		CapacityBytes: bios.CapacityBytes,
		Description:   bios.Description,
		Metadata:      bios.Metadata,
		Oem:           bios.Oem,
		ProductName:   prodName,
		SizeBytes:     bios.SizeBytes,
	}).MustJSON()

	// update the component attribute
	if err := updateAnyAttribute(ctx, exec, false, sc.ID, namespace, attrData); err != nil {
		return errors.Wrap(errAttribute, "bios: "+err.Error())
	}

	// compose the versioned attributes
	biosVA := &versionedAttributes{
		Firmware: bios.Firmware,
		Status:   bios.Status,
	}

	if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, namespace, biosVA.MustJSON()); err != nil {
		return errors.Wrap(errVersionedAttr, "bios: "+err.Error())
	}

	return nil
}

func (dv *DeviceView) writeBMC(ctx context.Context, exec boil.ContextExecutor) error {
	typeID := dbtools.MustComponentTypeID(ctx, exec, common.SlugBMC)

	bmc := dv.Inv.BMC
	sc := &models.ServerComponent{
		Name:                  null.StringFrom(common.SlugBMC),
		Vendor:                null.NewString(bmc.Vendor, bmc.Vendor != ""),
		Model:                 null.NewString(bmc.Model, bmc.Model != ""),
		Serial:                null.NewString(bmc.Serial, bmc.Serial != ""),
		ServerID:              dv.DeviceID.String(),
		ServerComponentTypeID: typeID,
	}

	prodName := strings.TrimSpace(bmc.ProductName)
	if sc.Model.IsZero() && prodName != "" {
		sc.Model.SetValid(prodName)
	}
	if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
		return errors.Wrap(errComponent, "bmc: "+err.Error())
	}

	namespace := inbandComponentNamespace
	if !dv.Inband {
		namespace = outofbandComponentNamespace
	}

	attrData := (&attributes{
		Capabilities: bmc.Capabilities,
		Description:  bmc.Description,
		Metadata:     bmc.Metadata,
		Oem:          bmc.Oem,
		ProductName:  prodName,
	}).MustJSON()

	if err := updateAnyAttribute(ctx, exec, false, sc.ID, namespace, attrData); err != nil {
		return errors.Wrap(errAttribute, "bmc: "+err.Error())
	}

	// compose the versioned attributes
	bmcVA := &versionedAttributes{
		Firmware: bmc.Firmware,
		Status:   bmc.Status,
	}

	if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, namespace, bmcVA.MustJSON()); err != nil {
		return errors.Wrap(errVersionedAttr, "bmc: "+err.Error())
	}

	return nil
}

func (dv *DeviceView) writeMainboard(ctx context.Context, exec boil.ContextExecutor) error {
	typeID := dbtools.MustComponentTypeID(ctx, exec, common.SlugMainboard)

	mb := dv.Inv.Mainboard
	sc := &models.ServerComponent{
		Name:                  null.StringFrom(common.SlugMainboard),
		Vendor:                null.NewString(mb.Vendor, mb.Vendor != ""),
		Model:                 null.NewString(mb.Model, mb.Model != ""),
		Serial:                null.NewString(mb.Serial, mb.Serial != ""),
		ServerID:              dv.DeviceID.String(),
		ServerComponentTypeID: typeID,
	}

	prodName := strings.TrimSpace(mb.ProductName)
	if sc.Model.IsZero() && prodName != "" {
		sc.Model.SetValid(prodName)
	}

	if sc.Serial.IsZero() {
		sc.Serial = null.StringFrom("0")
	}

	if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
		return errors.Wrap(errComponent, "mb: "+err.Error())
	}

	namespace := inbandComponentNamespace
	if !dv.Inband {
		namespace = outofbandComponentNamespace
	}

	attrData := (&attributes{
		Capabilities: mb.Capabilities,
		Description:  mb.Description,
		Metadata:     mb.Metadata,
		Oem:          mb.Oem,
		PhysicalID:   mb.PhysicalID,
		ProductName:  prodName,
	}).MustJSON()

	if err := updateAnyAttribute(ctx, exec, false, sc.ID, namespace, attrData); err != nil {
		return errors.Wrap(errAttribute, "mb: "+err.Error())
	}

	// compose the versioned attributes
	mbVA := &versionedAttributes{
		Firmware: mb.Firmware,
		Status:   mb.Status,
	}

	if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, namespace, mbVA.MustJSON()); err != nil {
		return errors.Wrap(errVersionedAttr, "mb: "+err.Error())
	}

	return nil
}

func (dv *DeviceView) writeDimms(ctx context.Context, exec boil.ContextExecutor) error {
	typeID := dbtools.MustComponentTypeID(ctx, exec, common.SlugPhysicalMem)

	for idx, dimm := range dv.Inv.Memory {
		// skip bogus dimms
		if dimm.Vendor == "" &&
			dimm.ProductName == "" &&
			dimm.SizeBytes == 0 &&
			dimm.ClockSpeedHz == 0 {
			continue
		}

		sc := &models.ServerComponent{
			Name:                  null.StringFrom(common.SlugPhysicalMem),
			Vendor:                null.NewString(dimm.Vendor, dimm.Vendor != ""),
			Model:                 null.NewString(dimm.Model, dimm.Model != ""),
			Serial:                null.NewString(dimm.Serial, dimm.Serial != ""),
			ServerID:              dv.DeviceID.String(),
			ServerComponentTypeID: typeID,
		}

		// set incrementing serial when one isn't found
		if sc.Serial.IsZero() {
			sc.Serial.SetValid(fmt.Sprintf("%d", idx))
		}

		prodName := strings.TrimSpace(dimm.ProductName)
		if sc.Model.IsZero() && prodName != "" {
			sc.Model.SetValid(prodName)
		}

		if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
			return errors.Wrap(errComponent, "dimm: "+err.Error())
		}

		namespace := inbandComponentNamespace
		if !dv.Inband {
			namespace = outofbandComponentNamespace
		}

		attrData := (&attributes{
			Capabilities: dimm.Capabilities,
			ClockSpeedHz: dimm.ClockSpeedHz,
			Description:  dimm.Description,
			FormFactor:   dimm.FormFactor,
			Metadata:     dimm.Metadata, // maybe this should be versioned?
			PartNumber:   dimm.PartNumber,
			ProductName:  prodName,
			SizeBytes:    dimm.SizeBytes,
			Slot:         strings.TrimPrefix(dimm.Slot, "DIMM.Socket."),
		}).MustJSON()

		// update the component attribute
		if err := updateAnyAttribute(ctx, exec, false, sc.ID, namespace, attrData); err != nil {
			return errors.Wrap(errAttribute, "dimm: "+err.Error())
		}

		// compose the versioned attributes for this dimm
		dimmVA := &versionedAttributes{
			Firmware: dimm.Firmware,
			Status:   dimm.Status,
		}

		if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, namespace, dimmVA.MustJSON()); err != nil {
			return errors.Wrap(errVersionedAttr, "dimm: "+err.Error())
		}
	}
	return nil
}
