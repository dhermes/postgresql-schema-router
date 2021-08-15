package server

const (
	// DefaultPort is the default value for `Config.Port`.
	DefaultPort = 5397
)

// Config represents the values needed to configure a server.
type Config struct {
	// Port is the port where the proxy should expose the server
	Port int
}
