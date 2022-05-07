package pkg

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	devSql := "SELECT id, asset_id, tenant_id, types, max_values FROM devices WHERE sn=$1"
	err := db.QueryRow(devSql, deviceSN).Scan(
		&devDBObj.ID,
		&devDBObj.AssetID,
		&devDBObj.TenantID,
		pq.Array(&devDBObj.Types),
		pq.Array(&devDBObj.MaxValues))
	if err != nil {
		fmt.Print(err.Error() + "\n")
	}

	return devDBObj
}

func InsertTelemetryDB(telemetry Telemetry) {
	sqlStatement := `
INSERT INTO device_telemetries (values, timestamp, device_id, asset_id, tenant_id)
VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement,
		telemetry.values,
		telemetry.timestamp,
		telemetry.Device.ID,
		telemetry.Device.AssetID,
		telemetry.Device.TenantID)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The telemetry inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}
}

func insertAlertDB(alert Alert, msg string) {
	sqlStatement := `
INSERT INTO device_alerts (telemetry_key, telemetry_value, severity_type, severity, message, device_id, asset_id, tenant_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(sqlStatement,
		alert.TelemetryKey,
		alert.TelemetryValue,
		alert.SeverityType,
		alert.Severity,
		msg,
		alert.Device.ID,
		alert.Device.AssetID,
		alert.Device.TenantID)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The alert inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}
}

func CheckDeviceValues(device Device, deviceMap map[string]interface{}) {
	values := deviceMap["values"].(map[string]interface{})

	for key, val := range values {
		telemetryValue, _ := strconv.ParseFloat(val.(string), 64)

		for i := 0; i < len(device.Types); i++ {
			if key == device.Types[i] {
				maxValue, _ := strconv.ParseFloat(device.MaxValues[i], 32)

				if telemetryValue >= (maxValue + maxValue/4.0) {
					AlertStruct := Alert{device, key, telemetryValue, "max", "critical"}
					msg := AlertStruct.PrepareAlertMessage()
					insertAlertDB(AlertStruct, msg)
					//log
				} else if telemetryValue >= maxValue {
					AlertStruct := Alert{device, key, telemetryValue, "max", "warning"}
					msg := AlertStruct.PrepareAlertMessage()
					insertAlertDB(AlertStruct, msg)
					//log
				}

			}
		}
	}
}
