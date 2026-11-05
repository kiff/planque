package browser

import "testing"

// Every form a browser reads as an address is refused, written out or not;
// a site's name in another script is refused, its xn-- form isn't; and a
// listed site still opens, as do its subdomains.
func TestAllowedRefusesEveryFormOfAnAddress(t *testing.T) {
	for u, want := range map[string]string{
		"https://127.1/":                 CodeAddressPrivate,
		"https://0x7f.1/":                CodeAddressPrivate,
		"https://0177.0.0.1/":            CodeAddressPrivate,
		"https://1.2.3.0x4/":             CodeAddressPrivate,
		"https://2130706433/":            CodeAddressPrivate,
		"https://127.0.0。1/":             CodeAddressPrivate,
		"https://169.254.169．254/":       CodeAddressPrivate,
		"https://10.0.0｡1/":              CodeAddressPrivate,
		"https://localhost。/":            CodeAddressPrivate,
		"https://[::ffff:127.0.0.1]/":    CodeAddressPrivate,
		"https://bücher.example/":        CodeBadURL,
		"https://ｓhop.example/":          CodeBadURL,
		"https://xn--bcher-kva.example/": "",
		"https://shop.example/":          "",
		"https://shop.example./":         "",
		"https://a.shop.example/item":    "",
		"https://shop.example:8443/":     "",
		"https://shop1.example/":         "",
		"https://123shop.example/":       "",
	} {
		if got, why := Allowed(u, nil); got != want {
			t.Errorf("%s: %q (%s), want %q", u, got, why, want)
		}
	}
	if got, _ := Allowed("https://shop。example/", []string{"shop.example"}); got != "" {
		t.Errorf("a listed site written with an ideographic full stop: %q", got)
	}
	if got, _ := Allowed("https://other。example/", []string{"shop.example"}); got != CodeSiteNotListed {
		t.Errorf("an unlisted site written with an ideographic full stop: %q", got)
	}
}
