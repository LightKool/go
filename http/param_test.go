package http

import (
	"testing"
)

func TestParam(t *testing.T) {
	p := NewParam()
	p.Add("key", "1", "2", "7", "9")
	t.Log(p.Count("key"))
	p.Add("api_key", "apikeycontent")
	t.Log(p.data)
	t.Log(p.keys)

	t.Log(p.Replace("key2", "hello world"))
	t.Log(p.Replace("key", "hello world"))
	t.Log(p.data)
	t.Log(p.keys)
}
