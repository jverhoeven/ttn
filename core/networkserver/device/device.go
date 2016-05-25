package device

import (
	"fmt"
	"time"

	"github.com/TheThingsNetwork/ttn/core/storage"
	"github.com/TheThingsNetwork/ttn/core/types"
)

// Options for the specified device
type Options struct {
	DisableFCntCheck bool // Disable Frame counter check (insecure)
	Uses32BitFCnt    bool // Use 32-bit Frame counters
}

// Device contains the state of a device
type Device struct {
	DevEUI      types.DevEUI
	AppEUI      types.AppEUI
	DevAddr     types.DevAddr
	NwkSKey     types.NwkSKey
	FCntUp      uint32
	FCntDown    uint32
	LastSeen    time.Time
	Options     Options
	Utilization Utilization
}

// DeviceProperties contains all properties of a Device that can be stored in Redis.
var DeviceProperties = []string{
	"dev_eui",
	"app_eui",
	"dev_addr",
	"nwk_s_key",
	"f_cnt_up",
	"f_cnt_down",
	"last_seen",
	"options",
	"utilization",
}

// ToStringStringMap converts the given properties of Device to a
// map[string]string for storage in Redis.
func (device *Device) ToStringStringMap(properties ...string) (map[string]string, error) {
	output := make(map[string]string)
	for _, p := range properties {
		property, err := device.formatProperty(p)
		if err != nil {
			return output, err
		}
		if property != "" {
			output[p] = property
		}
	}
	return output, nil
}

// FromStringStringMap imports known values from the input to a Device.
func (device *Device) FromStringStringMap(input map[string]string) error {
	for k, v := range input {
		device.parseProperty(k, v)
	}
	return nil
}

func (device *Device) formatProperty(property string) (formatted string, err error) {
	switch property {
	case "dev_eui":
		formatted = device.DevEUI.String()
	case "app_eui":
		formatted = device.AppEUI.String()
	case "dev_addr":
		formatted = device.DevAddr.String()
	case "nwk_s_key":
		formatted = device.NwkSKey.String()
	case "f_cnt_up":
		formatted = storage.FormatUint32(device.FCntUp)
	case "f_cnt_down":
		formatted = storage.FormatUint32(device.FCntDown)
	case "last_seen":
		formatted = device.LastSeen.Format(time.RFC3339Nano)
	case "options":
		// TODO
	case "utilization":
		// TODO
	default:
		err = fmt.Errorf("Property %s does not exist in Status", property)
	}
	return
}

func (device *Device) parseProperty(property string, value string) error {
	if value == "" {
		return nil
	}
	switch property {
	case "dev_eui":
		val, err := types.ParseDevEUI(value)
		if err != nil {
			return err
		}
		device.DevEUI = val
	case "app_eui":
		val, err := types.ParseAppEUI(value)
		if err != nil {
			return err
		}
		device.AppEUI = val
	case "dev_addr":
		val, err := types.ParseDevAddr(value)
		if err != nil {
			return err
		}
		device.DevAddr = val
	case "nwk_s_key":
		val, err := types.ParseNwkSKey(value)
		if err != nil {
			return err
		}
		device.NwkSKey = val
	case "f_cnt_up":
		val, err := storage.ParseUint32(value)
		if err != nil {
			return err
		}
		device.FCntUp = val
	case "f_cnt_down":
		val, err := storage.ParseUint32(value)
		if err != nil {
			return err
		}
		device.FCntDown = val
	case "last_seen":
		val, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return err
		}
		device.LastSeen = val
	case "options":
		// TODO
	case "utilization":
		// TODO
	}
	return nil
}
