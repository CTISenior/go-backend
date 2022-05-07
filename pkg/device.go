package pkg

import (
	"fmt"

	"github.com/google/uuid"
)

type I interface{}

type Telemetry struct {
	Values    string //string(valueObj)
	Timestamp interface{}
	//values map[string]string
}

type Alert struct {
	TelemetryKey   string
	TelemetryValue float64
	SeverityType   string
	Severity       string
}

type Device struct {
	ID        *uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	AssetID   *uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	TenantID  *string
	Types     []string
	MaxValues []string
	Telemetry
	Alert
}

func (a Alert) PrepareAlertMessage() string {
	return fmt.Sprintf("%s is %.2f", a.TelemetryKey, a.TelemetryValue)
}
