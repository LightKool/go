package http

import (
	"bytes"
	"net/url"
	"sort"
)

// Param hold url paramters for http requests.
// It's a map-like structure but holds multiple values for
// a single key while no duplicate key-value pair is allowd.
// Also the values of every single key are sorted.
// All keys and values are strings.
type Param struct {
	data map[string][]string
	keys []string
}

func NewParam() *Param {
	return &Param{make(map[string][]string), make([]string, 0)}
}

// Return the values of the key or nil if not exists.
func (p *Param) Get(key string) []string {
	return p.data[key]
}

// Get the count of the values of the key.
func (p *Param) Count(key string) int {
	exist, ok := p.data[key]
	if ok {
		return len(exist)
	} else {
		return 0
	}
}

// Add all values of the key.
func (p *Param) Add(key string, values ...string) bool {
	if key == "" || len(values) == 0 {
		return false
	}
	exist, ok := p.data[key]
	if ok {
		buf := make([]string, 0, len(values))
		for _, v := range values {
			i := sort.SearchStrings(exist, v)
			if v != exist[i] {
				buf = append(buf, v)
			}
		}
		if len(buf) == 0 {
			return false
		}
		exist = append(exist, buf...)
		sort.Strings(exist)
		p.data[key] = exist
	} else {
		sort.Strings(values)
		p.data[key] = values
		p.keys = append(p.keys, key)
		sort.Strings(p.keys)
	}
	return true
}

// Replace values of the key. If the key didn't exist,
// it will be added just as Param.Add(key, values).
func (p *Param) Replace(key string, values ...string) bool {
	if key == "" || len(values) == 0 {
		return false
	}
	_, ok := p.data[key]
	sort.Strings(values)
	p.data[key] = values
	if !ok {
		p.keys = append(p.keys, key)
		sort.Strings(p.keys)
		return false
	}
	return true
}

// Remove all values of the key.
func (p *Param) RemoveAll(key string) bool {
	_, ok := p.data[key]
	if ok {
		delete(p.data, key)
		p.keys = remove(p.keys, key)
		return true
	}
	return false
}

// Remove all the values of the key. If all the values
// have been removed after the operation, the key will
// also be removed just as the Param.RemoveAll(key).
func (p *Param) Remove(key string, values ...string) {
	exist, ok := p.data[key]
	if ok {
		for _, v := range values {
			exist = remove(exist, v)
		}
		if len(exist) == 0 {
			delete(p.data, key)
			p.keys = remove(p.keys, key)
		} else {
			p.data[key] = exist
		}
	}
}

// Generate the content use to calculate the signature
// according to the SMP API specification.
func (p *Param) ContentToSign() []byte {
	content := make([]byte, 0, 64)
	for _, key := range p.keys {
		content = append(content, key...)
		for _, v := range p.data[key] {
			content = append(content, v...)
		}
	}
	return content
}

// Generate http query string. Copied from url.Value.Encode().
func (p *Param) QueryString() string {
	var buf bytes.Buffer
	for _, k := range p.keys {
		vs := p.data[k]
		prefix := url.QueryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// Private function to remove a string from a sorted slice of strings.
func remove(sorted []string, value string) []string {
	i := sort.SearchStrings(sorted, value)
	if i < len(sorted) && value == sorted[i] {
		return append(sorted[:i], sorted[(i+1):]...)
	} else {
		return sorted
	}
}
