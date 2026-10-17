package egress

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"
)

func TestPublic(t *testing.T) {
	refused := []string{"127.0.0.1", "10.1.2.3", "172.16.0.1", "192.168.1.1", "169.254.169.254", "100.64.0.1",
		"0.0.0.0", "0.1.2.3", "224.0.0.1", "255.255.255.255", "192.0.2.1", "198.18.0.1", "::1", "::", "fe80::1",
		"fc00::1", "fd12:3456::1", "::ffff:127.0.0.1", "::ffff:10.0.0.1", "64:ff9b::a00:1", "2001:db8::1", "ff02::1"}
	for _, s := range refused {
		if Public(netip.MustParseAddr(s)) {
			t.Errorf("%s is treated as public", s)
		}
	}
	for _, s := range []string{"8.8.8.8", "54.230.1.1", "2606:4700::1111"} {
		if !Public(netip.MustParseAddr(s)) {
			t.Errorf("%s is refused", s)
		}
	}
}

// The check is on the connection, so it holds whatever the URL says.
func TestClientRefusesLocalServers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("internal")) }))
	defer srv.Close()
	for _, u := range []string{srv.URL, "http://localhost:" + srv.URL[len("http://127.0.0.1:"):]} {
		_, err := Client(5 * time.Second).Get(u)
		if !errors.Is(err, ErrBlocked) {
			t.Errorf("%s: %v, want ErrBlocked", u, err)
		}
	}
}
