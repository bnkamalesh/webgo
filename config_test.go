package webgo

import "testing"

func TestConfig(t *testing.T) {
	cfg := Config{}
	cfg.Load("tests/config.json")

	cfg.Port = "a"
	if cfg.Validate() != ErrInvalidPort {
		t.Log("Port validation failed")
		t.Fail()
	}
	cfg.Port = "65536"
	if cfg.Validate() != ErrInvalidPort {
		t.Log("Port validation failed")
		t.Fail()
	}
}
