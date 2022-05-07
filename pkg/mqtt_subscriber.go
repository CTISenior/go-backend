package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const ( // .env
	PORT  = 1883
	TOPIC = "v1/gateway/telemetry"
	QOS   = 1
)

var db = SetupDB()

//var DeviceDBStruct DeviceDB

var connectHandler mqtt.OnConnectHandler = func(c mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
	t := c.Subscribe(TOPIC, QOS, messagePubHandler)
	go func() {
		_ = t.Wait() // '<-t.Done()'
		if t.Error() != nil {
			fmt.Printf("Error occurred while subscribing: %s\n", t.Error())
		} else {
			fmt.Println("Subscribed to: ", TOPIC)
		}
	}()
}

var messagePubHandler mqtt.MessageHandler = func(c mqtt.Client, msg mqtt.Message) {

	logMsg := fmt.Sprintf("\nReceived MQTT message: %s from topic: %s", msg.Payload(), msg.Topic())
	fmt.Printf(logMsg + "\n")
	//InfoLogger.Println(logMsg)

	var res map[string]interface{ I }
	json.Unmarshal([]byte(msg.Payload()), &res)
	key := reflect.ValueOf(res).MapKeys()
	deviceSN := key[0].Interface().(string)

	deviceMap := res[deviceSN].([]interface{})[0].(map[string]interface{})

	jsonObj, err := json.Marshal(deviceMap)
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		//ErrorLogger.Println(err.Error())
	} else {
		// Next -> automatically create database record based on device metadata

		DeviceStruct := GetDeviceInfo(deviceSN)

		//if DeviceDBStruct.ID == nil {
		if DeviceStruct.IsStructureEmpty() {
			fmt.Printf("DeviceDB structure is empty\n")
			//ErrorLogger.Println(err.Error())
		} else {
			// Kafka Producer
			ProduceMessage(DeviceStruct.ID.String(), string(jsonObj))

			// DB Insert Operation
			valueObj, _ := json.Marshal(deviceMap["values"])

			TelemetryDBStruct := Telemetry{DeviceStruct, string(valueObj), deviceMap["ts"]}
			//async
			InsertTelemetryDB(TelemetryDBStruct)       //defer
			CheckDeviceValues(DeviceStruct, deviceMap) //defer
		}
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("MQTT Connection lost: %v", err)
	//WarningLogger.Println(logMsg)
}

var reconnectHandler mqtt.ReconnectHandler = func(mqtt.Client, *mqtt.ClientOptions) {
	fmt.Println("Attempting to reconnect to the MQTT broker")
	//WarningLogger.Println(logMsg)
}

func InitMQTTClient() {
	opts := mqtt.NewClientOptions()

	brokerPort, _ := strconv.ParseInt(GetEnvVariable("MQTT_PORT"), 10, 64)

	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", GetEnvVariable("MQTT_HOST"), brokerPort))
	opts.SetClientID(GetEnvVariable("MQTT_CLIENT"))
	opts.SetUsername(GetEnvVariable("MQTT_USER"))
	opts.SetPassword(GetEnvVariable("MQTT_PASSWORD"))
	// TLS/SSL
	opts.SetOrderMatters(true)

	opts.ConnectTimeout = time.Second
	opts.WriteTimeout = time.Second
	opts.KeepAlive = 10
	opts.PingTimeout = time.Second
	opts.ConnectRetry = true
	opts.AutoReconnect = true

	// Handlers
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.OnReconnecting = reconnectHandler

	client := mqtt.NewClient(opts)
	//client.AddRoute(TOPIC, h.handle)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("MQTT Connection is up")

}
