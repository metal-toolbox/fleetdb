package inventory

import (
	"encoding/json"

	common "github.com/metal-toolbox/bmc-common"
	rivets "github.com/metal-toolbox/rivets/v2/types"
)

func mustAttributesJSON(ca *rivets.ComponentAttributes) []byte {
	if ca == nil {
		return nil
	}
	byt, err := json.Marshal(ca)
	if err != nil {
		panic("bad component attributes payload")
	}
	return byt
}

func componentAttributesFromJSON(byt []byte) (*rivets.ComponentAttributes, error) {
	if byt == nil {
		return nil, nil
	}
	ca := &rivets.ComponentAttributes{}
	if err := json.Unmarshal(byt, ca); err != nil {
		return nil, err
	}
	return ca, nil
}

// for historical reasons, alloy used to store firmware data like this using a different API
// we create this here to maintain compatibility with old data
type firmwareContainer struct {
	Firmware *common.Firmware `json:"firmware,omitempty"`
}

func mustFirmwareJSON(fw *common.Firmware) []byte {
	if fw == nil {
		return nil
	}

	fwc := &firmwareContainer{
		Firmware: fw,
	}
	byt, err := json.Marshal(fwc)
	if err != nil {
		panic("bad firmware payload")
	}
	return byt
}

func firmwareFromJSON(byt []byte) (*common.Firmware, error) {
	if byt == nil {
		return nil, nil
	}
	fwc := &firmwareContainer{}
	err := json.Unmarshal(byt, fwc)
	if err != nil {
		return nil, err
	}
	return fwc.Firmware, nil
}

type statusContainer struct {
	Status *common.Status `json:"status,omitempty"`
}

func mustStatusJSON(st *common.Status) []byte {
	if st == nil {
		return nil
	}

	// XXX: In order to remain backward-compatible with data collected with
	// earlier versions of alloy, status is encoded as an array.
	ary := []*statusContainer{
		{
			Status: st,
		},
	}
	byt, err := json.Marshal(ary)
	if err != nil {
		panic("bad status payload")
	}
	return byt
}

func statusFromJSON(byt []byte) (*common.Status, error) {
	if byt == nil {
		return nil, nil
	}
	dataAry := []*statusContainer{}
	err := json.Unmarshal(byt, &dataAry)
	if err != nil {
		// skip malformed garbage
		return nil, err
	}

	var st *common.Status
	for _, o := range dataAry {
		if o.Status != nil {
			// take the first valid status
			st = o.Status
			break
		}
	}
	return st, nil
}
