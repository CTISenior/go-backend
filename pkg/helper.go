package pkg

import (
	"reflect"
	"strconv"
	"strings"
)

func (d Device) IsStructureEmpty() bool {
	return reflect.DeepEqual(d, Device{})
}

func CheckDeviceValues(deviceMap map[string]interface{}) {
	values := deviceMap["values"].(map[string]interface{})

	for key, val := range values {
		telemetryValue, _ := strconv.ParseFloat(val.(string), 64)

		for i := 0; i < len(DeviceStruct.SensorTypes); i++ {
			if strings.ToLower(key) == strings.ToLower(DeviceStruct.SensorTypes[i]) {
				maxValue, _ := strconv.ParseFloat(DeviceStruct.MaxValues[i], 64)
				//minValue, _ := strconv.ParseFloat(DeviceStruct.MinValues[i], 64)

				DeviceStruct.Alert.TelemetryKey = key
				DeviceStruct.Alert.TelemetryValue = telemetryValue

				if telemetryValue >= (maxValue + maxValue/4.0) {
					DeviceStruct.Alert.SeverityType = "max"
					DeviceStruct.Alert.Severity = "critical"
					msg := DeviceStruct.Alert.PrepareAlertMessage()
					insertAlertDB(msg)
					//log
				} else if telemetryValue >= maxValue {
					DeviceStruct.Alert.SeverityType = "max"
					DeviceStruct.Alert.Severity = "warning"
					msg := DeviceStruct.Alert.PrepareAlertMessage()
					insertAlertDB(msg)
					//log
				}

				/*switch {
					case telemetryValue >= (maxValue + maxValue/4.0):
					case telemetryValue >= maxValue:
					case telemetryValue <= (minValue - minValue/2.0):
					case telemetryValue <= (maxValue + maxValue/4.0):
				}*/

			}
		}
	}
}
