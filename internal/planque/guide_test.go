package planque

import (
	"strings"
	"testing"
)

// An app's adapted guide is served by its catalog entry, to an agent of the
// account: a provider's app by its toolkit, a reviewed entry by its name.
// An app with no adapted guide yet, or a name that isn't an app's, is a 404
// with words; another account's agent is refused.
func TestAnAppsGuideIsServedByItsEntry(t *testing.T) {
	e := newEnv(t)
	base := "/v1/accounts/" + tenant + "/catalog/"
	for entry, app := range map[string]string{"connector:composio:gmail": "gmail", "connector:composio:googlesheets": "googlesheets", StripeEntry: "stripe", "qonto-key@2": "qonto"} {
		st, out := e.do("GET", base+entry+"/guide", nil, nil, "k_setup")
		if st != 200 || out["app"] != app || out["checked"] == "" || !strings.HasPrefix(out["guide"].(string), "# ") {
			t.Fatalf("%s: %d %v", entry, st, out["app"])
		}
	}
	for _, entry := range []string{"connector:composio:no_such_app", "README", "apps", "connector:composio:..%2Fgmail"} {
		if st, out := e.do("GET", base+entry+"/guide", nil, nil, "k_setup"); st != 404 {
			t.Fatalf("%s: %d %v", entry, st, out)
		}
	}
	if st, _ := e.do("GET", base+"connector:composio:gmail/guide", nil, nil, "k_other"); st == 200 {
		t.Fatal("another account's agent read a guide on this account")
	}
}
