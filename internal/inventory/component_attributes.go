//nolint:all  // XXX remove this!
package inventory

import (
	"encoding/json"

	"github.com/bmc-toolbox/common"
)

type attributes struct {
	Architecture                 string               `json:"architecture,omitempty"`
	BlockSizeBytes               int64                `json:"block_size_bytes,omitempty"`
	BusInfo                      string               `json:"bus_info,omitempty"`
	Capabilities                 []*common.Capability `json:"capabilities,omitempty"`
	CapableSpeedGbps             int64                `json:"capable_speed_gbps,omitempty"`
	CapacityBytes                int64                `json:"capacity_bytes,omitempty" diff:"immutable"`
	ChassisType                  string               `json:"chassis_type,omitempty"`
	ClockSpeedHz                 int64                `json:"clock_speed_hz,omitempty"`
	Cores                        int                  `json:"cores,omitempty"`
	Description                  string               `json:"description,omitempty"`
	DriveType                    string               `json:"drive_type,omitempty"`
	FormFactor                   string               `json:"form_factor,omitempty"`
	ID                           string               `json:"id,omitempty"`
	InterfaceType                string               `json:"interface_type,omitempty"`
	MacAddress                   string               `json:"macaddress,omitempty"`
	Metadata                     map[string]string    `json:"metadata,omitempty"`
	NegotiatedSpeedGbps          int64                `json:"negotiated_speed_gbps,omitempty"`
	Oem                          bool                 `json:"oem,omitempty"`
	OemID                        string               `json:"oem_id,omitempty"`
	PartNumber                   string               `json:"part_number,omitempty"`
	PhysicalID                   string               `json:"physid,omitempty"`
	PowerCapacityWatts           int64                `json:"power_capacity_watts,omitempty"`
	ProductName                  string               `json:"product_name,omitempty"`
	Protocol                     string               `json:"protocol,omitempty"`
	SizeBytes                    int64                `json:"size_bytes,omitempty"`
	Slot                         string               `json:"slot,omitempty"`
	SmartErrors                  []string             `json:"smart_errors,omitempty"`
	SmartStatus                  string               `json:"smart_status,omitempty"`
	SpeedBits                    int64                `json:"speed_bits,omitempty"`
	SpeedGbps                    int64                `json:"speed_gbps,omitempty"`
	StorageController            string               `json:"storage_controller,omitempty"`
	SupportedControllerProtocols string               `json:"supported_controller_protocol,omitempty"`
	SupportedDeviceProtocols     string               `json:"supported_device_protocol,omitempty"`
	SupportedRAIDTypes           string               `json:"supported_raid_types,omitempty"`
	Threads                      int                  `json:"threads,omitempty"`
	WWN                          string               `json:"wwn,omitempty"`
}

func (a *attributes) MustJSON() []byte {
	byt, err := json.Marshal(a)
	if err != nil {
		panic("bad attributes")
	}
	return byt
}

func (a *attributes) FromJSON(byt []byte) error {
	return json.Unmarshal(byt, a)
}

func mustFirmwareJSON(fw *common.Firmware) []byte {
	if fw == nil {
		panic("missing firmware payload")
	}
	byt, err := json.Marshal(fw)
	if err != nil {
		panic("bad firmware payload")
	}
	return byt
}

func firmwareFromJSON(byt []byte) (*common.Firmware, error) {
	fw := &common.Firmware{}
	err := json.Unmarshal(byt, fw)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func mustStatusJSON(st *common.Status) []byte {
	if st == nil {
		panic("missing status payload")
	}
	byt, err := json.Marshal(st)
	if err != nil {
		panic("bad status payload")
	}
	return byt
}

func statusFromJSON(byt []byte) (*common.Status, error) {
	st := &common.Status{}
	err := json.Unmarshal(byt, st)
	if err != nil {
		return nil, err
	}
	return st, nil
}
