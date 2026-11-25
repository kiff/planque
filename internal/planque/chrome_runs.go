package planque

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/kiff/planque/internal/store"
)

// The runs that worked in a Chrome, for its side panel (kiff/tafs
// docs/chrome.md, "The side panel's chats"): the panel shows a run that
// worked there, wherever its message was sent from, and only that run.
//
// Each step planque sends a Chrome marks that Chrome and the kiff/agents run
// whose call sent it: the run is the call's permit's, which KIFF signed, and
// the Chrome is the one the step went to, so a connection moved to another
// Chrome later carries nothing with it. The panel asks with its own KIFF
// session: KIFF's /v1/me names the session's device, and planque finds the
// paired Chrome that device is by the key the extension offered for it at
// pairing (Chrome.KiffPub), which KIFF matches to its device
// (POST /v1/devices/by-key). Nothing the extension or a request says picks
// the Chrome. Unpairing drops the marks with the Chrome; a Chrome paired
// before KiffPub was kept has no device, and lists nothing until it pairs
// again.
//
// At 10,000 users: one mark per run per Chrome, written on the run's first
// step only, and one read per panel list of the account's Chromes (at most a
// few) and that Chrome's marks, cut to chromeWorkedListed. KIFF is asked
// about a Chrome only until it names its device, and at most once a
// chromeKiffAskEvery while it names none.

// chromeWorkedKept is how long a mark is kept: as long as tafs keeps a
// Chrome's mark on a chat it began, so a run's chat leaves no sooner.
const chromeWorkedKept = 400 * 24 * time.Hour

// chromeWorkedListed is the most runs a panel's read lists, newest first.
const chromeWorkedListed = 50

// chromeKiffAskEvery is how often KIFF is asked again for a Chrome's device
// while it names none (the extension claims its KIFF pairing after the
// phone's confirm, so the first ask may come before the device exists).
const chromeKiffAskEvery = time.Minute

func skChromeWorked(chrome, run string) string { return "CHROMEWORKED#" + chrome + "#" + run }

// chromeWorked is one run that worked in a Chrome.
type chromeWorked struct {
	Run string    `json:"run"`
	Taf string    `json:"taf"`
	At  time.Time `json:"at"`
}

// kiffOfferPub is the key in the extension's KIFF offer (kiff.pair.v1, as
// JSON text), or "" for an offer that isn't the extension's.
func kiffOfferPub(offer string) string {
	var o struct {
		V        string `json:"v"`
		Pub      string `json:"pub"`
		Platform string `json:"platform"`
	}
	if offer == "" || json.Unmarshal([]byte(offer), &o) != nil || o.V != "kiff.pair.v1" || o.Platform != "chrome" {
		return ""
	}
	return o.Pub
}

// markChromeWorked keeps that run worked in chrome. A call with no run (a
// permit from before runs were named) marks nothing. A mark that can't be
// kept is logged: the step still goes.
func (s *Server) markChromeWorked(ctx context.Context, tenant, chrome, run, taf string) {
	if run == "" {
		return
	}
	now := s.now().UTC()
	err := putDoc(ctx, s.KV, store.Item{PK: tenantPK(tenant), SK: skChromeWorked(chrome, run), TTL: now.Add(chromeWorkedKept).Unix()},
		chromeWorked{Run: run, Taf: taf, At: now}, store.Cond{})
	if err != nil {
		s.log().ErrorContext(ctx, "a run's mark in its Chrome wasn't kept", slog.String("tenant", tenant), slog.String("chrome", chrome), slog.String("err", err.Error()))
	}
}

// dropChromeWorked removes an unpaired Chrome's marks. One left is never
// read: the Chrome is gone, and its mark ends on its own.
func (s *Server) dropChromeWorked(ctx context.Context, tenant, chrome string) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), skChromeWorked(chrome, ""))
	if err != nil {
		s.log().ErrorContext(ctx, "an unpaired Chrome's runs weren't listed", slog.String("tenant", tenant), slog.String("chrome", chrome), slog.String("err", err.Error()))
		return
	}
	for _, it := range items {
		_, _, _ = s.KV.Delete(ctx, it.PK, it.SK)
	}
}

