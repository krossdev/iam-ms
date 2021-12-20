package msc

import "github.com/nats-io/nats.go"

// asyncErrorHandler sets the async error handler (e.g. slow consumer errors)
func asyncErrorHandler(c *nats.Conn, s *nats.Subscription, e error) {
	Logger.WithError(e).Error("Message broker async error")
}

// disconnectedErrorHandler sets the disconnected error handler that is called
// whenever the connection is disconnected.
// Disconnected error could be nil, for instance when user explicitly closes the
// connection.
func disconnectedErrorHandler(c *nats.Conn, e error) {
	Logger.WithError(e).Warn("Message broker disconnected")
}

// reconnectedHandler sets the reconnected handler called whenever
// the connection is successfully reconnected.
func reconnectedHandler(c *nats.Conn) {
	Logger.Info("Message broker reconnected")
}

// discoveredServersHandler sets the callback that is invoked whenever a new
// server has joined the cluster.
func discoveredServersHandler(c *nats.Conn) {
	Logger.Info("Message broker discover new server")
}

// closedHandler sets the closed handler that is called when a client will
// no longer be connected.
func closedHandler(c *nats.Conn) {
	Logger.Info("Message broker connection closed")
}
