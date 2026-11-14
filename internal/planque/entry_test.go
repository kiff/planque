package planque

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/mcp"
)

// installEntries is the reviewed catalog as planque matches grants to it:
// each entry's server, and how it's reached.
func installEntries(e *env, entries ...CatalogEntry) {
	e.s.CatalogEntries = map[string]CatalogEntry{}
	for _, c := range entries {
		e.s.CatalogEntries[c.ID] = c
	}
}

// listings numbers each listing, so two in the same second are two signed
// requests, not one replayed.
var listings atomic.Int64

// listedByPhone is the grants list as the owner's app reads it, by id.
func listedByPhone(t *testing.T, e *env) map[string]map[string]any {
	t.Helper()
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants?n="+strconv.FormatInt(listings.Add(1), 10), nil, e.phone, "")
	if st != 200 {
		t.Fatalf("list: %d %v", st, out)
	}
	byID := map[string]map[string]any{}
	for _, it := range out["grants"].([]any) {
		v := it.(map[string]any)
		byID[v["id"].(string)] = v
	}
	return byID
}

// stored is a grant as kept: its request's exact bytes, and its version.
func stored(t *testing.T, e *env, id string) ([]byte, int64) {
	t.Helper()
	it, ok, err := e.kv.Get(context.Background(), tenantPK(tenant), skGrant(id))
	if err != nil || !ok {
		t.Fatalf("grant %s: %v %v", id, ok, err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(it.Doc, &doc); err != nil {
		t.Fatal(err)
	}
	return doc["request"], it.Ver
}

func refundsEntry(e *env, id string) CatalogEntry {
	return CatalogEntry{ID: id, Vendor: "Refunds", Title: "Refunds", Server: e.vendor.MCPURL()}
}

// A grant made before requests named their app gets the entry recorded
// beside its request: the request the owner signed is kept byte for byte,
// its digest still matches the approval, and the backfill writes once.
func TestBackfillRecordsTheAppBesideTheSignedRequest(t *testing.T) {
	e := newEnv(t)
	g := e.connect() // no catalog yet: a grant as made before this change
	if g.Request.CatalogEntry != "" || g.Entry != "" {
		t.Fatalf("an old grant has an entry already: %q %q", g.Request.CatalogEntry, g.Entry)
	}
	before, _ := stored(t, e, g.ID)

	installEntries(e, refundsEntry(e, "refunds@1"))
	listed := listedByPhone(t, e)[g.ID]
	if listed["entry"] != "refunds@1" {
		t.Fatalf("the app lists it by entry: %v", listed["entry"])
	}
	if req := listed["request"].(map[string]any); req["catalog_entry"] != nil {
		t.Fatalf("the backfill went into the signed request: %v", req)
	}
	after, ver := stored(t, e, g.ID)
	if !bytes.Equal(before, after) {
		t.Fatalf("the signed request changed:\n%s\n%s", before, after)
	}
	back, err := e.s.loadGrant(context.Background(), tenant, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if back.Entry != "refunds@1" || back.Request.CatalogEntry != "" {
		t.Fatalf("not recorded on the grant: %q %q", back.Entry, back.Request.CatalogEntry)
	}
	d, err := digest(back.Request)
	if err != nil || d != g.Digest || back.Digest != g.Digest || listed["digest"] != g.Digest {
		t.Fatalf("digest: %s, kept %s, was %s, listed %v", d, back.Digest, g.Digest, listed["digest"])
	}
	if back.Approval == nil || !strings.Contains(back.Approval.Body, g.Digest) {
		t.Fatalf("the approval no longer covers the request's digest: %+v", back.Approval)
	}

	// Idempotent: listing again writes nothing, and says the same.
	if again := listedByPhone(t, e)[g.ID]; again["entry"] != "refunds@1" {
		t.Fatalf("second listing: %v", again["entry"])
	}
	if _, v := stored(t, e, g.ID); v != ver {
		t.Fatalf("the backfill wrote again: version %d, then %d", ver, v)
	}
	// The connection keeps working as signed.
	if st, out := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup"); st != 200 || out["entry"] != "refunds@1" || out["status"] != StatusConnected {
		t.Fatalf("view: %d %v", st, out)
	}
}

// Where the app is unclear, the grant is left alone: no entry for the same
// server twice in the catalog, and none for a server it doesn't have.
func TestBackfillLeavesAGrantAloneWhenItsAppIsUnclear(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	before, ver := stored(t, e, g.ID)
	for name, entries := range map[string][]CatalogEntry{
		"two entries for its server": {refundsEntry(e, "refunds@1"), refundsEntry(e, "refunds@2")},
		"none for its server":        {{ID: "other@1", Vendor: "Other", Title: "Other", Server: "https://mcp.other.example/mcp"}},
		"none at all":                nil,
	} {
		installEntries(e, entries...)
		if listed := listedByPhone(t, e)[g.ID]; listed["entry"] != nil {
			t.Fatalf("%s: listed under %v", name, listed["entry"])
		}
		after, v := stored(t, e, g.ID)
		if !bytes.Equal(before, after) || v != ver {
			t.Fatalf("%s: the grant was written (version %d, then %d)", name, ver, v)
		}
		if back, _ := e.s.loadGrant(context.Background(), tenant, g.ID); back.Entry != "" {
			t.Fatalf("%s: guessed %q", name, back.Entry)
		}
	}
}

// A managed app's grant made before its request named it gets its
// provider's toolkit as its entry; its request stays as signed.
func TestBackfillRecordsAManagedAppsToolkit(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	g := &Grant{ID: "grt_oldmanaged", Tenant: tenant, Status: StatusConnected, VendorName: "Stripe",
		Request: Request{Version: "planque.grant.v1", Grant: "grt_oldmanaged", Tenant: tenant, Server: "https://backend.composio.dev",
			Agent: "agt_taf", Tools: []ToolRequest{}, AskedBy: "agt_setup", CreatedAt: "2026-10-01T09:00:00Z", Kind: "connector",
			Connector: &ConnectorRequest{Provider: "composio", Toolkit: "stripe", Title: "Stripe"}},
		Events: []Event{{At: time.Now().UTC(), What: "connected"}}}
	g.Digest, _ = digest(g.Request)
	if err := e.s.saveGrant(context.Background(), g); err != nil {
		t.Fatal(err)
	}
	before, _ := stored(t, e, g.ID)
	if listed := listedByPhone(t, e)[g.ID]; listed["entry"] != "connector:composio:stripe" {
		t.Fatalf("entry: %v", listed["entry"])
	}
	after, _ := stored(t, e, g.ID)
	back, _ := e.s.loadGrant(context.Background(), tenant, g.ID)
	if !bytes.Equal(before, after) || back.Entry != "connector:composio:stripe" || back.Request.CatalogEntry != "" {
		t.Fatalf("recorded %q, request %q, unchanged %v", back.Entry, back.Request.CatalogEntry, bytes.Equal(before, after))
	}
}

// A grant asked by server from now on names its app in the request the
// owner signs, and the phone approves it as before.
func TestANewGrantNamesItsAppInTheSignedRequest(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	installEntries(e, refundsEntry(e, "refunds@1"))
	id, d := e.request()
	g, err := e.s.loadGrant(context.Background(), tenant, id)
	if err != nil {
		t.Fatal(err)
	}
	if g.Request.CatalogEntry != "refunds@1" || g.Entry != "" {
		t.Fatalf("request %q, grant %q", g.Request.CatalogEntry, g.Entry)
	}
	if signed, _ := digest(g.Request); signed != d || g.Digest != d {
		t.Fatal("the digest the phone signs doesn't cover the entry")
	}
	if st, out := e.approveAndSignIn(id, d); st != 200 || out["status"] != StatusConnected || out["entry"] != "refunds@1" {
		t.Fatalf("connect: %d %v", st, out)
	}
	// A server the catalog doesn't have: no entry, and nothing guessed.
	installEntries(e, CatalogEntry{ID: "other@1", Vendor: "Other", Title: "Other", Server: "https://mcp.other.example/mcp"})
	id, _ = e.request()
	if g, _ := e.s.loadGrant(context.Background(), tenant, id); g.Request.CatalogEntry != "" {
		t.Fatalf("guessed %q", g.Request.CatalogEntry)
	}
}

// The catalog shows an app as connected when a connection Coach started to
// it is, as for one asked from the catalog.
func TestTheCatalogCountsACoachStartedConnection(t *testing.T) {
	e := newEnv(t)
	g := e.connect()
	installEntries(e, refundsEntry(e, "refunds@1"))
	st, out := e.do("GET", "/v1/accounts/"+tenant+"/catalog", nil, nil, "k_setup")
	if st != 200 {
		t.Fatalf("catalog: %d %v", st, out)
	}
	entries := out["entries"].([]any)
	if len(entries) != 3 || entries[0].(map[string]any)["entry"] != QontoEntry || entries[2].(map[string]any)["entry"] != StripeEntry {
		t.Fatalf("entries: %v", entries)
	}
	// planque's own Stripe and Qonto entries carry their logos, as the managed apps do.
	if icon := entries[2].(map[string]any)["icon"]; icon != stripeIcon || !catalogHTTPS(stripeIcon) {
		t.Fatalf("Stripe's icon: %v", icon)
	}
	if icon := entries[0].(map[string]any)["icon"]; icon != qontoIcon || !catalogHTTPS(qontoIcon) {
		t.Fatalf("Qonto's icon: %v", icon)
	}
	if v := entries[1].(map[string]any); v["entry"] != "refunds@1" || v["state"] != "connected" || v["grant"] != g.ID {
		t.Fatalf("state: %v", v)
	}
}

func TestSameServer(t *testing.T) {
	for _, c := range []struct {
		a, b string
		same bool
	}{
		{"https://mcp.stripe.com", "https://mcp.stripe.com/", true},
		{"https://MCP.Stripe.com/mcp", "HTTPS://mcp.stripe.com/mcp/", true},
		{"https://mcp.stripe.com:443/mcp", "https://mcp.stripe.com/mcp", true},
		{"https://mcp.stripe.com/mcp", "https://mcp.stripe.com/v2/mcp", false},
		{"https://mcp.stripe.com", "https://stripe.attacker.com", false},
		{"https://mcp.stripe.com:8443", "https://mcp.stripe.com", false},
		{"http://mcp.stripe.com", "https://mcp.stripe.com", false},
		{"https://mcp.stripe.com/mcp?x=1", "https://mcp.stripe.com/mcp", false},
		{"https://user@mcp.stripe.com", "https://mcp.stripe.com", false},
		{"", "https://mcp.stripe.com", false},
	} {
		if got := sameServer(c.a, c.b); got != c.same {
			t.Errorf("sameServer(%q, %q) = %v", c.a, c.b, got)
		}
	}
}

// vendorOrigin is the fake vendor's origin, as a row of the reviewed list.
func vendorOrigin(t *testing.T, e *env) string {
	t.Helper()
	u, err := url.Parse(e.vendor.MCPURL())
	if err != nil {
		t.Fatal(err)
	}
	return u.Scheme + "://" + u.Host
}

// A vendor's own server on the reviewed list gets the managed app the list
// names: an older grant beside its request, a new one in it. A server not on
// the list stays without an entry.
func TestTheServerListNamesAVendorsOwnServer(t *testing.T) {
	e := newEnv(t)
	old := e.connect()
	before, _ := stored(t, e, old.ID)

	e.s.ServerApps = map[string]string{"https://mcp.elsewhere.example": "connector:composio:other"}
	if listed := listedByPhone(t, e)[old.ID]; listed["entry"] != nil {
		t.Fatalf("a server not on the list got %v", listed["entry"])
	}

	e.s.ServerApps = map[string]string{vendorOrigin(t, e): "connector:composio:refunds"}
	if listed := listedByPhone(t, e)[old.ID]; listed["entry"] != "connector:composio:refunds" {
		t.Fatalf("listed under %v", listed["entry"])
	}
	after, _ := stored(t, e, old.ID)
	back, _ := e.s.loadGrant(context.Background(), tenant, old.ID)
	if !bytes.Equal(before, after) || back.Entry != "connector:composio:refunds" || back.Request.CatalogEntry != "" {
		t.Fatalf("recorded %q, request %q, unchanged %v", back.Entry, back.Request.CatalogEntry, bytes.Equal(before, after))
	}

	id, d := e.request()
	g, _ := e.s.loadGrant(context.Background(), tenant, id)
	if g.Request.CatalogEntry != "connector:composio:refunds" {
		t.Fatalf("a new grant names %q", g.Request.CatalogEntry)
	}
	if signed, _ := digest(g.Request); signed != d {
		t.Fatal("the digest doesn't cover the entry")
	}
}

// The list says which app a server is, whatever the catalog holds for it;
// two rows that match with different apps say nothing.
func TestTheServerListComesFirstAndNeverGuesses(t *testing.T) {
	e := newEnv(t)
	e.enrol()
	installEntries(e, refundsEntry(e, "refunds@1"))
	e.s.ServerApps = map[string]string{vendorOrigin(t, e): "connector:composio:refunds"}
	g := &Grant{Request: Request{Server: e.vendor.MCPURL()}}
	if got := e.s.deriveEntry(g); got != "connector:composio:refunds" {
		t.Fatalf("the list didn't come first: %q", got)
	}
	e.s.ServerApps[e.vendor.MCPURL()] = "connector:composio:another"
	if got := e.s.deriveEntry(g); got != "" {
		t.Fatalf("two rows with two apps gave %q", got)
	}
	e.s.ServerApps[e.vendor.MCPURL()] = "connector:composio:refunds"
	if got := e.s.deriveEntry(g); got != "connector:composio:refunds" {
		t.Fatalf("two rows with one app gave %q", got)
	}
	// A key grant to the same server is the same app.
	if got := e.s.deriveEntry(&Grant{Request: Request{Kind: KindAPIKey, Server: e.vendor.MCPURL()}}); got != "connector:composio:refunds" {
		t.Fatalf("a key grant gave %q", got)
	}
}

func TestServerListRows(t *testing.T) {
	for _, c := range []struct {
		row, server string
		match       bool
	}{
		{"https://mcp.stripe.com", "https://mcp.stripe.com", true},
		{"https://mcp.stripe.com", "https://mcp.stripe.com/mcp", true},
		{"https://mcp.stripe.com/", "https://MCP.stripe.com/v1/sse", true},
		{"https://mcp.stripe.com/mcp", "https://mcp.stripe.com/mcp/", true},
		{"https://mcp.stripe.com/mcp", "https://mcp.stripe.com/other", false},
		{"https://mcp.stripe.com", "https://stripe.attacker.com/mcp", false},
		{"https://mcp.stripe.com", "https://mcp.stripe.com.attacker.com/mcp", false},
		{"https://mcp.stripe.com", "http://mcp.stripe.com/mcp", false},
		{"https://mcp.stripe.com", "https://mcp.stripe.com:8443/mcp", false},
		{"https://mcp.stripe.com", "https://user@mcp.stripe.com/mcp", false},
		{"https://mcp.stripe.com", "", false},
	} {
		if got := rowMatches(c.row, c.server); got != c.match {
			t.Errorf("rowMatches(%q, %q) = %v", c.row, c.server, got)
		}
	}
}

// The list deploy/aws.sh ships with the function decodes, so a bad edit to
// it fails here, in CI, and never refuses planque's start at deploy.
func TestTheShippedServerListDecodes(t *testing.T) {
	raw, err := os.ReadFile("../../deploy/server-apps.json")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := DecodeServerApps(raw)
	if err != nil {
		t.Fatalf("deploy/server-apps.json would refuse planque's start: %v", err)
	}
	if rows["https://mcp.stripe.com"] != StripeEntry || rows["https://mcp.qonto.com/mcp"] != QontoEntry || rows["https://mcp.agentmail.to"] != "connector:composio:agent_mail" {
		t.Fatalf("the shipped list's first rows: %v", rows)
	}
}

// A list that isn't exactly servers and entries is refused whole.
func TestDecodeServerApps(t *testing.T) {
	if _, err := DecodeServerApps([]byte(`{"https://mcp.example.com": "pay@1"}`)); err != nil {
		t.Fatalf("a reviewed entry's id: %v", err)
	}
	many := map[string]string{}
	for i := 0; i <= MaxServerApps; i++ {
		many["https://mcp"+strconv.Itoa(i)+".example.com"] = "connector:composio:app"
	}
	tooMany, _ := json.Marshal(many)
	for name, bad := range map[string]string{
		"plain http":    `{"http://mcp.example.com": "connector:composio:app"}`,
		"a query":       `{"https://mcp.example.com?x=1": "connector:composio:app"}`,
		"a fragment":    `{"https://mcp.example.com#a": "connector:composio:app"}`,
		"user info":     `{"https://u@mcp.example.com": "connector:composio:app"}`,
		"not an entry":  `{"https://mcp.example.com": "pay"}`,
		"half an entry": `{"https://mcp.example.com": "connector:composio:"}`,
		"not an object": `["https://mcp.example.com"]`,
		"trailing data": `{"https://mcp.example.com": "connector:composio:app"} {}`,
		"empty":         ``,
		"too many":      string(tooMany),
	} {
		if _, err := DecodeServerApps([]byte(bad)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

// A host that only looks like a listed one is never listed under its app:
// a subdomain, a suffix, a longer host, another port or plain http.
func TestTheServerListRefusesLookAlikes(t *testing.T) {
	s := &Server{ServerApps: map[string]string{"https://mcp.stripe.com": "connector:composio:stripe"}}
	if entry, listed := s.serverApp("https://mcp.stripe.com/mcp"); !listed || entry != "connector:composio:stripe" {
		t.Fatalf("the listed server: %q %v", entry, listed)
	}
	for _, server := range []string{
		"https://x.mcp.stripe.com",
		"https://x.mcp.stripe.com/mcp",
		"https://evilmcp.stripe.com",
		"https://evil-mcp.stripe.com/mcp",
		"https://mcp.stripe.com.evil.example",
		"https://mcp.stripe.com.evil.example/mcp",
		"https://mcp.stripe.com:8443",
		"https://mcp.stripe.com:8443/mcp",
		"http://mcp.stripe.com",
		"http://mcp.stripe.com/mcp",
	} {
		if entry, listed := s.serverApp(server); listed || entry != "" {
			t.Errorf("%s was listed under %q", server, entry)
		}
		if got := s.deriveEntry(&Grant{Request: Request{Server: server}}); got != "" {
			t.Errorf("%s got the entry %q", server, got)
		}
	}
}

// A server's connection is called by its reviewed app's name, never by the
// name its server gives itself (Stripe's says "stripe-mcp", which the owner
// read as the connection's name on 14 Nov). A server with no reviewed app
// keeps its own title, else its name.
func TestAServersConnectionIsCalledByItsReviewedApp(t *testing.T) {
	e := newEnv(t)
	stripe := mcp.Server{Name: "stripe-mcp"}
	for _, tc := range []struct {
		entry  string
		server mcp.Server
		want   string
	}{
		{StripeEntry, stripe, "Stripe"},
		{QontoEntry, mcp.Server{Name: "qonto-mcp-server"}, "Qonto"},
		{"", stripe, "stripe-mcp"},
		{"connector:none:unknown", mcp.Server{Name: "x", Title: "Example Desk"}, "Example Desk"},
	} {
		g := &Grant{Request: Request{CatalogEntry: tc.entry}}
		if got := e.s.serverVendorName(g, tc.server); got != tc.want {
			t.Errorf("entry %q, server %+v: %q, want %q", tc.entry, tc.server, got, tc.want)
		}
	}
	// A Stripe connection stored with its server's own name reads as Stripe,
	// its entry by planque's own id or by the managed app's id the reviewed
	// server list gives mcp.stripe.com (the owner's, 14 Nov).
	for _, entry := range []string{StripeEntry, "connector:composio:stripe", "connector:other:stripe"} {
		stored := &Grant{VendorName: "stripe-mcp", Request: Request{CatalogEntry: entry, Server: "https://mcp.stripe.com"}}
		if got := vendorName(stored); got != "Stripe" || connectionTitle(stored) != "Stripe" {
			t.Fatalf("a stored Stripe connection, entry %q: %q, %q", entry, got, connectionTitle(stored))
		}
	}
	if got := e.s.serverVendorName(&Grant{Request: Request{CatalogEntry: "connector:composio:qonto"}}, mcp.Server{Name: "qonto-mcp-server"}); got != "Qonto" {
		t.Fatalf("Qonto by its managed id: %q", got)
	}
	// Another managed app's id names nothing here: the stored name stays.
	other := &Grant{VendorName: "Gmail", Request: Request{CatalogEntry: "connector:composio:gmail"}}
	if got := vendorName(other); got != "Gmail" {
		t.Fatalf("another managed app: %q", got)
	}
}

// The owner's Stripe connection still read "stripe-mcp" after #210 and #212:
// its grant was stored before its request named its entry, so the request has
// none, and the entry recorded beside it (Grant.Entry) is the managed form the
// reviewed server list gave mcp.stripe.com then, planque's own, or none at
// all. Each reads "Stripe" in the grant's view, the account's connections and
// an ended taf's rows, with nothing rewritten. Another server whose own name
// ends in "-mcp" keeps it.
func TestAnOlderStoredStripeConnectionIsCalledStripe(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name, server, entry, stored, want string
	}{
		{"entry recorded as the managed app", "https://mcp.stripe.com", "connector:composio:stripe", "stripe-mcp", "Stripe"},
		{"entry recorded as planque's own", "https://mcp.stripe.com", StripeEntry, "stripe-mcp", "Stripe"},
		{"no entry, Stripe's server", "https://mcp.stripe.com", "", "stripe-mcp", "Stripe"},
		{"no entry, a path on Stripe's server", "https://mcp.stripe.com/v1/mcp", "", "stripe-mcp", "Stripe"},
		{"no entry, Qonto's server", "https://mcp.qonto.com/mcp", "", "qonto-mcp-server", "Qonto"},
		// NEGATIVE: another server's own "-mcp" name is kept, nothing stripped.
		{"another server named -mcp", "https://mcp.example.com/mcp", "", "example-mcp", "example-mcp"},
		{"a look-alike of Stripe's host", "https://mcp.stripe.com.example.net", "", "stripe-mcp", "stripe-mcp"},
		{"Qonto's host on another path", "https://mcp.qonto.com/other", "", "qonto-mcp-server", "qonto-mcp-server"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := newEnv(t)
			e.kiff.Keys["k_person"] = kiff.Me{TenantID: tenant}
			g := e.connect()
			g.Request.Server, g.Request.CatalogEntry, g.Entry, g.VendorName = tc.server, "", tc.entry, tc.stored
			if err := e.s.saveGrant(ctx, g); err != nil {
				t.Fatal(err)
			}
			title := tc.want + " · " + label(accountLabel(g)) // "Stripe · acct_…"
			if _, view := e.do("GET", "/v1/accounts/"+tenant+"/grants/"+g.ID, nil, nil, "k_setup"); view["vendor"] != tc.want || view["title"] != title {
				t.Errorf("the grant's view: vendor %v, title %v", view["vendor"], view["title"])
			}
			if _, cs := e.connections("k_person"); len(cs) != 1 || cs[0]["vendor"] != tc.want || cs[0]["title"] != title {
				t.Errorf("the connections: %v", cs)
			}
			if got := e.s.appName(g); got != tc.want {
				t.Errorf("an ended taf's rows: %q", got)
			}
			stored, err := e.s.loadGrant(ctx, tenant, g.ID)
			if err != nil || stored.VendorName != tc.stored || stored.Entry != tc.entry || stored.Request.CatalogEntry != "" {
				t.Errorf("the stored grant was rewritten: %+v, %v", stored, err)
			}
		})
	}
}
