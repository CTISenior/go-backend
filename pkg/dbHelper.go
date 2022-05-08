package pkg

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetDeviceInfo(deviceSN string) Device {
	//var devDBObj DeviceDB
	devDBObj := Device{}

	devSql := "SELECT id, asset_id, tenant_id, sensor_types, max_values FROM devices WHERE sn=$1"
	err := db.QueryRow(devSql, deviceSN).Scan(
		&devDBObj.ID,
		&devDBObj.AssetID,
		&devDBObj.TenantID,
		pq.Array(&devDBObj.SensorTypes),
		pq.Array(&devDBObj.MaxValues))
	if err != nil {
		fmt.Print(err.Error() + "\n")
	}

	return devDBObj
}

func InsertTelemetryDB() {
	sqlStatement := `
INSERT INTO device_telemetries (values, timestamp, device_id, asset_id, tenant_id)
VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement,
		DeviceStruct.Telemetry.Values,
		DeviceStruct.Telemetry.Timestamp,
		DeviceStruct.ID,
		DeviceStruct.AssetID,
		DeviceStruct.TenantID)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The telemetry inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}
}

func insertAlertDB(msg string) {
	sqlStatement := `
INSERT INTO device_alerts (telemetry_key, telemetry_value, severity_type, severity, message, device_id, asset_id, tenant_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(sqlStatement,
		DeviceStruct.Alert.TelemetryKey,
		DeviceStruct.Alert.TelemetryValue,
		DeviceStruct.Alert.SeverityType,
		DeviceStruct.Alert.Severity,
		msg,
		DeviceStruct.ID,
		DeviceStruct.AssetID,
		DeviceStruct.TenantID)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The alert inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}
}

func CheckDeviceValues(deviceMap map[string]interface{}) {
	values := deviceMap["values"].(map[string]interface{})

	for key, val := range values {
		telemetryValue, _ := strconv.ParseFloat(val.(string), 64)

		for i := 0; i < len(DeviceStruct.SensorTypes); i++ {
			if strings.ToLower(key) == strings.ToLower(DeviceStruct.SensorTypes[i]) {
				maxValue, _ := strconv.ParseFloat(DeviceStruct.MaxValues[i], 64)

				if telemetryValue >= (maxValue + maxValue/4.0) {
					DeviceStruct.Alert.TelemetryKey = key
					DeviceStruct.Alert.TelemetryValue = telemetryValue
					DeviceStruct.Alert.SeverityType = "max"
					DeviceStruct.Alert.Severity = "critical"

					msg := DeviceStruct.Alert.PrepareAlertMessage()
					insertAlertDB(msg)
					//log
				} else if telemetryValue >= maxValue {
					DeviceStruct.Alert.TelemetryKey = key
					DeviceStruct.Alert.TelemetryValue = telemetryValue
					DeviceStruct.Alert.SeverityType = "max"
					DeviceStruct.Alert.Severity = "warningl"

					msg := DeviceStruct.Alert.PrepareAlertMessage()
					insertAlertDB(msg)
					//log
				}

			}
		}
	}
}
