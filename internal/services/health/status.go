package health

type Status int

const (
	Healthy = iota
	Unhealthy
	Unknown
)

// MarshalText satisfies the encoding.TextMarshaler interface.
func (s Status) MarshalText() (text []byte, err error) {
	switch s {
	case Healthy:
		return []byte("Healthy"), nil
	case Unhealthy:
		return []byte("Unhealthy"), nil
	default:
		return []byte("Unknown"), nil
	}
}
