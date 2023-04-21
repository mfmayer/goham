package goham

import (
	"fmt"

	"github.com/mfmayer/goham/internal/stringutils"
)

// Sensor interface
type Sensor interface {
	// Name returns the sensor name
	Name() string
	// ID returns the sensor id (consists only of the following characters: [a-zA-Z0-9_-])
	ID() string
}

// sensor represents a basic sensor
type sensor struct {
	name   string
	id     string
	values []SensorValue
}

type sensorOption struct {
	id string
}

type newSensorOption func(opt *sensorOption)

func withSensorID(id string) newSensorOption {
	return func(opt *sensorOption) {
		opt.id = id
	}
}

// newSensor creates a new sensor with given opts (e.g. WithSensorID)
func newSensor(sensorName string, opts ...newSensorOption) *sensor {
	option := sensorOption{
		id: sensorName,
	}
	for _, opt := range opts {
		opt(&option)
	}
	return &sensor{
		name:   sensorName,
		id:     stringutils.Sanitize(option.id),
		values: []SensorValue{},
	}
}

// Name returns sensor's name
func (s *sensor) Name() string {
	return s.name
}

// ID returns sensor's id
func (s *sensor) ID() string {
	return s.id
}

type SensorValue interface {
	// Name should return a combination of sensor name and value name in a readable form
	Name() string
	// ID should return a combination of sensor name and value name using [a-zA-Z0-9_-] (alphanumerics, underscore and hyphen)
	ID() string
}

// sensorValue represents a basic sensor value
type sensorValue struct {
	sensor Sensor
	name   string
	id     string
}

type sensorValueOption struct {
	id string
}

type newSensorValueOption func(opt *sensorValueOption)

func withSensorValueID(id string) newSensorValueOption {
	return func(opt *sensorValueOption) {
		opt.id = id
	}
}

// newSensorValue creates a new sensor value
func newSensorValue(sensor Sensor, valueName string, opts ...newSensorValueOption) *sensorValue {
	option := sensorValueOption{
		id: valueName,
	}
	for _, opt := range opts {
		opt(&option)
	}
	return &sensorValue{
		sensor: sensor,
		name:   valueName,
		id:     stringutils.Sanitize(option.id),
	}
}

// Name returns "<sensor.name> <value.name>"
func (sv *sensorValue) Name() string {
	return fmt.Sprintf("%s_%s", sv.sensor.Name(), sv.name)
}

// ID returns "<sensor.name>_<value.name>"
func (sv *sensorValue) ID() string {
	return fmt.Sprintf("%s_%s", sv.sensor.ID(), sv.id)
}
