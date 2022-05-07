package pkg

import (
	"fmt"
	"reflect"
)

func (d Device) IsStructureEmpty() bool {
	return reflect.DeepEqual(d, Device{})
}

func (a Alert) PrepareAlertMessage() string {
	return fmt.Sprintf("%s is %.2f", a.TelemetryKey, a.TelemetryValue)
}
