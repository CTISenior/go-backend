package pkg

import (
	"github.com/google/uuid"
)

type Device struct {
	sn     string
	ts     int
	values map[string]string
}

type DeviceDB struct {
	ID        *uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	AssetID   *uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	TenantID  *string
	Types     []string
	MaxValues []string
}

type Telemetry struct {
	ts     int
	values map[string]string
}

type I interface{}
