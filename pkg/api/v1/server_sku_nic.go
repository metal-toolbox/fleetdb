package fleetdbapi

import (
	"github.com/metal-toolbox/fleetdb/internal/models"
)

// Nic represents a Sku Nic Layout for a Server
type Nic struct {
	Vendor        string `json:"vendor" binding:"required"`
	Model         string `json:"model" binding:"required"`
	PortBandwidth int64  `json:"port_bandwidth" binding:"required"`
	PortCount     int64  `json:"port_count" binding:"required"`
	Count         int64  `json:"count" binding:"required"`
}

// toDBModelServerSkuNic converts a Nic struct into a models.ServerSkuNic
func (nic *Nic) toDBModelServerSkuNic() *models.ServerSkuNic {
	model := &models.ServerSkuNic{
		Vendor:        nic.Vendor,
		Model:         nic.Model,
		PortBandwidth: nic.PortBandwidth,
		PortCount:     nic.PortCount,
		Count:         nic.Count,
	}

	return model
}

// fromDBModelServerSkuNic converts a models.ServerSkuNic into a Nic struct
func (nic *Nic) fromDBModelServerSkuNic(model *models.ServerSkuNic) {
	nic.Vendor = model.Vendor
	nic.Model = model.Model
	nic.PortBandwidth = model.PortBandwidth
	nic.PortCount = model.PortCount
	nic.Count = model.Count
}
