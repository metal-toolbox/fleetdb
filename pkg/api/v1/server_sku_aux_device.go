package fleetdbapi

import (
	"github.com/metal-toolbox/fleetdb/internal/models"
)

// AuxDevice represents a Sku Aux Device Layout for a Server
type AuxDevice struct {
	Vendor     string `json:"vendor" binding:"required"`
	Model      string `json:"model" binding:"required"`
	DeviceType string `json:"device_type" binding:"required"`
	Details    []byte `json:"details" binding:"required"`
}

// toDBModelServerSkuAuxDevice converts a AuxDevice struct into a models.ServerSkuAuxDevice
func (aux *AuxDevice) toDBModelServerSkuAuxDevice() *models.ServerSkuAuxDevice {
	model := &models.ServerSkuAuxDevice{
		Vendor:     aux.Vendor,
		Model:      aux.Model,
		DeviceType: aux.DeviceType,
		Details:    aux.Details,
	}

	return model
}

// fomDBModelServerSkuAuxDevice converts a models.ServerSkuAuxDevice into a AuxDevice struct
func (aux *AuxDevice) fromDBModelServerSkuAuxDevice(model *models.ServerSkuAuxDevice) {
	aux.Vendor = model.Vendor
	aux.Model = model.Model
	aux.DeviceType = model.DeviceType
	aux.Details = model.Details
}
