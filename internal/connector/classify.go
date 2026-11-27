package connector

import "strings"

// Class is what a managed app's tool does, read from its name and the
// provider's read-only mark. A tool in a held class (delete, share,
// send-to-many, money) or one planque can't place isn't served in wave 1a: it
// never comes back from a search or an app's tools, can't be asked for, and
// a call to one is never sent. They open in wave 1b, when a Card can hold
// them under "always asks you". The words are a floor, never a label the
// provider can lower: a tool named for deleting stays held whatever it says
// about itself.
type Class string

const (
	ClassRead     Class = "read"
	ClassWrite    Class = "write"    // creates or changes one thing: served
	ClassDelete   Class = "delete"   // deletes, empties, clears, revokes: held
	ClassShare    Class = "share"    // shares, publishes, changes who has access: held
	ClassSendMany Class = "sendmany" // a campaign, a broadcast, a bulk send: held
	ClassMoney    Class = "money"    // pays, refunds, transfers, charges: held
	ClassBulk     Class = "bulk"     // changes every item at once (archive all, move all): held
	ClassUnplaced Class = "unplaced" // a write planque can't read the verb of: held
)

// HoldClasses switches the hold on. It's a variable so a test of something
// else can serve any fake tool; nothing else changes it.
var HoldClasses = true

var (
	deleteWords   = words("delete deletes remove removes destroy erase purge trash empty clear wipe drop revoke unlink terminate truncate suspend ban kick")
	shareWords    = words("share shared sharing unshare permission permissions invite invites acl collaborator collaborators publish unpublish public ownership owner role roles access member members visibility grant grants")
	sendManyWords = words("campaign campaigns broadcast bulk mass newsletter blast audience audiences mailing everyone announce subscribers")
	moneyWords    = words("pay payment payments refund refunds transfer transfers payout payouts charge charges withdraw disburse paid void writeoff capture credit debit deposit purchase buy checkout")
	moneyVerbs    = words("pay refund transfer charge withdraw disburse void capture debit deposit purchase buy checkout")
	// dangerVerbs are words that are only ever verbs. Beside a read verb they
	// can only mean the tool also does that (GET_AND_DELETE, EXPORT_AND_PURGE);
	// refund, charge, purchase or checkout are mostly nouns there (LIST_CHARGE_REFUNDS).
	dangerVerbs = words("delete remove destroy erase purge wipe truncate pay withdraw disburse ban kick suspend")
	readVerbs   = words("get list fetch search find read query retrieve check count describe lookup show view download export")
	writeVerbs  = words("create add update edit set append insert comment assign move copy rename reply post upload import duplicate label tag attach mark complete close reopen start stop pin unpin like follow save submit send draft log record schedule reschedule restore archive unarchive resolve batch upsert format sort resize mutate replace write put patch modify change enable disable activate deactivate merge split convert apply run execute trigger subscribe unsubscribe snooze approve reject accept decline")
)

func words(s string) map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(s) {
		m[w] = true
	}
	return m
}

// Reads says t only reads: the provider marks it read-only and its name
// is a read (Classify). A provider that tags a write read-only
// (NOTION_UPDATE_PAGE, CONVERT_PURCHASE_ORDER_TO_BILL, APPLY_MACRO) gets a
// write everywhere planque says what a tool does: its catalog, the check on
// a read-only request, KIFF's declaration and tafs cloud's holds, which skip
// a read.
//
// A name that starts with a read verb is a read, as Classify has it, unless a
// write verb is joined on after it ("or", "and", "then"). One that starts with a
// noun (NOTION_PAGE_CREATE, SLACK_CHAT_POST_MESSAGE, HUBSPOT_CONTACTS_MERGE) is a
// write when any word in it is a write verb, and with no verb at all it follows
// the mark. A noun that is also a write verb after a read verb (GET_RECORD,
// RETRIEVE_POST_COMMENTS) never turns a read into a write.
func (t Tool) Reads() bool {
	if !t.ReadOnly || Classify(t) != ClassRead {
		return false
	}
	tokens := nameTokens(t)
	// A write verb joined on ("or", "and", "then") is a second act, wherever it comes:
	// GET_OR_CREATE_RECORD, SEARCH_AND_UPDATE_TICKETS, LIST_THEN_ARCHIVE.
	for i := 1; i < len(tokens); i++ {
		switch tokens[i-1] {
		case "and", "or", "then":
			if writeVerbs[tokens[i]] {
				return false
			}
		}
	}
	if len(tokens) == 0 || readVerbs[tokens[0]] || writeVerbs[tokens[0]] {
		// A write verb first is a read here only as "batch" before a read verb (BATCH_GET).
		return true
	}
	// A name that starts with a noun is a write when any of its words is a write
	// verb (PAGE_CREATE, PAGE_CREATE_AND_GET, SETTINGS_SEND_AS_GET): no read a
	// reviewed guide marks starts with a noun and carries one.
	for _, w := range tokens {
		if writeVerbs[w] {
			return false
		}
	}
	return true
}

