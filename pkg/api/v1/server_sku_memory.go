package fleetdbapi

import (
	"time"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuMemory represents a Sku Memory Layout for a Server
type ServerSkuMemory struct {
	ID        string    `json:"id"`
	SkuID     string    `json:"sku_id" binding:"required"`
	Bytes     int64     `json:"bytes" binding:"required"`
	Count     int64     `json:"count" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// toDBModelServerSkuMemory converts a ServerSkuMemory into a models.ServerSkuMemory
func (mem *ServerSkuMemory) toDBModelServerSkuMemory() *models.ServerSkuMemory {
	model := &models.ServerSkuMemory{
		ID:    mem.ID,
		SkuID: mem.SkuID,
		Bytes: mem.Bytes,
		Count: mem.Count,
	}

	return model
}

// fromDBModelServerSkuMemory converts a models.ServerSkuMemory into a ServerSkuMemory
func (mem *ServerSkuMemory) fromDBModelServerSkuMemory(model *models.ServerSkuMemory) {
	mem.ID = model.ID
	mem.SkuID = model.SkuID
	mem.Bytes = model.Bytes
	mem.Count = model.Count
	mem.CreatedAt = model.CreatedAt.Time
	mem.UpdatedAt = model.UpdatedAt.Time
}
