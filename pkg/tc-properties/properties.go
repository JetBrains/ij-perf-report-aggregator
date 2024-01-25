package tc_properties

import (
  "fmt"
)

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

// GetString returns the expanded value for the given key if exists or
// the default value otherwise.
func (p *Properties) GetString(key, def string) string {
  if v, ok := p.Get(key); ok {
    return v
  }
  return def
}

// Keys returns all keys in the same order as in the input.
func (p *Properties) Keys() []string {
  keys := make([]string, 0, len(p.m))
  for k := range p.m {
    keys = append(keys, k)
  }
  return keys
}

// String returns a string of all expanded 'key = value' pairs.
func (p *Properties) String() string {
  var s string
  for key := range p.m {
    value, _ := p.Get(key)
    s = fmt.Sprintf("%s%s = %s\n", s, key, value)
  }
  return s
}
