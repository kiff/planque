package planque

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kiff/planque/internal/device"
	"github.com/kiff/planque/internal/fakes"
	"github.com/kiff/planque/internal/kiff"
	"github.com/kiff/planque/internal/openapi"
	"github.com/kiff/planque/internal/profile"
)

func refundProfile() profile.Profile {
	return profile.Profile{V: profile.Version, Vendor: "Pay API", Transport: "rest", Pin: profile.Pin{Kind: "base_path", Value: "/v1"}, Source: "drafted",
		Operations: []profile.Operation{{Tool: "refund", Title: "Refund a charge", Kind: "write", Effect: profile.EffectMovesMoney, Shape: "single",
			Resource: &profile.Resource{Arg: "charge", Noun: "charge"}, Amount: &profile.Amount{From: "argument", Arg: "amount", Format: "decimal_string", Scale: 2, Unit: "EUR"},
			RequiredArguments: []string{"amount"}, Marker: &profile.Marker{Field: "metadata.kiff_op"}}}}
}

func profileResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}

func TestProfileApprovalLifecycle(t *testing.T) {
	for _, flow := range []string{"key", "approve"} {
		for _, scenario := range []string{"accepted", "mismatch", "unavailable", "wrong_receipt", "accepted_header"} {
			if flow == "approve" && scenario == "accepted_header" {
				continue
			}
			t.Run(flow+"/"+scenario, func(t *testing.T) {
				outcome := scenario
				headerPin := scenario == "accepted_header"
				if headerPin {
					outcome = "accepted"
				}
				e := newInProcessEnv(t)
				e.enrol()
				ctx := context.Background()
				var vendorCalls atomic.Int32
				var sent map[string]any
				description := strings.ReplaceAll(restDesc, `"amount": {"type": "integer"}`, `"amount": {"type": "string"}, "metadata": {"type":"object","additionalProperties":{"type":"string"}}`)
				e.s.REST.HTTP = &http.Client{Transport: applicationTransport(func(r *http.Request) (*http.Response, error) {
					if r.URL.Path == "/openapi.json" {
						return profileResponse(200, description), nil
					}
					if r.Header.Get("X-API-Key") != "test-vendor-key" {
						t.Error("vendor did not receive its stored key")
					}
					if headerPin && r.Header.Get("Pay-Version") != "2030-01-01" {
						t.Error("the approved version header was not sent")
					}
					vendorCalls.Add(1)
					if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
						t.Error(err)
					}
					return profileResponse(200, `{"id":"refund_one"}`), nil
				})}
				tools := []ToolRequest{{Tool: "createRefund", Name: "refund", Title: "Refund a charge"}}
				p := refundProfile()
				if headerPin {
					p.Pin = profile.Pin{Kind: "header", Value: "2030-01-01"}
					p.Source = "catalog"
					p.CatalogEntry = "pay@1"
				}
				body := map[string]any{"kind": KindAPIKey, "agent": "agt_taf", "tools": tools, "profile": p,
					"api": map[string]any{"url": "https://pay.example/openapi.json", "base": "https://pay.example/v1", "auth": map[string]string{"in": "header", "name": "X-API-Key"}}}
				if flow == "approve" {
					body["tools"] = []ToolRequest{}
					delete(body, "profile")
				}
				path := "/v1/accounts/" + tenant + "/grants"
				createPath := path
				if headerPin {
					raw, _ := json.Marshal(p)
					sum := sha256.Sum256([]byte(description))
					entry := CatalogEntry{ID: "pay@1", Vendor: p.Vendor, Title: "Refund payments", Kind: KindAPIKey, VersionHeader: "Pay-Version", Profile: raw, Tools: tools,
						API: &apiAsk{URL: "https://pay.example/openapi.json", Base: "https://pay.example/v1", Auth: openapi.Auth{In: "header", Name: "X-API-Key"}}, DescriptionSHA256: hex.EncodeToString(sum[:])}
					conf, _ := json.Marshal(map[string]CatalogEntry{entry.ID: entry})
					var err error
					e.s.CatalogEntries, err = DecodeCatalog(conf, nil)
					if err != nil {
						t.Fatal(err)
					}
					body = map[string]any{"agent": "agt_taf"}
					createPath = "/v1/accounts/" + tenant + "/catalog/pay@1/grants"
				}
				status, requested := e.do("POST", createPath, body, nil, "k_setup")
				if headerPin && status == 201 {
					requested = requested["grant"].(map[string]any)
				}
				if status != 201 {
					t.Fatalf("request: %d %v", status, requested)
				}
				id := requested["id"].(string)
				keyRequest := func(approved *ProfileDigests, digest string) (int, map[string]any) {
					sum := sha256.Sum256([]byte("test-vendor-key"))
					payload := map[string]any{"digest": digest, "key_sha256": hex.EncodeToString(sum[:])}
					if approved != nil {
						payload["kiff_profile"] = approved
					}
					b, _ := json.Marshal(payload)
					r, _ := http.NewRequest("POST", e.url+path+"/"+id+"/key", bytes.NewReader(b))
					e.phone.Sign(r, tenant, b)
					r.Header.Set(KeyHeader, "test-vendor-key")
					res, err := e.client.Do(r)
					if err != nil {
						t.Fatal(err)
					}
					defer res.Body.Close()
					var out map[string]any
					_ = json.NewDecoder(res.Body).Decode(&out)
					return res.StatusCode, out
				}
				if flow == "approve" {
					status, out := keyRequest(nil, requested["digest"].(string))
					if status != 200 || out["status"] != StatusListed {
						t.Fatalf("list: %d %v", status, out)
					}
					status, requested = e.do("POST", path+"/"+id+"/tools", map[string]any{"tools": tools, "profile": p}, nil, "k_setup")
					if status != 200 {
						t.Fatalf("set tools: %d %v", status, requested)
					}
				}
				g, err := e.s.loadGrant(ctx, tenant, id)
				if err != nil {
					t.Fatal(err)
				}
				if headerPin && (g.Request.API.VersionPin == nil || g.Request.API.VersionPin.Name != "Pay-Version" || g.Request.API.VersionPin.Value != p.Pin.Value) {
					t.Fatal("the signed grant lost its version pin")
				}
				profileHash, _ := profile.Digest(g.Request.Profile)
				approved := &ProfileDigests{ProfileSHA256: profileHash, WordsSHA256: strings.Repeat("a", 64)}
				var applications atomic.Int32
				base := e.client.Transport
				e.client.Transport = applicationTransport(func(r *http.Request) (*http.Response, error) {
					if r.URL.Host != "kiff.example" || r.URL.Path != "/v1/profiles" {
						return base.RoundTrip(r)
					}
					applications.Add(1)
					var req struct {
						Profile  json.RawMessage      `json:"profile"`
						Approval kiff.ProfileApproval `json:"approval"`
					}
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						t.Fatal(err)
					}
					message, _ := base64.StdEncoding.DecodeString(req.Approval.Message)
					signedBody, _ := base64.StdEncoding.DecodeString(req.Approval.Body)
					lines := strings.Split(string(message), "\n")
					if len(lines) != 7 || !strings.HasSuffix(lines[4], "/"+flow) || !e.phone.Verifies(message, req.Approval.Signature) || req.Approval.DeviceID != e.phone.DeviceID {
						t.Fatal("profile application changed the signed approval evidence")
					}
					if !bytes.Equal(message, device.Message(e.url, tenant, "POST", path+"/"+id+"/"+flow, signedBody, lines[6])) {
						t.Fatal("profile evidence body was not the signed request")
					}
					gotHash, _ := profile.Digest(req.Profile)
					if gotHash != profileHash || r.Header.Get("Authorization") != "Bearer k_connector" {
						t.Fatal("wrong profile or applying credential")
					}
					conn, ok := e.kiff.Connection(tenant, "refund")
					if !ok || conn.AmountFormat != "decimal_string" || conn.AmountArg != "amount" || conn.AmountScale == nil || *conn.AmountScale != 2 || conn.AmountUnit != "EUR" || conn.ResourceArg != "charge" {
						t.Fatal("profile was not applied after registering its measurement settings")
					}
					pending, _ := e.s.loadGrant(ctx, tenant, id)
					args := json.RawMessage(`{"amount":"80.50","charge":"charge_one","planque_op":"during-application"}`)
					permit := e.kiff.Mint(fakes.Permit{Audience: pending.Audience, Tenant: tenant, Tool: "refund", Sub: "agt_taf", Op: "during-application", Args: args})
					result := e.s.call(ctx, pending, "refund", args, permit)
					if !strings.Contains(mustJSONString(result), "not been applied") || vendorCalls.Load() != 0 {
						t.Fatal("relay forwarded before KIFF applied the profile")
					}
					if outcome == "mismatch" {
						return profileResponse(409, `{"error":"profile_mismatch"}`), nil
					}
					if outcome == "unavailable" {
						return profileResponse(503, `{"error":"unavailable"}`), nil
					}
					if outcome == "wrong_receipt" {
						gotHash = strings.Repeat("b", 64)
					}
					out, _ := json.Marshal(kiff.ProfileRevision{Revision: 3, ProfileSHA256: gotHash, WordsSHA256: approved.WordsSHA256})
					return profileResponse(200, string(out)), nil
				})
				approve := func(d *ProfileDigests) (int, map[string]any) {
					if flow == "key" {
						return keyRequest(d, requested["digest"].(string))
					}
					payload := map[string]any{"digest": requested["digest"]}
					if d != nil {
						payload["kiff_profile"] = d
					}
					return e.do("POST", path+"/"+id+"/approve", payload, e.phone, "")
				}
				if status, _ := approve(nil); status != 409 || applications.Load() != 0 {
					t.Fatal("profile approved without its signed digests")
				}
				bad := *approved
				bad.ProfileSHA256 = strings.Repeat("b", 64)
				if status, _ := approve(&bad); status != 409 || applications.Load() != 0 {
					t.Fatal("another profile's digest was accepted")
				}
				status, out := approve(approved)
				if applications.Load() != 1 {
					t.Fatalf("applications=%d", applications.Load())
				}
				g, _ = e.s.loadGrant(ctx, tenant, id)
				if outcome != "accepted" {
					if status != 502 || g.Status == StatusConnected || len(e.kiff.Connections(tenant)) != 0 || vendorCalls.Load() != 0 {
						t.Fatalf("failed profile left authority active: %d %v", status, out)
					}
					return
				}
				if status != 200 || g.Status != StatusConnected || !profileReady(g) {
					t.Fatalf("connect: %d %v", status, out)
				}
				conn, _ := e.kiff.Connection(tenant, "refund")
				forward := func(args json.RawMessage, operation string) json.RawMessage {
					permit := e.kiff.Mint(fakes.Permit{Audience: g.Audience, Tenant: tenant, Tool: "refund", Sub: "agt_taf", Op: operation, Args: args})
					result, err := e.kiff.Forward(ctx, conn, args, permit)
					if err != nil {
						t.Fatal(err)
					}
					return result
				}
				args := json.RawMessage(`{"amount":"80.50","charge":"charge_one","planque_op":"op_one"}`)
				forward(args, "op_one")
				first := forward(args, "op_one")
				// A later profile cannot reinterpret or re-send an earlier call.
				original := g.Request.Profile
				g.Request.Profile = json.RawMessage(strings.Replace(string(original), "Pay API", "Pay API changed", 1))
				if err := e.s.saveGrant(ctx, g); err != nil {
					t.Fatal(err)
				}
				if replay := forward(args, "op_one"); !bytes.Equal(first, replay) {
					t.Fatal("a later profile changed an earlier operation's replay")
				}
				g.Request.Profile = original
				if err := e.s.saveGrant(ctx, g); err != nil {
					t.Fatal(err)
				}
				metadata, _ := sent["metadata"].(map[string]any)
				if vendorCalls.Load() != 1 || sent["amount"] != "80.50" || metadata["kiff_op"] != "op_one" {
					t.Fatalf("decimal, marker or replay: %v calls=%d", sent, vendorCalls.Load())
				}
				result := forward(json.RawMessage(`{"charge":"charge_one","planque_op":"missing_amount"}`), "missing_amount")
				if !strings.Contains(string(result), "profile requires argument amount") || vendorCalls.Load() != 1 {
					t.Fatal("missing profile argument was sent")
				}
				result = forward(json.RawMessage(`{"amount":"1.00","charge":"charge_one","metadata":{"kiff_op":"caller-marker"},"planque_op":"spoof"}`), "spoof")
				if !strings.Contains(string(result), "not a caller") || vendorCalls.Load() != 1 {
					t.Fatal("caller supplied a trusted operation marker")
				}
				var record operation
				_, err = getDoc(ctx, e.kv, tenantPK(tenant), skOp(id, "op_one"), &record)
				if err != nil || record.ProfileRevision == nil || record.ProfileRevision.Revision != 3 || !bytes.Equal(record.Profile, g.Request.Profile) {
					t.Fatal("operation did not retain its approved profile revision")
				}
			})
		}
	}
}

func TestDraftCannotAssertCatalogOrTestEvidence(t *testing.T) {
	for _, field := range []string{"catalog", "test"} {
		p := refundProfile()
		if field == "catalog" {
			p.Source = "catalog"
			p.CatalogEntry = "unreviewed@1"
		} else {
			p.TestProven = true
		}
		raw, _ := json.Marshal(p)
		if _, err := profileFields(raw, []ToolRequest{{Name: "refund", Tool: "createRefund"}}); err == nil {
			t.Fatal("proposal asserted trusted evidence")
		}
	}
}
