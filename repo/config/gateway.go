package config

// Gateway contains options for the HTTP gateway server.
type Gateway struct {
	RootRedirect string
	Writable     bool

	// The url of a newline delimited list of keys that the gateway should not serve
	BlackList string

	// The url of a newline delimited list of keys that the gateway should only serve
	WhiteList string
}