// nameTokens is t's slug without the toolkit's own prefix, lower case, split
// at underscores, dashes and spaces.
func nameTokens(t Tool) []string {
	slug := strings.TrimPrefix(strings.ToLower(t.Slug), strings.ToLower(t.Toolkit)+"_")
	return strings.FieldsFunc(slug, func(r rune) bool { return r == '_' || r == '-' || r == ' ' })
}

// Classify is the class of t. The toolkit's own prefix on the slug
// (GOOGLEDRIVE_DELETE_CHILD) is dropped; the rest is split at underscores.
func Classify(t Tool) Class {
	tokens := nameTokens(t)
	has := func(set map[string]bool) bool {
		for _, w := range tokens {
			if set[w] {
				return true
			}
		}
		return false
	}
	// A read by its verb and by the provider's mark may carry a noun: listing
	// permissions, a trash, payments, refunds or a block changes nothing. It may
	// not carry a second, dangerous verb: the mark never lowers a name that
	// also deletes or moves money (GET_AND_DELETE_FILES, GET_AND_REFUND), whatever
	// the provider says. A second verb is one that follows "and"/"or"/"then", or
	// an unambiguous verb that ends the name.
	if t.ReadOnly && len(tokens) > 0 && readVerbs[tokens[0]] && !alsoActs(tokens) {
		return ClassRead
	}
	// "block" shuts someone out only as the verb that starts the name
	// (BLOCK_USER); inside a name it's mostly the noun (UPDATE_BLOCK).
	if len(tokens) > 0 && tokens[0] == "block" {
		return ClassDelete
	}
	switch {
	case has(moneyWords):
		return ClassMoney
	case has(sendManyWords), has(map[string]bool{"batch": true}) && has(map[string]bool{"send": true}):
		return ClassSendMany
	case has(deleteWords):
		return ClassDelete
	case has(shareWords):
		return ClassShare
	}
	// A write over every item at once (TRELLO_ARCHIVE_ALL_LIST_CARDS) is one
	// call that changes many things: held, whatever its verb.
	if !t.ReadOnly && has(map[string]bool{"all": true}) {
		return ClassBulk
	}
	// The provider's mark alone never makes a read of a name that starts with
	// a write verb (UPDATE_…, CONVERT_…_TO_BILL, APPLY_MACRO): a provider that
	// tags a write read-only would otherwise let it run as a read. "batch"
	// before a read verb (BATCH_GET) is a read all the same.
	// A write over every item (ARCHIVE_ALL_…) is held as bulk, as it is
	// unmarked; a read over all of them (LIST_ALL_…) never got here.
	if t.ReadOnly && len(tokens) > 0 && writeVerbs[tokens[0]] &&
		(tokens[0] != "batch" || len(tokens) < 2 || !readVerbs[tokens[1]] || alsoActs(tokens[1:])) {
		if has(map[string]bool{"all": true}) {
			return ClassBulk
		}
		return ClassWrite
	}
	if t.ReadOnly {
		return ClassRead
	}
	if has(writeVerbs) {
		return ClassWrite
	}
	return ClassUnplaced
}

// Held says t is in a class wave 1a doesn't serve, and which.
func Held(t Tool) (Class, bool) {
	if !HoldClasses {
		return "", false
	}
	switch c := Classify(t); c {
	case ClassDelete, ClassShare, ClassSendMany, ClassMoney, ClassBulk, ClassUnplaced:
		return c, true
	}
	return "", false
}

// alsoActs is whether a name that starts with a read verb goes on to delete or
// move money: a delete or money word after "and", "or" or "then", or a word that
// is only ever a verb at the end of the name.
func alsoActs(tokens []string) bool {
	for i := 1; i < len(tokens); i++ {
		switch tokens[i-1] {
		case "and", "or", "then":
			if deleteWords[tokens[i]] || moneyVerbs[tokens[i]] {
				return true
			}
		}
	}
	return len(tokens) > 1 && dangerVerbs[tokens[len(tokens)-1]]
}
