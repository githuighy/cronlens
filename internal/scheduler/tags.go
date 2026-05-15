package scheduler

import (
	"fmt"
	"sort"
	"strings"
)

// Tags holds a set of key-value metadata labels attached to a Schedule.
type Tags map[string]string

// Set adds or overwrites a tag.
func (t Tags) Set(key, value string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("tag key must not be empty")
	}
	t[key] = value
	return nil
}

// Get returns the value for a tag key and whether it was found.
func (t Tags) Get(key string) (string, bool) {
	v, ok := t[key]
	return v, ok
}

// Delete removes a tag by key. It is a no-op if the key does not exist.
func (t Tags) Delete(key string) {
	delete(t, key)
}

// Keys returns a sorted slice of all tag keys.
func (t Tags) Keys() []string {
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// String returns a human-readable representation of the tags.
func (t Tags) String() string {
	if len(t) == 0 {
		return "(no tags)"
	}
	parts := make([]string, 0, len(t))
	for _, k := range t.Keys() {
		parts = append(parts, fmt.Sprintf("%s=%s", k, t[k]))
	}
	return strings.Join(parts, ", ")
}

// Clone returns a shallow copy of the Tags map.
func (t Tags) Clone() Tags {
	copy := make(Tags, len(t))
	for k, v := range t {
		copy[k] = v
	}
	return copy
}
