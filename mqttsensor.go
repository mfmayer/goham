package goham

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type sensorValueMQTTConfig struct {
	Name              string `json:"name"`
	UniqueID          string `json:"unique_id"`
	StateTopic        string `json:"state_topic"`
	ValueTemplate     string `json:"value_template"`
	DeviceClass       string `json:"device_class,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	Icon              string `json:"icon,omitempty"`
}

// MQTTPublisherFunc type is an adapter to allow the use of ordinary functions as
// MQTT publisher. If f is a fuction with the appropriate signature, MQTTPublisherFunc(f)
// is a MQTTPublisher that calls f
type MQTTPublisherFunc func(topic string, qos byte, retained bool, payload interface{})

// Publish calls f(topic, qos, retained, payload)
func (f MQTTPublisherFunc) Publish(topic string, qos byte, retained bool, payload interface{}) {
	f(topic, qos, retained, payload)
}

type MQTTPublisher interface {
	Publish(topic string, qos byte, retained bool, payload interface{})
}

type MQTTSensor struct {
	*sensor
	mqttClient MQTTPublisher
	// hassConfig HassConfig
	// handleMeterReading func(message *gosml.ListEntry)
	// handlePowerReading func(message *gosml.ListEntry)
}

type mqttSensorOption struct {
	newSensorOptions []newSensorOption
}

type newMQTTSensorOption func(*mqttSensorOption)

func WithSensorID(id string) newMQTTSensorOption {
	return func(opt *mqttSensorOption) {
		opt.newSensorOptions = append(opt.newSensorOptions, withSensorID(id))
	}
}

func NewMQTTSensor(client MQTTPublisher, sensorName string, opts ...newMQTTSensorOption) *MQTTSensor {
	option := mqttSensorOption{}
	for _, opt := range opts {
		opt(&option)
	}
	return &MQTTSensor{
		sensor:     newSensor(sensorName, option.newSensorOptions...),
		mqttClient: client,
	}
}

type mqttValueOption struct {
	deviceClass        string
	unitOfMeasurement  string
	icon               string
	sensorValueOptions []newSensorValueOption
}

type addValueOption func(o *mqttValueOption)

func WithSensorValueID(id string) addValueOption {
	return func(o *mqttValueOption) {
		o.sensorValueOptions = append(o.sensorValueOptions, withSensorValueID(id))
	}
}

// WithDeviceClass to set a device class for the sensor value (see: https://www.home-assistant.io/integrations/sensor/#device-class)
func WithDeviceClass(deviceClass string) addValueOption {
	return func(o *mqttValueOption) {
		o.deviceClass = deviceClass
	}
}

// WithUnitAndIcon can be used alternatively to WithDeviceClass to add a custom unit and icon
func WithUnitOfMeasurement(unitOfMeasurement string) addValueOption {
	return func(o *mqttValueOption) {
		o.unitOfMeasurement = unitOfMeasurement
	}
}

// WithUnitAndIcon can be used alternatively to WithDeviceClass to add a custom unit and icon
func WithIcon(icon string) addValueOption {
	return func(o *mqttValueOption) {
		o.icon = icon
	}
}

// AddValue adds a new mqtt value to the mqtt sensor
func (ms *MQTTSensor) AddValue(valueName string, opts ...addValueOption) MQTTSensorValue {
	option := &mqttValueOption{}
	for _, opt := range opts {
		opt(option)
	}

	sensorValue := &mqttSensorValue{
		sensorValue: newSensorValue(ms, valueName, option.sensorValueOptions...),
	}

	topicBase := fmt.Sprintf("homeassistant/sensor/%s/", sensorValue.ID())
	sensorValue.configTopic = topicBase + "config"
	sensorValue.stateTopic = topicBase + "state"

	config := sensorValueMQTTConfig{
		Name:              sensorValue.Name(),
		UniqueID:          sensorValue.ID(),
		StateTopic:        sensorValue.stateTopic,
		ValueTemplate:     "{{ value }}",
		DeviceClass:       option.deviceClass,
		UnitOfMeasurement: option.unitOfMeasurement,
		Icon:              option.icon,
	}
	jsonConfig, _ := json.Marshal(config)
	sensorValue.configMessage = jsonConfig

	ms.values = append(ms.values, sensorValue)
	return sensorValue
}

type MQTTSensorValue interface {
	SensorValue
	// Update and publish a sensor value
	Update(float64)
	// PublishConfig sends autodiscover configuration that can be understood by home assistant
	PublishConfig()
}

type mqttSensorValue struct {
	*sensorValue
	lastlyPublishedConfig time.Time
	configMessage         []byte
	configTopic           string
	stateTopic            string
}

// Update and publish sensor value
// In case no autodiscover config has been sent yet or last config has been published more than 10 minutes
// before, publish a new autodiscover config message to homa assistant
func (sv *mqttSensorValue) Update(value float64) {
	if sv.lastlyPublishedConfig.IsZero() || time.Since(sv.lastlyPublishedConfig) > (10*time.Minute) {
		sv.PublishConfig()
		sv.lastlyPublishedConfig = time.Now()
	}
	if mqttSensor, ok := sv.sensor.(*MQTTSensor); ok {
		mqttSensor.mqttClient.Publish(sv.stateTopic, 0, false, strconv.FormatFloat(value, 'f', 4, 64))
	}
}

func (sv *mqttSensorValue) PublishConfig() {
	if mqttSensor, ok := sv.sensor.(*MQTTSensor); ok {
		mqttSensor.mqttClient.Publish(sv.configTopic, 0, false, sv.configMessage)
	}
}
