package pkg

import (
	"github.com/google/uuid"
)

type Device struct {
	ID        *uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	AssetID   *uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	TenantID  *string
	Types     []string
	MaxValues []string
}

type Alert struct {
	Device
	TelemetryKey   string
	TelemetryValue float64
	SeverityType   string
	Severity       string
}

type Telemetry struct {
	Device
	values    string //string(valueObj)
	timestamp interface{}
	//values map[string]string
}

type I interface{}
