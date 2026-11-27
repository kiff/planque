package connector

import (
	"encoding/json"
	"testing"
)

func tool(toolkit, slug string, readOnly bool) Tool {
	return Tool{Slug: slug, Toolkit: toolkit, Version: "20261101_00", Name: slug, ReadOnly: readOnly, InputSchema: json.RawMessage(`{"type":"object"}`)}
}

// A tool's class comes from its name and the provider's read-only mark. The
// words are a floor: a tool named for deleting stays held when the provider
// marks it read-only; a read by its verb may name permissions or a trash;
// a write whose verb planque can't read is unplaced, so held.
func TestAToolsClassIsReadFromItsName(t *testing.T) {
	for _, c := range []struct {
		toolkit, slug string
		readOnly      bool
		want          Class
	}{
		{"googledrive", "GOOGLEDRIVE_GET_PERMISSION", true, ClassRead},
		{"googledrive", "GOOGLEDRIVE_LIST_TRASHED_FILES", true, ClassRead},
		{"googledrive", "GOOGLEDRIVE_DELETE_FOLDER_OR_FILE_ACTION", true, ClassDelete}, // the provider's mark doesn't lower it
		{"googledrive", "GOOGLEDRIVE_CREATE_PERMISSION", false, ClassShare},
		{"googledrive", "GOOGLEDRIVE_EMPTY_TRASH", false, ClassDelete},
		{"googlesheets", "GOOGLESHEETS_CLEAR_VALUES", false, ClassDelete},
		{"googlesheets", "GOOGLESHEETS_SPREADSHEETS_VALUES_APPEND", false, ClassWrite},
		{"googlesheets", "GOOGLESHEETS_UPSERT_ROWS", false, ClassWrite},
		{"mailchimp", "MAILCHIMP_SEND_CAMPAIGN", false, ClassSendMany},
		{"gmail", "GMAIL_SEND_EMAIL", false, ClassWrite},
		{"gmail", "GMAIL_FETCH_EMAILS", true, ClassRead},
		{"stripe", "STRIPE_PAY_INVOICE", false, ClassMoney},
		{"stripe", "STRIPE_CREATE_REFUND", false, ClassMoney},
		{"notion", "NOTION_ADD_PAGE_CONTENT", false, ClassWrite},
		{"notion", "NOTION_DELETE_BLOCK", false, ClassDelete},
		{"asana", "ASANA_ADD_MEMBER_TO_PROJECT", false, ClassShare},
		{"x", "X_FROBNICATE_THING", false, ClassUnplaced},
		{"x", "X_FROBNICATE_THING", true, ClassRead}, // an unknown verb the provider marks as a read
		{"x", "X_DELETED_ITEMS_REPORT", true, ClassRead},
		// The provider's mark never lowers a name that also deletes or moves money.
		{"x", "X_GET_AND_DELETE_FILES", true, ClassDelete},
		{"x", "X_EXPORT_AND_PURGE", true, ClassDelete},
		{"x", "X_GET_AND_REFUND", true, ClassMoney},
		{"x", "X_LIST_PAYMENTS", true, ClassRead}, // a noun alone is still a read
		// Sharing by other words.
		{"x", "X_UPDATE_SHARING_SETTINGS", false, ClassShare},
		{"x", "X_SET_VISIBILITY", false, ClassShare},
		{"x", "X_GRANT_ACCESS_LEVEL", false, ClassShare},
		// A send in a batch is a send to many; more money words; people shut out.
		{"gmail", "GMAIL_BATCH_SEND", false, ClassSendMany},
		{"bexio", "BEXIO_CREATE_CREDIT_NOTE", false, ClassMoney},
		{"shop", "SHOP_CREATE_CHECKOUT_SESSION", false, ClassMoney},
		{"x", "X_CAPTURE_PAYMENT", false, ClassMoney},
		{"x", "X_SUSPEND_USER", false, ClassDelete},
		{"x", "X_BAN_USER", false, ClassDelete},
		// A shared link is a share (the live catalog, 19 Nov); reading one stays a read.
		{"dropbox", "DROPBOX_CREATE_SHARED_LINK", false, ClassShare},
		{"dropbox", "DROPBOX_MODIFY_SHARED_LINK_SETTINGS", false, ClassShare},
		{"dropbox", "DROPBOX_LIST_SHARED_LINKS", true, ClassRead},
		{"dropbox", "DROPBOX_GET_SHARED_LINK_METADATA", true, ClassRead},
		// A write over every item at once is bulk; one item, or a read of all, isn't.
		{"trello", "TRELLO_ARCHIVE_ALL_LIST_CARDS", false, ClassBulk},
		{"trello", "TRELLO_MOVE_ALL_LIST_CARDS", false, ClassBulk},
		{"trello", "TRELLO_MARK_ALL_NOTIFICATIONS_READ", false, ClassBulk},
		{"linear", "LINEAR_ARCHIVE_ISSUE", false, ClassWrite},
		{"jira", "JIRA_GET_ALL_GROUPS", true, ClassRead},
		{"cal", "CAL_FETCH_ALL_BOOKINGS", true, ClassRead},
	} {
		if got := Classify(tool(c.toolkit, c.slug, c.readOnly)); got != c.want {
			t.Errorf("%s (read-only %v): %s, want %s", c.slug, c.readOnly, got, c.want)
		}
	}
}

