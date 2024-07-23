package fleetdbapi

import (
	"time"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuNic represents a SKU for a Server
type ServerSkuNic struct {
	ID            string    `json:"id"`
	SkuID         string    `json:"sku_id" binding:"required"`
	PortBandwidth int64     `json:"port_bandwidth" binding:"required"`
	PortCount     int64     `json:"port_count" binding:"required"`
	Count         int64     `json:"count" binding:"required"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

// toDBModelServerSkuNic converts a ServerSkuNic into a models.ServerSkuNic
func (nic *ServerSkuNic) toDBModelServerSkuNic() *models.ServerSkuNic {
	model := &models.ServerSkuNic{
		ID:            nic.ID,
		SkuID:         nic.SkuID,
		PortBandwidth: nic.PortBandwidth,
		PortCount:     nic.PortCount,
		Count:         nic.Count,
	}

	return model
}

// fromDBModelServerSkuNic converts a models.ServerSkuNic into a ServerSkuNic
func (nic *ServerSkuNic) fromDBModelServerSkuNic(model *models.ServerSkuNic) {
	nic.ID = model.ID
	nic.SkuID = model.SkuID
	nic.PortBandwidth = model.PortBandwidth
	nic.PortCount = model.PortCount
	nic.Count = model.Count
	nic.CreatedAt = model.CreatedAt.Time
	nic.UpdatedAt = model.UpdatedAt.Time
}
