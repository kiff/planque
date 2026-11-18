package connector

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The hold ships on: TestMain turns it off for the older tests, so what it
// started as is read here, when the package's variables are set up.
var shippedHold = HoldClasses

func TestTheClassHoldShipsOn(t *testing.T) {
	if !shippedHold {
		t.Fatal("HoldClasses is off by default")
	}
}

// Refund, charge, purchase, checkout, deposit and block are mostly nouns in a
// tool's name. A read by its verb and the provider's mark stays a read with
// one; Notion's UPDATE_BLOCK is an edit, not a shutting-out.
func TestANounInAReadsNameDoesNotHoldIt(t *testing.T) {
	for _, c := range []struct {
		toolkit, slug string
		readOnly      bool
		want          Class
	}{
		{"stripe", "STRIPE_GET_REFUNDS_REFUND", true, ClassRead},
		{"stripe", "STRIPE_LIST_CHARGE_REFUNDS", true, ClassRead},
		{"stripe", "STRIPE_LIST_CHECKOUT_SESSIONS", true, ClassRead},
		{"shopify", "SHOPIFY_GET_ORDER_REFUND_BY_ID", true, ClassRead},
		{"quickbooks", "QUICKBOOKS_GET_PURCHASE", true, ClassRead},
		{"quickbooks", "QUICKBOOKS_GET_TRANSFER", true, ClassRead},
		{"moneybird", "MONEYBIRD_LIST_PURCHASE_INVOICES", true, ClassRead},
		{"notion", "NOTION_FETCH_BLOCK_CONTENTS", true, ClassRead},
		{"dropbox", "DROPBOX_CHECK_DELETE_BATCH", true, ClassRead},
		{"clay", "CLAY_GET_PRIVACY_DELETE_STATUS", true, ClassRead},
		{"notion", "NOTION_UPDATE_BLOCK", false, ClassWrite},
		// The verbs stay held in verb position.
		{"x", "X_GET_AND_REFUND", true, ClassMoney},
		{"x", "X_FETCH_THEN_DELETE", true, ClassDelete},
		{"x", "X_EXPORT_THEN_PURGE", true, ClassDelete},
		{"x", "X_LIST_OR_REMOVE_ITEMS", true, ClassDelete},
		{"x", "X_EXPORT_AND_PURGE", true, ClassDelete},
		{"x", "X_GET_FILES_AND_WIPE", true, ClassDelete},
		{"x", "X_CHECK_AND_PAY", true, ClassMoney},
		{"x", "X_GET_ITEM_DELETE", true, ClassDelete},
		{"github", "GITHUB_BLOCK_USER", false, ClassDelete},
		{"whatsapp", "WHATSAPP_BLOCK_USERS", false, ClassDelete},
		{"notion", "NOTION_DELETE_BLOCK", false, ClassDelete},
		// Words that the older cases only met beside another held word.
		{"x", "X_GRANT_LEVEL", false, ClassShare},
		{"x", "X_CAPTURE", false, ClassMoney},
		{"x", "X_KICK_MEMBER", false, ClassDelete},
		{"x", "X_KICK", false, ClassDelete},
		{"x", "X_BLOCK", false, ClassDelete},
	} {
		if got := Classify(tool(c.toolkit, c.slug, c.readOnly)); got != c.want {
			t.Errorf("%s (read-only %v): %s, want %s", c.slug, c.readOnly, got, c.want)
		}
	}
}

var toolRow = regexp.MustCompile("^\\| `([A-Z0-9_]+)` \\|.*\\| (read|write) \\|")

// Over the shipped lists of 72 apps (deploy/guides/generated): a tool the
// provider marks as a read, whose first word is a read verb, is a read unless
// its name goes on to delete or pay, and no tool marked write comes out as a
// read. Counted, so a list that stops being read fails here.
func TestTheShippedToolListsKeepTheirReads(t *testing.T) {
	files, _ := filepath.Glob("../../deploy/guides/generated/*.md")
	if len(files) < 60 {
		t.Fatalf("%d lists", len(files))
	}
	reads, held := 0, 0
	for _, f := range files {
		raw, _ := os.ReadFile(f)
		toolkit := strings.TrimSuffix(filepath.Base(f), ".md")
		for _, line := range strings.Split(string(raw), "\n") {
			m := toolRow.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			c := Classify(tool(toolkit, m[1], m[2] == "read"))
			if m[2] == "write" && c == ClassRead {
				t.Errorf("%s: a write came out a read", m[1])
			}
			if m[2] != "read" {
				continue
			}
			first := strings.Split(strings.ToLower(strings.TrimPrefix(m[1], strings.ToUpper(toolkit)+"_")), "_")
			if !readVerbs[first[0]] {
				continue
			}
			reads++
			if c != ClassRead {
				held++
				if !alsoActs(first) {
					t.Errorf("%s: a read held as %s", m[1], c)
				}
			}
		}
	}
	if reads < 3000 {
		t.Errorf("only %d reads found", reads)
	}
	t.Logf("%d reads by verb, %d held for a second verb", reads, held)
	if c := Classify(tool("notion", "NOTION_UPDATE_BLOCK", false)); c != ClassWrite {
		t.Errorf("NOTION_UPDATE_BLOCK is %s", c)
	}
}