// Only the four held classes and the unplaced are held; reads and ordinary
// writes are served; and the switch off serves all.
func TestOnlyHeldClassesAreHeld(t *testing.T) {
	old := HoldClasses
	defer func() { HoldClasses = old }()
	HoldClasses = true
	for slug, held := range map[string]bool{"N_DELETE_X": true, "N_SHARE_X": true, "N_SEND_CAMPAIGN": true, "N_REFUND_X": true, "N_ARCHIVE_ALL_X": true, "N_FROBNICATE_X": true, "N_CREATE_X": false, "N_LIST_X": false} {
		_, got := Held(tool("n", slug, slug == "N_LIST_X"))
		if got != held {
			t.Errorf("%s held %v, want %v", slug, got, held)
		}
	}
	HoldClasses = false
	if _, got := Held(tool("n", "N_DELETE_X", false)); got {
		t.Error("held with the switch off")
	}
}

// A held tool isn't a definition: it can't be granted, listed or run, the way
// a never-served tool isn't, and a served one validates as before.
func TestAHeldToolIsNotADefinition(t *testing.T) {
	old := HoldClasses
	defer func() { HoldClasses = old }()
	HoldClasses = true
	if err := tool("notion", "NOTION_DELETE_BLOCK", false).Validate(); err != ErrDefinition {
		t.Errorf("a delete validated: %v", err)
	}
	if err := tool("notion", "NOTION_ADD_PAGE_CONTENT", false).Validate(); err != nil {
		t.Errorf("a write refused: %v", err)
	}
}

// The provider's read-only mark doesn't make a read of a name that starts with
// a write verb: a provider that tags a write read-only (seen in the generated
// guides: CONVERT_PURCHASE_ORDER_TO_BILL, APPLY_ZENDESK_MACRO) would let it
// run unasked. A batch read stays a read, and a name with no write verb first
// still follows the mark.
func TestTheReadOnlyMarkDoesntMakeAWriteVerbARead(t *testing.T) {
	for _, c := range []struct {
		toolkit, slug string
		want          Class
	}{
		{"notion", "NOTION_UPDATE_PAGE", ClassWrite},
		{"zoho_books", "ZOHO_BOOKS_CONVERT_PURCHASE_ORDER_TO_BILL", ClassMoney},
		{"zendesk", "ZENDESK_APPLY_ZENDESK_MACRO", ClassWrite},
		{"intercom", "INTERCOM_CREATE_DATA_EXPORT", ClassWrite},
		{"googlesheets", "GOOGLESHEETS_BATCH_GET", ClassRead},
		{"googlesheets", "GOOGLESHEETS_SPREADSHEETS_VALUES_GET", ClassRead},
		{"gmail", "GMAIL_LIST_MESSAGES", ClassRead},
		{"gmail", "GMAIL_SETTINGS_SEND_AS_GET", ClassRead},
		{"desk", "DESK_BATCH_GET_AND_DELETE_TICKETS", ClassDelete},
		{"desk", "DESK_BATCH_UPDATE_TICKETS", ClassWrite},
	} {
		if got := Classify(tool(c.toolkit, c.slug, true)); got != c.want {
			t.Errorf("%s marked read-only: %s, want %s", c.slug, got, c.want)
		}
	}
}
