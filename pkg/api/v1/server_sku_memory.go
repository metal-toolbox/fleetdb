package fleetdbapi

import (
	"github.com/metal-toolbox/fleetdb/internal/models"
)

// Memory represents a Sku Memory Layout for a Server
type Memory struct {
	Vendor string `json:"vendor" binding:"required"`
	Model  string `json:"model" binding:"required"`
	Bytes  int64  `json:"bytes" binding:"required"`
	Count  int64  `json:"count" binding:"required"`
}

// toDBModelServerSkuMemory converts a Memory struct into a models.ServerSkuMemory
func (mem *Memory) toDBModelServerSkuMemory() *models.ServerSkuMemory {
	model := &models.ServerSkuMemory{
		Vendor: mem.Vendor,
		Model:  mem.Model,
		Bytes:  mem.Bytes,
		Count:  mem.Count,
	}

	return model
}

// fromDBModelServerSkuMemory converts a models.ServerSkuMemory into a Memory struct
func (mem *Memory) fromDBModelServerSkuMemory(model *models.ServerSkuMemory) {
	mem.Vendor = model.Vendor
	mem.Model = model.Model
	mem.Bytes = model.Bytes
	mem.Count = model.Count
}
