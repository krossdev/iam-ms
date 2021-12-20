package msc

import "github.com/sirupsen/logrus"

// the default Logger has no any configuration(like log to file...),
// app should replace it with SetLogger()
var Logger = logrus.StandardLogger().WithField("realm", "ms")

// replace the default logger with customized logger
func SetLogger(e *logrus.Entry) {
	Logger = e
}

// Connect to message broker
func Connect(servers []string) {
	if len(servers) == 0 {
		Logger.Panic("no servers give to connect")
	}
	createBroker(servers)

	// try to make connection to message server, the connection is permanent
	// once connected, and will auto reconnect when connection is broken.
	//
	// if connect failed, for example,nats not startup yet, the lib will
	// try to make connection later when publish data to message server,
	// if that failed again, then report error to app.
	if err := broker.connect(); err != nil {
		Logger.Warnf("failed to connect: %v", err)
		return
	}
	Logger.Tracef("message broker connected")
}

// Disconnect connection from message broker
func Disconnect() {
	if broker != nil {
		broker.disconnect()
	}
}

// MaxPayloadSize return maximum allowed payload size.
func MaxPayloadSize() int64 {
	return broker.maxPayloadSize()
}
