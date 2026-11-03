package planque

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kiff/planque/internal/oauth"
)

func TestClientMetadataIsPublicAndDoesNotUseTheRequestHost(t *testing.T) {
	s := &Server{BaseURL: "https://planque.example"}
	r := httptest.NewRequest("GET", "https://untrusted.example/oauth/client-metadata.json", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	var doc oauth.ClientDocument
	if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil || w.Code != 200 {
		t.Fatalf("%d %s %v", w.Code, w.Body, err)
	}
	if doc.ClientID != s.BaseURL+"/oauth/client-metadata.json" || len(doc.RedirectURIs) != 1 || doc.RedirectURIs[0] != s.redirect() || strings.Contains(w.Body.String(), "untrusted") {
		t.Fatalf("metadata escaped its configured installation: %+v", doc)
	}
	// No KV, credentials or discovery client is configured: CIMD doesn't
	// register dynamically or read a grant just to identify this installation.
	id, err := s.client(context.Background(), oauth.Metadata{ClientMetadata: true})
	if err != nil || id != doc.ClientID {
		t.Fatalf("%q %v", id, err)
	}
}