// chromeOwnRuns serves GET /v1/accounts/{tenant}/chromes/self/runs: the runs
// that worked in the asking Chrome, newest first, as {runs: [{run, taf, at}]}.
// Only the tafs extension's own KIFF session asks: any other caller, a
// Chrome not paired here, or one whose device KIFF hasn't named, has none.
func (s *Server) chromeOwnRuns(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant")
	me, err := s.agentKey(r, tenant)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	runs := []chromeWorked{}
	if !me.DeviceSession || me.DevicePlatform != "chrome" || me.DeviceID == "" || !s.chromeOnFor(tenant) {
		writeJSON(w, http.StatusOK, map[string]any{"runs": runs})
		return
	}
	chrome, err := s.chromeOfDevice(r.Context(), tenant, me.DeviceID)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	if chrome != "" {
		items, err := s.KV.Query(r.Context(), tenantPK(tenant), skChromeWorked(chrome, ""))
		if err != nil {
			s.fail(w, r, err)
			return
		}
		now := s.now()
		for _, it := range items {
			var m chromeWorked
			if json.Unmarshal(it.Doc, &m) != nil || m.Run == "" || (it.TTL != 0 && now.Unix() >= it.TTL) {
				continue
			}
			runs = append(runs, m)
		}
		sort.SliceStable(runs, func(i, j int) bool { return runs[i].At.After(runs[j].At) })
		runs = runs[:min(len(runs), chromeWorkedListed)]
	}
	writeJSON(w, http.StatusOK, map[string]any{"runs": runs})
}

// chromeOfDevice is the account's paired Chrome whose KIFF device is device,
// or "" when none is. A Chrome whose device KIFF hasn't named yet is asked
// about by its kept key, with the account's connector key, and keeps the
// device KIFF names.
func (s *Server) chromeOfDevice(ctx context.Context, tenant, device string) (string, error) {
	items, err := s.KV.Query(ctx, tenantPK(tenant), "CHROME#")
	if err != nil {
		return "", err
	}
	for _, it := range items {
		var c Chrome
		if json.Unmarshal(it.Doc, &c) != nil || c.ID == "" || skChrome(c.ID) != it.SK || c.Tenant != tenant {
			continue
		}
		if c.KiffDevice != "" {
			if c.KiffDevice == device {
				return c.ID, nil
			}
			continue
		}
		if c.KiffPub == "" || c.KiffAsked != nil && s.now().Sub(*c.KiffAsked) < chromeKiffAskEvery {
			continue
		}
		named, err := s.kiffDeviceOf(ctx, tenant, c.KiffPub)
		if err != nil {
			s.log().WarnContext(ctx, "KIFF wasn't asked for a Chrome's device", slog.String("tenant", tenant), slog.String("chrome", c.ID), slog.String("err", err.Error()))
			continue
		}
		if named == "" {
			at := s.now().UTC()
			c.KiffAsked = &at
		} else {
			c.KiffDevice, c.KiffAsked = named, nil
		}
		ver := it.Ver
		// Kept only over the Chrome as it was read: an unpairing or another
		// write in between wins, and the next read asks again.
		_ = putDoc(ctx, s.KV, store.Item{PK: it.PK, SK: it.SK}, c, store.Cond{Ver: &ver})
		if named == device {
			return c.ID, nil
		}
	}
	return "", nil
}

// kiffDeviceOf asks KIFF which of the account's devices holds pub.
func (s *Server) kiffDeviceOf(ctx context.Context, tenant, pub string) (string, error) {
	key, err := s.connectorKey(ctx, tenant)
	if err != nil {
		return "", err
	}
	return s.KIFF.ChromeDevice(ctx, key, strings.TrimSpace(pub))
}
