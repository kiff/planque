package planque

import (
	"net/http/httptest"
	"testing"
)

func TestRobotsKeepsCrawlersOff(t *testing.T) {
	w := httptest.NewRecorder()
	(&Server{}).Handler().ServeHTTP(w, httptest.NewRequest("GET", "/robots.txt", nil))
	if w.Code != 200 || w.Body.String() != "User-agent: *\nDisallow: /\n" || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("%d %q %q", w.Code, w.Header().Get("Content-Type"), w.Body)
	}
}
