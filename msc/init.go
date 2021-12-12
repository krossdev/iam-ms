// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

// Initial setup message broker
func Initial(servers []string) {
	if len(servers) == 0 {
		logger.Panic("no servers give to connect")
	}
	createBroker(servers)

	// try to make connection to message server, the connection is permanent
	// once connected, and will auto reconnect when connection is broken.
	//
	// if connect failed, for example,nats not startup yet, the lib will
	// try to make connection later when publish data to message server,
	// if that failed again, then report error to app.
	if err := broker.connect(); err != nil {
		logger.Warnf("failed to connect: %v", err)
		return
	}
	logger.Tracef("message broker connected")
}

// Deinitial clean up message broker
func Deinitial() {
	if broker != nil {
		broker.disconnect()
	}
}

// MaxPayloadSize return maximum allowed payload size.
func MaxPayloadSize() int64 {
	return broker.maxPayloadSize()
}
