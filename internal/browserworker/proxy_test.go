package browserworker

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// connect asks the proxy for a tunnel to host:port and returns the status
// line it answered.
func connect(t *testing.T, proxy, hostport string) string {
	t.Helper()
	c, err := net.DialTimeout("tcp", proxy, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_, _ = c.Write([]byte("CONNECT " + hostport + " HTTP/1.1\r\nHost: " + hostport + "\r\n\r\n"))
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	line, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(line)
}

// A listed name that resolves to a private address is never fetched: the
// proxy resolves it itself, refuses it, and nothing connects. Neither is
// an address written out, a name off the sites, a port other than 443, or
// a plain-http request. (agentgabo's ask on kiff/planque#88 for step 2.)
func TestProxyNeverReachesAPrivateAddressWhateverTheName(t *testing.T) {
	var touched atomic.Int32
	trap, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer trap.Close()
	go func() {
		for {
			c, err := trap.Accept()
			if err != nil {
				return
			}
			touched.Add(1)
			c.Close()
		}
	}()
	trapPort := trap.Addr().(*net.TCPAddr).Port
	answers := map[string][]netip.Addr{
		"shop.example":       {netip.MustParseAddr("10.0.0.7")},
		"meta.shop.example":  {netip.MustParseAddr("169.254.169.254")},
		"mixed.shop.example": {netip.MustParseAddr("93.184.216.34"), netip.MustParseAddr("192.168.1.9")},
		"loop.shop.example":  {netip.MustParseAddr("127.0.0.1")},
		"six.shop.example":   {netip.MustParseAddr("fd12::1")},
		"gone.shop.example":  nil,
	}
	p := &Proxy{Sites: []string{"shop.example"}, Resolve: func(_ context.Context, host string) ([]netip.Addr, error) {
		ips, ok := answers[host]
		if !ok {
			return nil, errors.New("no such host")
		}
		return ips, nil
	}}
	addr, err := p.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	for _, c := range []struct{ target, want string }{
		{"shop.example:443", "403"},            // resolves to 10.x
		{"meta.shop.example:443", "403"},       // the cloud metadata address
		{"mixed.shop.example:443", "403"},      // one private answer among public ones
		{"loop.shop.example:443", "403"},       // loopback
		{"six.shop.example:443", "403"},        // a private IPv6
		{"gone.shop.example:443", "403"},       // doesn't resolve
		{"other.example:443", "403"},           // off the sites, never resolved
		{"10.0.0.7:443", "403"},                // an address written out
		{"127.0.0.1:" + itoa(trapPort), "403"}, // loopback written out, on a listening port
		{"shop.example:8443", "403"},           // not 443
	} {
		status := connect(t, addr, c.target)
		if !strings.Contains(status, c.want) {
			t.Errorf("CONNECT %s: %s, want %s", c.target, status, c.want)
		}
	}
	// A plain-http request through the proxy is refused too.
	c, _ := net.Dial("tcp", addr)
	_, _ = c.Write([]byte("GET http://shop.example/ HTTP/1.1\r\nHost: shop.example\r\n\r\n"))
	line, _ := bufio.NewReader(c).ReadString('\n')
	c.Close()
	if !strings.Contains(line, "403") {
		t.Errorf("plain http: %s", line)
	}
	if n := touched.Load(); n != 0 {
		t.Fatalf("something connected to the trap %d times", n)
	}
	if p.Refused != 11 {
		t.Errorf("refused %d, want 11", p.Refused)
	}
}

// A listed name that resolves to a public address is tunnelled to that
// address, by address; a site not listed is refused before resolving.
func TestProxyTunnelsToThePublicAddressItResolved(t *testing.T) {
	// A server on loopback stands in for the public address: the proxy
	// dials what the resolver answered, and AllowLocal lets it be loopback
	// here. The resolver is asked only for the listed name.
	up, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer up.Close()
	go func() {
		for {
			c, err := up.Accept()
			if err != nil {
				return
			}
			_, _ = c.Write([]byte("hello from upstream"))
			c.Close()
		}
	}()
	var asked []string
	p := &Proxy{Sites: []string{"shop.example"}, AllowLocal: true, Resolve: func(_ context.Context, host string) ([]netip.Addr, error) {
		asked = append(asked, host)
		return []netip.Addr{netip.MustParseAddr("127.0.0.1")}, nil
	}}
	// With AllowLocal the loopback answer passes the public check; the
	// port is the upstream's, so the test gives the real one.
	addr, err := p.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	port := up.Addr().(*net.TCPAddr).Port
	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	target := "shop.example:" + itoa(port)
	_, _ = c.Write([]byte("CONNECT " + target + " HTTP/1.1\r\nHost: " + target + "\r\n\r\n"))
	rd := bufio.NewReader(c)
	line, _ := rd.ReadString('\n')
	if !strings.Contains(line, "200") {
		t.Fatalf("CONNECT: %s", line)
	}
	for l, _ := rd.ReadString('\n'); l != "\r\n" && l != ""; l, _ = rd.ReadString('\n') {
	}
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	body, _ := rd.ReadString('m')
	if !strings.HasPrefix(body, "hello") {
		t.Errorf("through the tunnel: %q", body)
	}
	if len(asked) != 1 || asked[0] != "shop.example" {
		t.Errorf("resolved: %v", asked)
	}
	if status := connect(t, addr, "other.example:"+itoa(port)); !strings.Contains(status, "403") || len(asked) != 1 {
		t.Errorf("an unlisted site: %s (resolved %v)", status, asked)
	}
}

// Tests only: a plain-http shop on loopback is reachable with AllowLocal,
// so the real browser can be run against an httptest server.
func TestProxyForwardsALocalShopForTests(t *testing.T) {
	shop := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html><title>Shop</title></html>"))
	}))
	defer shop.Close()
	p := &Proxy{AllowLocal: true}
	addr, err := p.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	c, _ := net.Dial("tcp", addr)
	defer c.Close()
	_, _ = c.Write([]byte("GET " + shop.URL + "/ HTTP/1.1\r\nHost: " + strings.TrimPrefix(shop.URL, "http://") + "\r\n\r\n"))
	res, err := http.ReadResponse(bufio.NewReader(c), nil)
	if err != nil || res.StatusCode != 200 {
		t.Fatalf("local shop: %v %v", err, res)
	}
	// Without AllowLocal, the same is refused.
	q := &Proxy{}
	qa, _ := q.Start()
	defer q.Close()
	c2, _ := net.Dial("tcp", qa)
	defer c2.Close()
	_, _ = c2.Write([]byte("GET " + shop.URL + "/ HTTP/1.1\r\nHost: " + strings.TrimPrefix(shop.URL, "http://") + "\r\n\r\n"))
	line, _ := bufio.NewReader(c2).ReadString('\n')
	if !strings.Contains(line, "403") {
		t.Errorf("local shop in production: %s", line)
	}
}

func itoa(n int) string { return strconv.Itoa(n) }
