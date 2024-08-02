package fleetdbapi

import (
	"time"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuAuxDevice represents a SKU Aux Device for a Server
type ServerSkuAuxDevice struct {
	ID         string    `json:"id"`
	SkuID      string    `json:"sku_id" binding:"required"`
	Vendor     string    `json:"vendor" binding:"required"`
	Model      string    `json:"model" binding:"required"`
	DeviceType string    `json:"device_type" binding:"required"`
	Details    []byte    `json:"details" binding:"required"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// toDBModelServerSkuAuxDevice converts a ServerSkuAuxDevice into a models.ServerSkuAuxDevice (created by sqlboiler)
func (aux *ServerSkuAuxDevice) toDBModelServerSkuAuxDevice() *models.ServerSkuAuxDevice {
	model := &models.ServerSkuAuxDevice{
		ID:         aux.ID,
		SkuID:      aux.SkuID,
		Vendor:     aux.Vendor,
		Model:      aux.Model,
		DeviceType: aux.DeviceType,
		Details:    aux.Details,
	}

	return model
}

// fomDBModelServerSkuAuxDevice converts a models.ServerSkuAuxDevice (created by sqlboiler) into a ServerSkuAuxDevice
func (aux *ServerSkuAuxDevice) fromDBModelServerSkuAuxDevice(model *models.ServerSkuAuxDevice) {
	aux.ID = model.ID
	aux.SkuID = model.SkuID
	aux.Vendor = model.Vendor
	aux.Model = model.Model
	aux.DeviceType = model.DeviceType
	aux.Details = model.Details
	aux.CreatedAt = model.CreatedAt.Time
	aux.UpdatedAt = model.UpdatedAt.Time
}
