package planque

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
)

// A connect that ended failed answers with its reason, for the asking agent to read and fix (Coach reads it as
// "Not connected: <reason>"), never as planque's "could not complete this request" (a resource type found only at
// approve, 21 Nov).
func TestAFailedConnectAnswersWithItsReason(t *testing.T) {
	s := &Server{}
	for _, err := range []error{connectFailed("x's resource must be a string in its schema"), fmt.Errorf("approve: %w", connectFailed("x's resource must be a string in its schema"))} {
		w := httptest.NewRecorder()
		s.fail(w, httptest.NewRequest("POST", "/v1/accounts/t/grants/g/approve", nil), err)
		var out map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		if w.Code != 409 || out["code"] != "connect_failed" || out["error"] != "x's resource must be a string in its schema" {
			t.Fatalf("%v: %d %v", err, w.Code, out)
		}
	}
}
