package fleetdbapi

import (
	"time"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuDisk represents a SKU Disk Layout for a Server
type ServerSkuDisk struct {
	ID        string    `json:"id"`
	SkuID     string    `json:"sku_id" binding:"required"`
	Bytes     int64     `json:"bytes" binding:"required"`
	Protocol  string    `json:"protocol" binding:"required"`
	Count     int64     `json:"count" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// toDBModelServerSkuDisk converts a ServerSkuDisk into a models.ServerSkuDisk (created by sqlboiler)
func (disk *ServerSkuDisk) toDBModelServerSkuDisk() *models.ServerSkuDisk {
	model := &models.ServerSkuDisk{
		ID:       disk.ID,
		SkuID:    disk.SkuID,
		Bytes:    disk.Bytes,
		Protocol: disk.Protocol,
		Count:    disk.Count,
	}

	return model
}

// fromDBModelServerSkuDisk converts a models.ServerSkuDisk (created by sqlboiler) into a ServerSkuDisk
func (disk *ServerSkuDisk) fromDBModelServerSkuDisk(model *models.ServerSkuDisk) {
	disk.ID = model.ID
	disk.SkuID = model.SkuID
	disk.Bytes = model.Bytes
	disk.Protocol = model.Protocol
	disk.Count = model.Count
	disk.CreatedAt = model.CreatedAt.Time
	disk.UpdatedAt = model.UpdatedAt.Time
}
