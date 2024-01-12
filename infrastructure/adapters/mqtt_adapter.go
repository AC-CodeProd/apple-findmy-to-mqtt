package adapters

import (
	"apple-findmy-to-mqtt/core/interfaces"
	"apple-findmy-to-mqtt/infrastructure/config"
	"bytes"
	"crypto/tls"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/fx"
)

type MqttClientParams struct {
	fx.In
	Config config.Config
}
type pahoMQTTClient struct {
	client MQTT.Client
}

func NewPahoMQTTClient(mcp MqttClientParams) interfaces.IMQTTClient {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	var raw_broker bytes.Buffer
	raw_broker.WriteString("tls://")
	raw_broker.WriteString(fmt.Sprintf("%s:%d", mcp.Config.Mqtt.Broker, mcp.Config.Mqtt.Port))
	opts := MQTT.NewClientOptions().AddBroker(raw_broker.String())
	opts.SetClientID(mcp.Config.Mqtt.ClientID)
	opts.SetUsername(mcp.Config.Mqtt.Username)
	opts.SetPassword(mcp.Config.Mqtt.Password)
	opts.SetTLSConfig(tlsConfig)

	return &pahoMQTTClient{
		client: MQTT.NewClient(opts),
	}
}

func (pqc *pahoMQTTClient) Connect() error {
	if token := pqc.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (pqc *pahoMQTTClient) Disconnect() {
	pqc.client.Disconnect(250)
}

func (pqc *pahoMQTTClient) Publish(topic string, payload []byte) error {
	token := pqc.client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (pqc *pahoMQTTClient) Subscribe(topic string, handler interfaces.MessageHandler) error {
	if token := pqc.client.Subscribe(topic, 0, func(client MQTT.Client, message MQTT.Message) {
		handler(message.Topic(), message.Payload())
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
