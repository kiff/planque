package browserworker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

// chrome is one headless Chromium, driven over its remote-debugging pipe
// (fds 3 and 4, JSON messages ended by a NUL): the Chrome DevTools
// Protocol with nothing but the standard library.
type chrome struct {
	cmd     *exec.Cmd
	w       *os.File
	dir     string
	next    atomic.Int64
	mu      sync.Mutex
	pending map[int64]chan cdpMessage
	events  []cdpMessage // kept for waits, bounded
	eventCh chan struct{}
	done    chan struct{}
	err     error
	// handlers run for every event, each in a goroutine of its own, so
	// one may call the browser back (the reader never waits on them).
	handlers []func(cdpMessage)
}

// handle adds a handler for the browser's events.
func (c *chrome) handle(h func(cdpMessage)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers = append(c.handlers, h)
}

type cdpMessage struct {
	ID        int64           `json:"id,omitempty"`
	SessionID string          `json:"sessionId,omitempty"`
	Method    string          `json:"method,omitempty"`
	Params    json.RawMessage `json:"params,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// chromePath finds the browser: PLANQUE_CHROME, then the usual names, then
// the Playwright install.
func chromePath(configured string) (string, error) {
	candidates := []string{configured, os.Getenv("PLANQUE_CHROME")}
	for _, name := range []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "chrome"} {
		if p, err := exec.LookPath(name); err == nil {
			candidates = append(candidates, p)
		}
	}
	candidates = append(candidates, "/opt/pw-browsers/chromium", "/usr/bin/chromium")
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", errors.New("no Chromium found: set PLANQUE_CHROME")
}

// launch starts a Chromium that goes out only through proxy, with a fresh
// profile directory it deletes on exit.
func launch(ctx context.Context, path, proxy string) (*chrome, error) {
	dir, err := os.MkdirTemp("", "planque-browser-")
	if err != nil {
		return nil, err
	}
	toChrome, w, err := os.Pipe()
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	r, fromChrome, err := os.Pipe()
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	cmd := exec.CommandContext(ctx, path, launchArgs(dir, proxy)...)
	cmd.ExtraFiles = []*os.File{toChrome, fromChrome}
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "HOME="+dir)
	if err := cmd.Start(); err != nil {
		os.RemoveAll(dir)
		return nil, fmt.Errorf("start the browser: %w", err)
	}
	toChrome.Close()
	fromChrome.Close()
	c := &chrome{cmd: cmd, w: w, dir: dir, pending: map[int64]chan cdpMessage{}, eventCh: make(chan struct{}, 1), done: make(chan struct{})}
	go c.read(bufio.NewReaderSize(r, 1<<20))
	return c, nil
}

// launchArgs are the browser's flags: headless, a profile of its own, and
// every way out through the proxy. WebRTC's UDP (a page's STUN or TURN
// candidates) would go around an HTTP proxy, so it is confined to what
// the proxy carries, which is nothing: a page gathers no candidate and
// probes no address.
func launchArgs(dir, proxy string) []string {
	return []string{
		"--headless=new", "--no-sandbox", "--disable-gpu", "--disable-dev-shm-usage",
		"--no-first-run", "--no-default-browser-check", "--disable-extensions", "--disable-sync",
		"--disable-background-networking", "--disable-component-update", "--disable-default-apps",
		"--disable-translate", "--metrics-recording-only", "--mute-audio", "--hide-scrollbars",
		"--window-size=1280,900", "--lang=en-US",
		"--remote-debugging-pipe",
		"--user-data-dir=" + dir,
		"--proxy-server=http://" + proxy,
		// Loopback too goes through the proxy: nothing bypasses it.
		"--proxy-bypass-list=<-loopback>",
		"--force-webrtc-ip-handling-policy=disable_non_proxied_udp",
		// Every site in a process of its own, cross-site frames included:
		// the boundary the per-target script relies on, and the stronger
		// isolation anyway.
		"--site-per-process",
		"about:blank",
	}
}

// read dispatches what the browser sends: answers to their callers, events
// to the kept list.
func (c *chrome) read(r *bufio.Reader) {
	defer close(c.done)
	for {
		line, err := r.ReadBytes(0)
		if err != nil {
			c.mu.Lock()
			c.err = errors.New("the browser ended")
			for id, ch := range c.pending {
				close(ch)
				delete(c.pending, id)
			}
			c.mu.Unlock()
			return
		}
		var m cdpMessage
		if json.Unmarshal(line[:len(line)-1], &m) != nil {
			continue
		}
		c.mu.Lock()
		if m.ID != 0 {
			if ch, ok := c.pending[m.ID]; ok {
				ch <- m
				delete(c.pending, m.ID)
			}
		} else if m.Method != "" {
			c.events = append(c.events, m)
			if len(c.events) > 500 {
				c.events = c.events[len(c.events)-500:]
			}
			select {
			case c.eventCh <- struct{}{}:
			default:
			}
			for _, h := range c.handlers {
				go h(m)
			}
		}
		c.mu.Unlock()
	}
}

// call sends one method and waits for its answer.
func (c *chrome) call(ctx context.Context, session, method string, params any, out any) error {
	id := c.next.Add(1)
	msg := map[string]any{"id": id, "method": method}
	if params != nil {
		msg["params"] = params
	}
	if session != "" {
		msg["sessionId"] = session
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ch := make(chan cdpMessage, 1)
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return c.err
	}
	c.pending[id] = ch
	c.mu.Unlock()
	if _, err := c.w.Write(append(b, 0)); err != nil {
		return fmt.Errorf("the browser isn't listening: %w", err)
	}
	select {
	case m, ok := <-ch:
		if !ok {
			return errors.New("the browser ended")
		}
		if m.Error != nil {
			return fmt.Errorf("%s: %s", method, m.Error.Message)
		}
		if out != nil && len(m.Result) > 0 {
			return json.Unmarshal(m.Result, out)
		}
		return nil
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return ctx.Err()
	}
}

// waitEvent waits for an event of a session, from the events kept since
// index from; it returns the index after it, for the next wait.
func (c *chrome) waitEvent(ctx context.Context, session, method string, from int) (json.RawMessage, int, error) {
	for {
		c.mu.Lock()
		for i := from; i < len(c.events); i++ {
			e := c.events[i]
			if e.Method == method && (session == "" || e.SessionID == session) {
				c.mu.Unlock()
				return e.Params, i + 1, nil
			}
		}
		from = len(c.events)
		c.mu.Unlock()
		select {
		case <-c.eventCh:
		case <-c.done:
			return nil, from, errors.New("the browser ended")
		case <-ctx.Done():
			return nil, from, ctx.Err()
		}
	}
}

// eventsSeen is how many events are kept now: the point a wait starts from.
func (c *chrome) eventsSeen() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.events)
}

// stop ends the browser and deletes its profile: nothing it wrote stays.
func (c *chrome) stop() {
	_ = c.call(context.Background(), "", "Browser.close", nil, nil)
	select {
	case <-c.done:
	case <-time.After(2 * time.Second):
		_ = c.cmd.Process.Kill()
	}
	_ = c.cmd.Wait()
	c.w.Close()
	_ = os.RemoveAll(c.dir)
}
