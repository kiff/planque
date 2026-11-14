package planque

import "net/http"

// toolRefusal is why planque refused a connection's tools, for the agent that
// fixes the request: a stable code, the tool and field at fault, the form
// planque wants there, and the words. The person never reads it: the agent
// turns it into plain words.
type toolRefusal struct {
	Msg   string `json:"error"`
	Code  string `json:"code"`
	Tool  string `json:"tool,omitempty"`
	Field string `json:"field,omitempty"`
	Want  string `json:"want,omitempty"`
}

func (r *toolRefusal) Error() string { return r.Msg }

// writeRefusal answers a refusal as 400 with its code and fields.
func writeRefusal(w http.ResponseWriter, r *toolRefusal) {
	writeJSON(w, http.StatusBadRequest, r)
}

// The codes a refusal of a connection's tools carries. They are stable: an
// agent acts on them, so a code is never renamed or reused for another fault.
// Field is the argument, pin or lookup at fault, when there is one.
const (
	// checkTools
	refAgentRequired       = "agent_required"         // no agent named
	refTooManyTools        = "too_many_tools"         // more than 20 tools
	refToolName            = "tool_name"              // a tool or its name isn't a valid name
	refToolTwice           = "tool_twice"             // the same tool name twice
	refToolTitle           = "tool_title"             // a title not on one line, or over 100 characters
	refTitleTwice          = "title_twice"            // two tools with the same title
	refReadOnlyMovesAmount = "read_only_moves_amount" // a read-only tool that moves an amount
	refIdempotencyReserved = "idempotency_reserved"   // amount_argument is planque's own idempotency argument
	refSetupReadsNotRead   = "setup_reads_not_read"   // setup_reads on a tool that isn't read-only

	// checkNarrowing, checkPin
	refPinShape         = "pin_shape"         // a pin's path isn't a dotted path, or its value isn't JSON; field = the pin
	refArgumentName     = "argument_name"     // an argument's name isn't valid
	refArgumentType     = "argument_type"     // an argument's type isn't one planque takes
	refArgumentSchema   = "argument_schema"   // an argument's reviewed schema is wrong, or a scalar has one
	refArgumentOneOf    = "argument_one_of"   // an argument's one_of is malformed
	refArgumentPath     = "argument_path"     // an argument's path isn't a dotted path
	refArgumentOverlap  = "argument_overlap"  // two arguments or pins at overlapping paths
	refAmountArgument   = "amount_argument"   // amount_argument isn't one of the tool's arguments of the right type
	refResourceArgument = "resource_argument" // resource_argument isn't valid, or not one of the tool's string arguments

	// checkLookups, checkOneOf
	refLookup = "lookup" // a unit_lookup or amount_lookup is incomplete or names no read; field = the lookup or its subfield
	refOneOf  = "one_of" // an argument's one_of names no read-only tool of the request

	// apiToolsCheck: the tools against the API's operations
	refNoOperation      = "no_operation"      // the API has no operation by the tool's name; field = the tool
	refNotARead         = "not_a_read"        // a read-only tool on an operation that isn't GET or HEAD
	refPinUnknown       = "pin_unknown"       // a pin isn't one of the operation's arguments; field = the pin
	refArgumentUnknown  = "argument_unknown"  // an argument's path isn't one of the operation's arguments
	refRequiredArgument = "required_argument" // a narrowed tool leaves out an argument the operation requires

	// notServed: a managed app's tool planque never serves (connector.NeverServed)
	refToolNotServed = "tool_not_served" // tool = its slug; field = "tools" (a catalog request) or "tool" (setting tools)
)

// notServed refuses a managed app's tool planque never serves, slug the
// tool and field where the request named it.
func notServed(slug, field string) *toolRefusal {
	return &toolRefusal{Code: refToolNotServed, Tool: slug, Field: field, Want: "another of the app's tools",
		Msg: slug + " is never served: what it does can't be undone (voiding, deleting, finalizing, paying or writing off a bill), so planque refuses it on this app; leave it out"}
}
