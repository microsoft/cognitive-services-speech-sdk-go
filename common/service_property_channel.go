package common

// ServicePropertyChannel defines channels used to pass property settings to service.
// Added in version 1.5.0.
type ServicePropertyChannel int

const (
	// URIQueryParameter uses URI query parameter to pass property settings to service.
	URIQueryParameter   ServicePropertyChannel = 0
)
