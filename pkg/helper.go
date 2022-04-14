package pkg

import (
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type Alert struct {
	Type    string
	Message string
}

func (a Alert) PrepareMessage() string {
	return fmt.Sprintf("[%s] - %s", a.Type, a.Message)
}

func InsertTelemetryDB(deviceSN string, valueObj []byte, timestamp interface{}) {
	sqlStatement := `
INSERT INTO device_telemetries (sn, device_id, value, timestamp)
VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement,
		deviceSN, 1, string(valueObj), timestamp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The telemetry is inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}
}

func InsertAlertDB(deviceSN string, msgType string, message string, telemetry_key string, timestamp interface{}) {
	sqlStatement := `
INSERT INTO device_alerts (sn, device_id, type, telemetry_key, message, timestamp)
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sqlStatement,
		deviceSN, 1, msgType, message, telemetry_key, timestamp)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		logMsg := "The telemetry is inserted into the database successfully"
		fmt.Printf(logMsg + "\n")
		//InfoLogger.Println(logMsg)
	}

	/*
	   id INT DEFAULT unique_rowid(),
	   sn STRING NOT NULL,
	   device_id INT8 NOT NULL,

	   type STRING(30) NOT NULL,
	   telemetry_key STRING(50) NOT NULL,
	   message TEXT NOT NULL,

	   timestamp INT8 NOT NULL,
	*/
}

func CheckDeviceValues(deviceSN string, deviceMap map[string]interface{}) {
	values := deviceMap["values"].(map[string]interface{})
	for k, v := range values {
		//if value, ok := v.(string); ok {
		value, _ := strconv.ParseInt(v.(string), 10, 64)
		msg := ""
		alert := false
		/*switch k {

		}*/

		// db - select statement
		if k == "temperature" && value > 50 {
			msg = "Warning - temperature"
			alert = true

		}
		if k == "humidity" && value > 70 {
			msg = "Warning - temperature"
			alert = true
		}

		if alert {
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
