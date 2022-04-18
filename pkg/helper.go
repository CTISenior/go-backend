package pkg

import (
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

type Alert struct {
	Type    string
	Message string
}

func (a Alert) PrepareMessage() string {
	return fmt.Sprintf("[%s] - %s", a.Type, a.Message)
}

var DeviceDBStruct DeviceDB

func getDeviceInfo(deviceSN string) DeviceDB {
	var devDBObj DeviceDB
	devSql := "SELECT id, asset_id, tenant_id, types, max_values FROM devices WHERE sn=$1"

	err := db.QueryRow(devSql, deviceSN).Scan(&devDBObj.ID, &devDBObj.AssetID, &devDBObj.TenantID, pq.Array(&devDBObj.Types), pq.Array(&devDBObj.MaxValues))
	if err != nil {
		fmt.Print(err.Error() + "\n")
	}

	return devDBObj
}

func InsertTelemetryDB(deviceSN string, valueObj []byte, timestamp interface{}) {
	DeviceDBStruct = getDeviceInfo(deviceSN)

	sqlStatement := `
INSERT INTO device_telemetries (sn, value, device_id, asset_id, tenant_id, timestamp)
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sqlStatement,
		deviceSN, string(valueObj), DeviceDBStruct.ID, DeviceDBStruct.AssetID, DeviceDBStruct.TenantID, timestamp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The telemetry inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}
}

func InsertAlertDB(deviceSN string, msgType string, message string, telemetry_key string, timestamp interface{}) {
	sqlStatement := `
INSERT INTO device_alerts (sn, type, telemetry_key, message, device_id, asset_id, tenant_id, timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(sqlStatement,
		deviceSN, msgType, telemetry_key, message, DeviceDBStruct.ID, DeviceDBStruct.AssetID, DeviceDBStruct.TenantID, timestamp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The alert inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}

}

func CheckDeviceValues(deviceSN string, deviceMap map[string]interface{}) {
	values := deviceMap["values"].(map[string]interface{})
	for k, v := range values {
		//if value, ok := v.(string); ok {
		value, _ := strconv.ParseInt(v.(string), 10, 64)
		/*switch k {

		}*/

		// DeviceDBStruct.Types
		// DeviceDBStruct.MaxValues

		msg := fmt.Sprintf("%s is %d", k, value)
		if k == "temperature" && value > 50 {
			fmt.Println(msg)
			//log
			InsertAlertDB(deviceSN, "warning", msg, k, deviceMap["ts"])
		}
		if k == "humidity" && value > 70 {
			fmt.Println(msg)
			//log
			InsertAlertDB(deviceSN, "warning", msg, k, deviceMap["ts"])
		}

		//}
	}
}

func InitHelper() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
