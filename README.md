# Goham

Goham is a Go module that provides an easy-to-use interface for creating MQTT sensors and publishing sensor data to a MQTT broker. The module is compatible with the Home Assistant autodiscovery feature, allowing the automatic creation of sensors in Home Assistant based on the published data.

## Installation

To install the Goham module, run the following command:

```bash
go get -u github.com/mfmayer/goham
```

## Usage

Here's a basic example of how to use the Goham module to create a MQTT sensor and publish sensor data:

```go
package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mfmayer/goham"
)

func main() {
	// Connect to the MQTT broker
	opts := mqtt.NewClientOptions().AddBroker("tcp://yourmqttbroker:1883")
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	// Create a publisher function that matches the MQTTPublisher interface
	mqttPublisher := goham.MQTTPublisherFunc(func(topic string, qos byte, retained bool, payload interface{}) {
		token := client.Publish(topic, qos, retained, payload)
    token.Wait()
	})

	// Create a new MQTT sensor
	sensor := goham.NewMQTTSensor(mqttPublisher, "kitchen")

	// Add a sensor value with a device class
	value := sensor.AddValue("temperature", goham.WithDeviceClass("temperature"))

	// Update and publish the sensor value
	value.Update(25.0)
}
```

Replace yourmqttbroker with the address of your MQTT broker and yourusername with your GitHub username or the username of the repository owner.

For more examples and available options, please refer to the source code and comments.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
