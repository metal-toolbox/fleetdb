package fleetdbapi

import (
	"time"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSku represents a SKU for a Server
type ServerSku struct {
	ID               string               `json:"id"`
	Name             string               `json:"name" binding:"required"`
	Version          string               `json:"version" binding:"required"`
	Vendor           string               `json:"vendor" binding:"required"`
	Chassis          string               `json:"chassis" binding:"required"`
	BMCModel         string               `json:"bmc_model" binding:"required"`
	MotherboardModel string               `json:"motherboard_model" binding:"required"`
	CPUVendor        string               `json:"cpu_vendor" binding:"required"`
	CPUModel         string               `json:"cpu_model" binding:"required"`
	CPUCores         int64                `json:"cpu_cores" binding:"required"`
	CPUHertz         int64                `json:"cpu_hertz" binding:"required"`
	CPUCount         int64                `json:"cpu_count" binding:"required"`
	AuxDevices       []ServerSkuAuxDevice `json:"aux_devices" binding:"required"`
	Disks            []ServerSkuDisk      `json:"disks" binding:"required"`
	Memory           []ServerSkuMemory    `json:"memory" binding:"required"`
	Nics             []ServerSkuNic       `json:"nics" binding:"required"`
	CreatedAt        time.Time            `json:"created_at,omitempty"`
	UpdatedAt        time.Time            `json:"updated_at,omitempty"`
}

// toDBModelServerSku converts a ServerSku into a models.ServerSku
func (sku *ServerSku) toDBModelServerSku() *models.ServerSku {
	model := &models.ServerSku{
		ID:               sku.ID,
		Name:             sku.Name,
		Version:          sku.Version,
		Vendor:           sku.Vendor,
		Chassis:          sku.Chassis,
		BMCModel:         sku.BMCModel,
		MotherboardModel: sku.MotherboardModel,
		CPUVendor:        sku.CPUVendor,
		CPUModel:         sku.CPUModel,
		CPUCores:         sku.CPUCores,
		CPUHertz:         sku.CPUHertz,
		CPUCount:         sku.CPUCount,
	}

	return model
}

// toDBModelServerSkuDeep converts a ServerSku into a models.ServerSku. It also includes all relations, doing a deep copy
func (sku *ServerSku) toDBModelServerSkuDeep() *models.ServerSku {
	dbSku := sku.toDBModelServerSku()

	if len(sku.AuxDevices) > 0 || len(sku.Disks) > 0 || len(sku.Memory) > 0 || len(sku.Nics) > 0 {
		dbSku.R = dbSku.R.NewStruct()

		for i := range sku.AuxDevices {
			dbSku.R.SkuServerSkuAuxDevices = append(dbSku.R.SkuServerSkuAuxDevices, sku.AuxDevices[i].toDBModelServerSkuAuxDevice())
		}

		for i := range sku.Disks {
			dbSku.R.SkuServerSkuDisks = append(dbSku.R.SkuServerSkuDisks, sku.Disks[i].toDBModelServerSkuDisk())
		}

		for i := range sku.Memory {
			dbSku.R.SkuServerSkuMemories = append(dbSku.R.SkuServerSkuMemories, sku.Memory[i].toDBModelServerSkuMemory())
		}

		for i := range sku.Nics {
			dbSku.R.SkuServerSkuNics = append(dbSku.R.SkuServerSkuNics, sku.Nics[i].toDBModelServerSkuNic())
		}
	}

	return dbSku
}

// fromDBModelServerSku converts a models.ServerSku into a ServerSku
func (sku *ServerSku) fromDBModelServerSku(model *models.ServerSku) {
	sku.ID = model.ID
	sku.Name = model.Name
	sku.Version = model.Version
	sku.Vendor = model.Vendor
	sku.Chassis = model.Chassis
	sku.BMCModel = model.BMCModel
	sku.MotherboardModel = model.MotherboardModel
	sku.CPUVendor = model.CPUVendor
	sku.CPUModel = model.CPUModel
	sku.CPUCores = model.CPUCores
	sku.CPUHertz = model.CPUHertz
	sku.CPUCount = model.CPUCount
	sku.CreatedAt = model.CreatedAt.Time
	sku.UpdatedAt = model.UpdatedAt.Time

	if model.R != nil {
		diskCount := len(model.R.SkuServerSkuDisks)
		if diskCount > 0 {
			sku.Disks = make([]ServerSkuDisk, diskCount)
			for i, disk := range model.R.SkuServerSkuDisks {
				sku.Disks[i].fromDBModelServerSkuDisk(disk)
			}
		}

		memoryCount := len(model.R.SkuServerSkuMemories)
		if memoryCount > 0 {
			sku.Memory = make([]ServerSkuMemory, memoryCount)
			for i, memory := range model.R.SkuServerSkuMemories {
				sku.Memory[i].fromDBModelServerSkuMemory(memory)
			}
		}

		nicCount := len(model.R.SkuServerSkuNics)
		if nicCount > 0 {
			sku.Nics = make([]ServerSkuNic, nicCount)
			for i, nic := range model.R.SkuServerSkuNics {
				sku.Nics[i].fromDBModelServerSkuNic(nic)
			}
		}

		auxDeviceCount := len(model.R.SkuServerSkuAuxDevices)
		if auxDeviceCount > 0 {
			sku.AuxDevices = make([]ServerSkuAuxDevice, auxDeviceCount)
			for i, auxDevice := range model.R.SkuServerSkuAuxDevices {
				sku.AuxDevices[i].fromDBModelServerSkuAuxDevice(auxDevice)
			}
		}
	}
}
