//nolint:all  // XXX remove this!
package inventory

import (
	"encoding/json"

	"github.com/bmc-toolbox/common"
	rivets "github.com/metal-toolbox/rivets/types"
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

func mustFirmwareJSON(fw *common.Firmware) []byte {
	if fw == nil {
		return nil
	}
	byt, err := json.Marshal(fw)
	if err != nil {
		panic("bad firmware payload")
	}
	return byt
}

func firmwareFromJSON(byt []byte) (*common.Firmware, error) {
	if byt == nil {
		return nil, nil
	}
	fw := &common.Firmware{}
	err := json.Unmarshal(byt, fw)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func mustStatusJSON(st *common.Status) []byte {
	if st == nil {
		return nil
	}
	byt, err := json.Marshal(st)
	if err != nil {
		panic("bad status payload")
	}
	return byt
}

func statusFromJSON(byt []byte) (*common.Status, error) {
	if byt == nil {
		return nil, nil
	}
	st := &common.Status{}
	err := json.Unmarshal(byt, st)
	if err != nil {
		return nil, err
	}
	return st, nil
}
