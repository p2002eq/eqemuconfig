package eqemuconfig

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	c, err := GetConfig()
	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}
	t.Log(c.Discord.Channels)
}
