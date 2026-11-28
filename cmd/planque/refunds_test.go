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

// PLANQUE_BROWSER_SESSIONS: unset is 1; 1 to 3 are taken; anything else
// stops planque.
func TestBrowserSessionsAreOneToThree(t *testing.T) {
	for v, want := range map[string]int{"": 1, "1": 1, " 2 ": 2, "3": 3} {
		if n, err := browserSessionsFrom(v); err != nil || n != want {
			t.Errorf("%q: %d, %v", v, n, err)
		}
	}
	for _, v := range []string{"0", "4", "two", "-1", "1.5"} {
		if _, err := browserSessionsFrom(v); err == nil {
			t.Errorf("%q taken", v)
		}
	}
}
