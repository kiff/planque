package browserworker

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kiff/planque/internal/browser"
)

// Server is the worker's API, which planque's client calls: a session per
// call, opened, stepped, screenshotted and closed. Every request carries
// the shared key; a session is its caller's for its life and ends on its
// own after idle time, so a planque that died mid-call leaves no browser.
//
//	POST   /sessions                  {grant, state, sites} → {session}
//	POST   /sessions/{id}/open        {url}                 → {page}
//	POST   /sessions/{id}/step        {step}                → {page}
//	POST   /sessions/{id}/screenshot                        → {png (base64)}
//	DELETE /sessions/{id}                                   → {state}
type Server struct {
	Local *Local
	// KeySHA256 is the SHA-256, hex, of the key planque presents in
	// KeyHeader. The key itself is only in planque's environment.
	KeySHA256 string
	// Idle is how long a session may go unused before the worker ends it.
	Idle   time.Duration
	Logger *slog.Logger

	mu       sync.Mutex
	sessions map[string]*held
}

// KeyHeader carries the shared key.
const KeyHeader = "Planque-Worker-Key"

type held struct {
	sess browser.Session
	last time.Time
	mu   sync.Mutex // one call at a time per session
}

func (s *Server) log() *slog.Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// Handler is the API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]bool{"ok": true}) })
	mux.HandleFunc("POST /sessions", s.keyed(s.start))
	mux.HandleFunc("POST /sessions/{id}/open", s.keyed(s.open))
	mux.HandleFunc("POST /sessions/{id}/step", s.keyed(s.step))
	mux.HandleFunc("POST /sessions/{id}/screenshot", s.keyed(s.screenshot))
	mux.HandleFunc("DELETE /sessions/{id}", s.keyed(s.close))
	return mux
}

// keyed checks the shared key before anything is read.
func (s *Server) keyed(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sum := sha256.Sum256([]byte(strings.TrimSpace(r.Header.Get(KeyHeader))))
		want := strings.ToLower(s.KeySHA256)
		got := hexOf(sum[:])
		if want == "" || len(got) != len(want) || subtle.ConstantTimeCompare([]byte(got), []byte(want)) != 1 {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "the worker's key is required"})
			return
		}
		h(w, r)
	}
}

func hexOf(b []byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, c := range b {
		out[i*2], out[i*2+1] = digits[c>>4], digits[c&15]
	}
	return string(out)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decode(r *http.Request, v any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, 4<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// startBody is what a session begins with.
type startBody struct {
	Grant string          `json:"grant"`
	State json.RawMessage `json:"state,omitempty"`
	Sites []string        `json:"sites,omitempty"`
}

func (s *Server) start(w http.ResponseWriter, r *http.Request) {
	var in startBody
	if err := decode(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "bad body")
		return
	}
	sess, err := s.Local.Start(r.Context(), browser.Start{Grant: in.Grant, State: in.State, Sites: in.Sites})
	if err != nil {
		s.log().Warn("session not started", slog.String("err", err.Error()))
		writeErr(w, http.StatusServiceUnavailable, "the browser couldn't be started")
		return
	}
	id := rand.Text()
	s.mu.Lock()
	if s.sessions == nil {
		s.sessions = map[string]*held{}
	}
	s.sessions[id] = &held{sess: sess, last: time.Now()}
	s.mu.Unlock()
	s.reapLater(id)
	writeJSON(w, http.StatusCreated, map[string]string{"session": id})
}

// reapLater ends the session when it has gone idle.
func (s *Server) reapLater(id string) {
	idle := s.Idle
	if idle == 0 {
		idle = 2 * time.Minute
	}
	go func() {
		for {
			time.Sleep(idle / 4)
			s.mu.Lock()
			h, ok := s.sessions[id]
			if !ok {
				s.mu.Unlock()
				return
			}
			if time.Since(h.last) < idle {
				s.mu.Unlock()
				continue
			}
			delete(s.sessions, id)
			s.mu.Unlock()
			h.mu.Lock()
			_, _ = h.sess.Close(context.Background())
			h.mu.Unlock()
			s.log().Info("idle session ended", slog.String("session", id))
			return
		}
	}()
}

// take is the session, held for one call.
func (s *Server) take(w http.ResponseWriter, r *http.Request) (*held, bool) {
	s.mu.Lock()
	h, ok := s.sessions[r.PathValue("id")]
	if ok {
		h.last = time.Now()
	}
	s.mu.Unlock()
	if !ok {
		writeErr(w, http.StatusNotFound, "no such session")
		return nil, false
	}
	h.mu.Lock()
	return h, true
}

func (s *Server) open(w http.ResponseWriter, r *http.Request) {
	var in struct {
		URL string `json:"url"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "bad body")
		return
	}
	h, ok := s.take(w, r)
	if !ok {
		return
	}
	defer h.mu.Unlock()
	p, err := h.sess.Open(r.Context(), in.URL)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"page": p})
}

func (s *Server) step(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Step browser.Step `json:"step"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "bad body")
		return
	}
	h, ok := s.take(w, r)
	if !ok {
		return
	}
	defer h.mu.Unlock()
	p, err := h.sess.Do(r.Context(), in.Step)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"page": p})
}

func (s *Server) screenshot(w http.ResponseWriter, r *http.Request) {
	h, ok := s.take(w, r)
	if !ok {
		return
	}
	defer h.mu.Unlock()
	png, err := h.sess.Screenshot(r.Context())
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"png": png})
}

func (s *Server) close(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.mu.Lock()
	h, ok := s.sessions[id]
	delete(s.sessions, id)
	s.mu.Unlock()
	if !ok {
		writeErr(w, http.StatusNotFound, "no such session")
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	st, err := h.sess.Close(r.Context())
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"state": json.RawMessage(st)})
}

// Shutdown ends every session.
func (s *Server) Shutdown(ctx context.Context) {
	s.mu.Lock()
	all := s.sessions
	s.sessions = nil
	s.mu.Unlock()
	for _, h := range all {
		h.mu.Lock()
		_, _ = h.sess.Close(ctx)
		h.mu.Unlock()
	}
}
