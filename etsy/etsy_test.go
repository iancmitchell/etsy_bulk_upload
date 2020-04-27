package etsy

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	_, ok := interface{}(client).(Client)
	if !ok {
		t.Errorf("Client not set correctly.")
	}
}

func TestAuthentizate(t *testing.T) {
	client := NewClient()
	client.Authenticate()
}
