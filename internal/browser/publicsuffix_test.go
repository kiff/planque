package browser

import (
	"strings"
	"testing"
)

// A shared host or a registry's suffix isn't a site a grant may name: under
// it, anyone has a site. One owner's site under it is.
func TestASharedDomainIsntASite(t *testing.T) {
	for _, shared := range []string{
		// The private section's shared hosts.
		"github.io", "vercel.app", "netlify.app", "myshopify.com", "notion.site", "blogspot.com",
		"herokuapp.com", "pages.dev", "web.app", "s3.amazonaws.com",
		// A wildcard rule's names, one label under it.
		"anything.compute.amazonaws.com",
		// ICANN's suffixes.
		"com", "co.uk", "com.au",
	} {
		if !PublicSuffix(shared) {
			t.Errorf("%s: not read as a public suffix", shared)
		}
		if err := CheckSite(shared); err == nil {
			if strings.Contains(shared, ".") {
				t.Errorf("%s: taken as a site", shared)
			}
		} else if strings.Contains(shared, ".") && !strings.Contains(err.Error(), "shared domain") {
			t.Errorf("%s: %v", shared, err)
		}
	}
	for _, site := range []string{
		"shop.example", "shop.github.io", "my-store.myshopify.com", "docs.example.org", "example.co.uk",
		"amazon.com", "GitHub.com",
	} {
		if PublicSuffix(site) {
			t.Errorf("%s: read as a public suffix", site)
		}
		if err := CheckSite(strings.ToLower(site)); err != nil {
			t.Errorf("%s: %v", site, err)
		}
	}
	if !strings.Contains(publicSuffixRaw, "===BEGIN PRIVATE DOMAINS===") || len(publicSuffixes) < 9000 {
		t.Fatalf("the list holds %d rules, without its private section?", len(publicSuffixes))
	}
}
