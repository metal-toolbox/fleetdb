package fleetdbapi

import (
	"github.com/metal-toolbox/fleetdb/internal/models"
)

// Disk represents a Sku Disk Layout for a Server
type Disk struct {
	Vendor   string `json:"vendor" binding:"required"`
	Model    string `json:"model" binding:"required"`
	Bytes    int64  `json:"bytes" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
	Count    int64  `json:"count" binding:"required"`
}

// toDBModelServerSkuDisk converts a Disk struct into a models.ServerSkuDisk
func (disk *Disk) toDBModelServerSkuDisk() *models.ServerSkuDisk {
	model := &models.ServerSkuDisk{
		Vendor:   disk.Vendor,
		Model:    disk.Model,
		Bytes:    disk.Bytes,
		Protocol: disk.Protocol,
		Count:    disk.Count,
	}

	return model
}

// fromDBModelServerSkuDisk converts a models.ServerSkuDisk into a Disk struct
func (disk *Disk) fromDBModelServerSkuDisk(model *models.ServerSkuDisk) {
	disk.Vendor = model.Vendor
	disk.Model = model.Model
	disk.Bytes = model.Bytes
	disk.Protocol = model.Protocol
	disk.Count = model.Count
}
