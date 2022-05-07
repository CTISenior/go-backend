package pkg

import (
	"reflect"
)

func (d Device) IsStructureEmpty() bool {
	return reflect.DeepEqual(d, Device{})
}
