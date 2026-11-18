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
	ClassUnplaced Class = "unplaced" // a write planque can't read the verb of: held
)

// HoldClasses switches the hold on. It's a variable so a test of something
// else can serve any fake tool; nothing else changes it.
var HoldClasses = true

var (
	deleteWords   = words("delete deletes remove removes destroy erase purge trash empty clear wipe drop revoke unlink terminate truncate")
	shareWords    = words("share unshare permission permissions invite invites acl collaborator collaborators publish unpublish public ownership owner role roles access member members")
	sendManyWords = words("campaign campaigns broadcast bulk mass newsletter blast audience audiences mailing everyone announce subscribers")
	moneyWords    = words("pay payment payments refund refunds transfer transfers payout payouts charge charges withdraw disburse paid void writeoff")
	readVerbs     = words("get list fetch search find read query retrieve check count describe lookup show view download export")
	writeVerbs    = words("create add update edit set append insert comment assign move copy rename reply post upload import duplicate label tag attach mark complete close reopen start stop pin unpin like follow save submit send draft log record schedule reschedule restore archive unarchive resolve batch upsert format sort resize mutate replace write put patch modify change enable disable activate deactivate merge split convert apply run execute trigger subscribe unsubscribe snooze approve reject accept decline")
)

func words(s string) map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(s) {
		m[w] = true
	}
	return m
}

// Classify is the class of t. The toolkit's own prefix on the slug
// (GOOGLEDRIVE_DELETE_CHILD) is dropped; the rest is split at underscores.
func Classify(t Tool) Class {
	slug := strings.ToLower(t.Slug)
	slug = strings.TrimPrefix(slug, strings.ToLower(t.Toolkit)+"_")
	tokens := strings.FieldsFunc(slug, func(r rune) bool { return r == '_' || r == '-' || r == ' ' })
	has := func(set map[string]bool) bool {
		for _, w := range tokens {
			if set[w] {
				return true
			}
		}
		return false
	}
	// A read by its verb and by the provider's mark may carry any noun: listing
	// permissions or a trash folder changes nothing.
	if t.ReadOnly && len(tokens) > 0 && readVerbs[tokens[0]] {
		return ClassRead
	}
	switch {
	case has(moneyWords):
		return ClassMoney
	case has(sendManyWords):
		return ClassSendMany
	case has(deleteWords):
		return ClassDelete
	case has(shareWords):
		return ClassShare
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
	case ClassDelete, ClassShare, ClassSendMany, ClassMoney, ClassUnplaced:
		return c, true
	}
	return "", false
}
