package mqtt

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/SoNim-LSCM/TKOH_OMS/errors"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

var server *mqtt.Server

func MqttSetup() {
	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	server = mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":1883", nil)
	err := server.AddListener(tcp)
	errors.CheckFatalError(err)

	err = server.AddHook(new(ExampleHook), map[string]any{})
	errors.CheckFatalError(err)

	go func() {
		err := server.Serve()
		errors.CheckFatalError(err)
	}()

	// Demonstration of using an inline client to directly subscribe to a topic and receive a message when
	// that subscription is activated. The inline subscription method uses the same internal subscription logic
	// as used for external (normal) clients.
	go func() {
		// Inline subscriptions can also receive retained messages on subscription.
		_ = server.Publish("direct/publish", []byte("test message"), false, 0)

		server.Log.Info("inline client subscribing")
		_ = server.Subscribe("direct/#", 1, exampleCallback)
	}()

	// Run server until interrupted
	<-done

	// Cleanup
}

func PublishMqtt(topic string, msg []byte) {
	errors.CheckError(server.Publish(topic, msg, false, 0), "mqtt publishing")
}

type ExampleHook struct {
	mqtt.HookBase
}

func exampleCallback(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	server.Log.Info("exampleCallback() received message from subscription", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
}
