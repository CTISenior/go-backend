package pkg

type Device struct {
	sn     string
	ts     int
	values map[string]string
}

type Telemetry struct {
	ts     int
	values map[string]string
}

type I interface{}
