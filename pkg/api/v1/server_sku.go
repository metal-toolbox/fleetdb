package fleetdbapi

import (
	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSku represents a SKU for a Server
type ServerSku struct {
	Name             string      `json:"name" binding:"required"`
	Version          string      `json:"version" binding:"required"`
	Vendor           string      `json:"vendor" binding:"required"`
	Chassis          string      `json:"chassis" binding:"required"`
	BMCModel         string      `json:"bmc_model" binding:"required"`
	MotherboardModel string      `json:"motherboard_model" binding:"required"`
	CPUVendor        string      `json:"cpu_vendor" binding:"required"`
	CPUModel         string      `json:"cpu_model" binding:"required"`
	CPUCores         int64       `json:"cpu_cores" binding:"required"`
	CPUHertz         int64       `json:"cpu_hertz" binding:"required"`
	CPUCount         int64       `json:"cpu_count" binding:"required"`
	AuxDevices       []AuxDevice `json:"aux_devices" binding:"required"`
	Disks            []Disk      `json:"disks" binding:"required"`
	Memory           []Memory    `json:"memory" binding:"required"`
	Nics             []Nic       `json:"nics" binding:"required"`
}

// toDBModelServerSku converts a ServerSku into a models.ServerSku
func (sku *ServerSku) toDBModelServerSku() *models.ServerSku {
	model := &models.ServerSku{
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
func (sku *ServerSku) toDBModelServerSkuDeep(id string) *models.ServerSku {
	dbSku := sku.toDBModelServerSku()

	if len(sku.AuxDevices) > 0 || len(sku.Disks) > 0 || len(sku.Memory) > 0 || len(sku.Nics) > 0 {
		dbSku.R = dbSku.R.NewStruct()

		for i := range sku.AuxDevices {
			dbSku.R.SkuServerSkuAuxDevices = append(dbSku.R.SkuServerSkuAuxDevices, sku.AuxDevices[i].toDBModelServerSkuAuxDevice())
			dbSku.R.SkuServerSkuAuxDevices[i].SkuID = id
		}

		for i := range sku.Disks {
			dbSku.R.SkuServerSkuDisks = append(dbSku.R.SkuServerSkuDisks, sku.Disks[i].toDBModelServerSkuDisk())
			dbSku.R.SkuServerSkuDisks[i].SkuID = id
		}

		for i := range sku.Memory {
			dbSku.R.SkuServerSkuMemories = append(dbSku.R.SkuServerSkuMemories, sku.Memory[i].toDBModelServerSkuMemory())
			dbSku.R.SkuServerSkuMemories[i].SkuID = id
		}

		for i := range sku.Nics {
			dbSku.R.SkuServerSkuNics = append(dbSku.R.SkuServerSkuNics, sku.Nics[i].toDBModelServerSkuNic())
			dbSku.R.SkuServerSkuNics[i].SkuID = id
		}
	}

	return dbSku
}

// fromDBModelServerSku converts a models.ServerSku into a ServerSku
func (sku *ServerSku) fromDBModelServerSku(model *models.ServerSku) {
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

	if model.R != nil {
		diskCount := len(model.R.SkuServerSkuDisks)
		if diskCount > 0 {
			sku.Disks = make([]Disk, diskCount)
			for i, disk := range model.R.SkuServerSkuDisks {
				sku.Disks[i].fromDBModelServerSkuDisk(disk)
			}
		}

		memoryCount := len(model.R.SkuServerSkuMemories)
		if memoryCount > 0 {
			sku.Memory = make([]Memory, memoryCount)
			for i, memory := range model.R.SkuServerSkuMemories {
				sku.Memory[i].fromDBModelServerSkuMemory(memory)
			}
		}

		nicCount := len(model.R.SkuServerSkuNics)
		if nicCount > 0 {
			sku.Nics = make([]Nic, nicCount)
			for i, nic := range model.R.SkuServerSkuNics {
				sku.Nics[i].fromDBModelServerSkuNic(nic)
			}
		}

		auxDeviceCount := len(model.R.SkuServerSkuAuxDevices)
		if auxDeviceCount > 0 {
			sku.AuxDevices = make([]AuxDevice, auxDeviceCount)
			for i, auxDevice := range model.R.SkuServerSkuAuxDevices {
				sku.AuxDevices[i].fromDBModelServerSkuAuxDevice(auxDevice)
			}
		}
	}
}
