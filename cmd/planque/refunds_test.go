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
