package tc_properties

// A Properties contains the key/value pairs from the properties input.
// All values are stored in unexpanded form and are expanded at runtime
type Properties struct {
	// Stores the key/value pairs
	m map[string]string

	// WriteSeparator specifies the separator of key and value while writing the properties.
	WriteSeparator string
}

// NewProperties creates a new Properties struct with the default
// configuration for "${key}" expressions.
func NewProperties() *Properties {
	return &Properties{
		m: map[string]string{},
	}
}

// Get returns the expanded value for the given key if exists.
// Otherwise, ok is false.
func (p *Properties) Get(key string) (string, bool) {
	v, ok := p.m[key]
	return v, ok
}
