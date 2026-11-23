package planque

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func suffixLists(t *testing.T) suffixList {
	t.Helper()
	var l suffixList
	if err := json.Unmarshal(publicSuffixesJSON, &l); err != nil || len(l.Registry) == 0 || len(l.Hosting) == 0 {
		t.Fatal("the list", err)
	}
	return l
}

// A mark may not take every name under a domain anyone can have a name under:
// one test for each suffix on the list, as "*.<suffix>"; and, under shared
// hosting, anything beneath it too ("*.s3.amazonaws.com", "*.acme.github.io").
func TestAWildcardOnAPublicSuffixIsRefused(t *testing.T) {
	l := suffixLists(t)
	e, _, devices := catalogEnv(t)
	refused := func(t *testing.T, host string) {
		t.Helper()
		st, out := markAsk(t, e, devices["op_ada"], "op_ada", "Vendor "+strings.NewReplacer(".", " ", "*", "w", "-", " ").Replace(host), []string{host}, "2026-11-17")
		if st != http.StatusBadRequest {
			t.Fatalf("%s: %d %v", host, st, out)
		}
	}
	for _, suffix := range append(slices.Clone(l.Registry), l.Hosting...) {
		t.Run(suffix, func(t *testing.T) { refused(t, "*."+suffix) })
	}
	for _, suffix := range l.Hosting {
		t.Run("beneath "+suffix, func(t *testing.T) {
			refused(t, "*.acme."+suffix)
			refused(t, "*.s3."+suffix)
		})
	}
	// The case a review named.
	refused(t, "*.s3.amazonaws.com")
}

// What stays allowed: a vendor's own per-customer domain; a vendor's own domain
// one label under a registry suffix; exact hosts, wherever they are.
func TestAVendorsOwnWildcardIsAccepted(t *testing.T) {
	e, _, devices := catalogEnv(t)
	for i, domain := range []string{"zendesk.com", "freshdesk.com", "myshopify.com", "gorgias.com", "teamwork.com", "acme.co.uk", "vendor.com.au"} {
		name := []string{"Zendesk", "Freshdesk", "Shopify", "Gorgias", "Teamwork", "Acme", "Vendor"}[i]
		if st, out := markAsk(t, e, devices["op_ada"], "op_ada", name, []string{"*." + domain}, "2026-11-17"); st != http.StatusOK {
			t.Errorf("*.%s: %d %v", domain, st, out)
		}
	}
	if st, out := markAsk(t, e, devices["op_ada"], "op_ada", "Moneybird", []string{"moneybird.com", "api.moneybird.com", "acme.github.io"}, "2026-11-17"); st != http.StatusOK {
		t.Errorf("exact hosts: %d %v", st, out)
	}
}

// The list itself: lowercase, sorted, distinct, in the form a mark takes,
// registry and hosting not overlapping, and no vendor's own domain the review
// checks give a wildcard to (the two would contradict).
func TestThePublicSuffixListIsSoundAndHoldsNoVendorsDomain(t *testing.T) {
	l := suffixLists(t)
	for name, list := range map[string][]string{"registry": l.Registry, "hosting": l.Hosting} {
		if !slices.IsSorted(list) {
			t.Errorf("%s is sorted", name)
		}
		for i, s := range list {
			if s != strings.ToLower(s) || strings.Count(s, ".") < 1 || !markHost.MatchString("*."+s) || slices.Contains(list[:i], s) {
				t.Errorf("%s %q: lowercase, distinct, two labels or more", name, s)
			}
		}
	}
	for _, s := range l.Registry {
		if slices.Contains(l.Hosting, s) {
			t.Errorf("%q is on both lists", s)
		}
	}
	for vendor, hosts := range vendorHosts {
		for _, h := range hosts {
			if domain, wild := strings.CutPrefix(h, "*."); wild && isPublicSuffix(domain) {
				t.Errorf("%s: its wildcard %s is on the public suffix list", vendor, h)
			}
		}
	}
}
