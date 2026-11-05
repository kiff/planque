// Package browserworker is the browser worker (kiff/planque#86, step 2):
// a headless Chromium driven over its debugging pipe with nothing but the
// standard library, one isolated browser context per session, an egress
// proxy the browser goes out through, and the HTTP API planque's client
// calls. It reports facts about a page; planque decides.
package browserworker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/egress"
)

// Proxy is the egress one session's browser goes out through: a CONNECT
// proxy on loopback that resolves each host itself, connects only to
// public addresses (egress.Public), and only to the session's sites when
// it has a list. https only: a plain request is refused. The browser
// connects to the address the proxy chose, by address, so a name rebound
// to a private address after resolution reaches nothing. One proxy per
// session, so a session's sites never apply to another's.
type Proxy struct {
	// Sites are the only hosts (and subdomains) reached; nil, any public site.
	Sites []string
	// Resolve looks a host up; nil, the system resolver. Tests give one.
	Resolve func(ctx context.Context, host string) ([]netip.Addr, error)
	// AllowLocal lets tests reach a shop on loopback, over plain http too.
	AllowLocal bool
	Logger     *slog.Logger

	mu sync.Mutex
	ln net.Listener
	// Refused counts what the proxy refused, for tests and the log.
	Refused int
}

// Start listens on loopback and serves until Close.
func (p *Proxy) Start() (addr string, err error) {
	p.ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	go func() {
		for {
			c, err := p.ln.Accept()
			if err != nil {
				return
			}
			go p.serve(c)
		}
	}()
	return p.ln.Addr().String(), nil
}

// Close stops listening.
func (p *Proxy) Close() error {
	if p.ln == nil {
		return nil
	}
	return p.ln.Close()
}

func (p *Proxy) log() *slog.Logger {
	if p.Logger != nil {
		return p.Logger
	}
	return slog.Default()
}

func (p *Proxy) refuse(c net.Conn, status int, why string) {
	p.mu.Lock()
	p.Refused++
	p.mu.Unlock()
	fmt.Fprintf(c, "HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", status, http.StatusText(status), len(why), why)
	c.Close()
}

// serve handles one connection from the browser: a CONNECT to a site's
// 443, checked, then bytes both ways.
func (p *Proxy) serve(c net.Conn) {
	_ = c.SetReadDeadline(time.Now().Add(30 * time.Second))
	req, err := http.ReadRequest(bufio.NewReader(c))
	if err != nil {
		c.Close()
		return
	}
	_ = c.SetReadDeadline(time.Time{})
	sites := p.Sites
	if req.Method != http.MethodConnect {
		if p.AllowLocal && req.URL != nil && req.URL.Scheme == "http" && isLoopback(req.URL.Hostname()) {
			p.plainLocal(c, req)
			return
		}
		p.refuse(c, http.StatusForbidden, "the browser reaches sites over https only")
		return
	}
	host, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		p.refuse(c, http.StatusBadRequest, "bad host")
		return
	}
	addr, why := p.resolve(req.Context(), host, port, sites)
	if why != "" {
		p.log().Info("egress refused", slog.String("host", host), slog.String("why", why))
		p.refuse(c, http.StatusForbidden, why)
		return
	}
	up, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		p.refuse(c, http.StatusBadGateway, "couldn't reach "+host)
		return
	}
	_, _ = io.WriteString(c, "HTTP/1.1 200 Connection Established\r\n\r\n")
	go func() {
		_, _ = io.Copy(up, c)
		up.Close()
	}()
	_, _ = io.Copy(c, up)
	c.Close()
}

// resolve says where host:port may be reached, as an address, or why not.
func (p *Proxy) resolve(ctx context.Context, host, port string, sites []string) (addr, why string) {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	if p.AllowLocal && isLoopback(host) {
		return net.JoinHostPort(host, port), ""
	}
	if _, err := netip.ParseAddr(host); err == nil {
		return "", "the browser opens sites by name, never an address"
	}
	if len(sites) > 0 && !listed(host, sites) {
		return "", host + " isn't one of the session's sites"
	}
	ips, err := p.lookup(ctx, host)
	if err != nil || len(ips) == 0 {
		return "", host + " couldn't be resolved"
	}
	// Every address must be public: a name that answers with one private
	// address among public ones is refused whole, not retried on another.
	local := true
	for _, ip := range ips {
		if p.AllowLocal && ip.IsLoopback() {
			continue
		}
		local = false
		if !egress.Public(ip) {
			return "", host + " resolves to an address that isn't public"
		}
	}
	if port != "443" && !local {
		return "", "the browser reaches sites on 443 only"
	}
	return net.JoinHostPort(ips[0].Unmap().String(), port), ""
}

func (p *Proxy) lookup(ctx context.Context, host string) ([]netip.Addr, error) {
	if p.Resolve != nil {
		return p.Resolve(ctx, host)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return net.DefaultResolver.LookupNetIP(ctx, "ip", host)
}

func listed(host string, sites []string) bool {
	for _, s := range sites {
		if host == s || strings.HasSuffix(host, "."+s) {
			return true
		}
	}
	return false
}

func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip, err := netip.ParseAddr(host)
	return err == nil && ip.IsLoopback()
}

// plainLocal forwards a plain-http request to a loopback shop: tests only.
func (p *Proxy) plainLocal(c net.Conn, req *http.Request) {
	defer c.Close()
	req.RequestURI = ""
	res, err := (&http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}).Do(req)
	if err != nil {
		p.refuse(c, http.StatusBadGateway, err.Error())
		return
	}
	defer res.Body.Close()
	_ = res.Write(c)
}
