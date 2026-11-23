package main

import (
	"reflect"
	"testing"
)

func TestRefundTenants(t *testing.T) {
	for v, want := range map[string]map[string]bool{
		"":                           nil,
		"on":                         nil,
		"off":                        nil,
		"tenants:":                   {},
		"tenants:user-a":             {"user-a": true},
		" tenants: user-a , ,user-b": {"user-a": true, "user-b": true},
		"tenant:user-a":              nil,
	} {
		if got := refundTenants(v); !reflect.DeepEqual(got, want) {
			t.Errorf("%q: %v, want %v", v, got, want)
		}
	}
}

// The hand-over is on unless the environment says exactly "off".
func TestChromeHandIsOnUnlessTurnedOff(t *testing.T) {
	for v, want := range map[string]bool{
		"":      true,
		"on":    true,
		"true":  true,
		"off":   false,
		" off ": false,
	} {
		if got := chromeHandOn(v); got != want {
			t.Errorf("chromeHandOn(%q) = %v, want %v", v, got, want)
		}
	}
}
