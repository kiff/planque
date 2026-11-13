package planque

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

// refusalsLogged logs each refusal of a connection request: the route, the
// entry when the path names one, the status and planque's words, as the
// person's Coach reads them. A taf's failed step then has its reason in
// planque's log too. A 500 is left to fail, which logs its cause. Only
// planque's own answer is logged, never the request's body: it can carry a
// key.
func (s *Server) refusalsLogged(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &answerKept{ResponseWriter: w}
		h(rec, r)
		if rec.status < 400 || rec.status == http.StatusInternalServerError {
			return
		}
		var got struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		_ = json.Unmarshal(rec.body.Bytes(), &got)
		attrs := []any{slog.String("path", loggedPath(r.URL.Path)), slog.Int("status", rec.status), slog.String("error", got.Error)}
		if entry := r.PathValue("entry"); entry != "" {
			attrs = append(attrs, slog.String("entry", entry))
		}
		if got.Code != "" {
			attrs = append(attrs, slog.String("code", got.Code))
		}
		s.log().WarnContext(r.Context(), "connection refused", attrs...)
	}
}

// answerKept passes an answer on and keeps its status, and the start of a
// refusal's body to log.
type answerKept struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (a *answerKept) WriteHeader(status int) {
	if a.status == 0 {
		a.status = status
	}
	a.ResponseWriter.WriteHeader(status)
}

func (a *answerKept) Write(b []byte) (int, error) {
	if a.status == 0 {
		a.status = http.StatusOK
	}
	if a.status >= 400 && a.body.Len() < 2<<10 {
		a.body.Write(b[:min(len(b), 2<<10-a.body.Len())])
	}
	return a.ResponseWriter.Write(b)
}
