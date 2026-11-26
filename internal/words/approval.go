package words

// Approval is one language's words for what the owner reads, and signs, to
// approve a connection, and for how a connection stands. A field ending in
// F is a format: its values are named beside it, in that order, and a
// language may take them in another with %[n]s. The vendor's, the app's
// and the taf's own names, the tools' titles, the exact terms under Details
// and every id stay as they are.
type Approval struct {
	ThisTaf      string
	ThisTafFirst string // to open a sentence

	// A connection's first step: what an MCP server offers.
	DiscoverF       string // the server's host, the taf, the sign-in's host
	ConnectionVault string

	// The sign-in, and each tool.
	KeyGiven     string
	SignedInToF  string // the vendor account
	SignedIn     string
	WillSignInF  string // the sign-in's host
	MayF         string // the taf, opening a sentence
	Changes      string
	ReadsOnly    string
	AnyChange    string
	MovesMoneyF  string // the amount's argument
	SetupReadF   string // the read, the reads a day
	SetupReadsF  string // the reads, the reads a day
	And          string // joins the last two of a list
	KeyInHeaderF string // the header
	KeyInQueryF  string // the query parameter
	// KeyAsBasicPasswordF: the key is HTTP Basic's password, with the username given (%s).
	KeyAsBasicPasswordF string
	// KeyAsBasicUser: the key is HTTP Basic's username.
	KeyAsBasicUser string
	// FixedHeaderF follows where the key goes: the header every call also
	// carries (%s), and its value (%s).
	FixedHeaderF string

	// An API key.
	APIAsReadF   string // the description's URL, its digest, the API's base
	KeyVaultF    string // the API's host, where the key goes
	DiscoverAPIF string // the API's host, the taf, APIAsRead, KeyVault
	WillGiveKeyF string // the API's host, APIAsRead
	KeyTail      string

	// Tools added to a connection.
	FirstToolsF string // the tools
	AddedToolsF string // the tools

	// A managed app's connection. Who keeps the sign-in, plainly (our sign-in
	// partner: the owner signs these words), with no internal names, and no
	// provider's tool names: the actions are named, by their titles, on the
	// next card.
	ConnectTitledF   string // the app, the taf
	ConnectReusedF   string // the app, the taf
	ReusedNext       string
	EnterKeyF        string // the app, the taf
	EnterTokenF      string // the app, the taf
	SignInThenChoose string
	OnlyYourAccountF string // the app
	MoreAccessF      string // the app; says why (SignInThenChoose says the taf gets only what's approved)
	OnlyWhatItUsesF  string // the app

	// A registered application.
	AccessAtF     string // the application
	ReviewedFromF string // where its definitions are from

	// A browser.
	AnySiteF     string // the taf, opening a sentence
	TheseSitesF  string // the taf, opening a sentence, then the sites
	BrowserFills string
	BrowserPaysF string // the currency
	BrowserTail  string
	// The owner's own Chrome, through the tafs extension.
	ChromeSitesF string // the taf, opening a sentence, then the sites, then the Chrome's name
	ChromeTail   string

	// A Telegram bot.
	TelegramHookF     string // the taf
	TelegramVault     string
	TelegramConnectF  string // the taf, TelegramHook, TelegramVault
	TelegramWillGiveF string // TelegramHook
	TelegramTail      string

	// How a connection stands.
	PausedF       string // why
	PausedChanged string // a tool the person approved changed at the app since (the why of PausedF)
	PausedGone    string // a tool the person approved is no longer offered by the app (the why of PausedF)
	Connected     string
	// Connecting: approved and signed in, its tools still being registered in
	// KIFF (a few seconds). A language without it reads English's.
	Connecting       string
	NewToolsWaiting  string
	ToolsWaiting     string
	WaitingForYou    string
	DoorEvents       string
	ToolsBeingChosen string
	SignedOutByF     string // the vendor
	SignedOut        string // signed out, no vendor named
	ApprovedWaitF    string // the vendor
	WaitingToConfirm string
	NotConnected     string
	Disconnected     string
}

// Approve is lang's approval words, English's for a language planque
// doesn't speak.
func Approve(lang string) Approval {
	if a, ok := approvals[lang]; ok {
		return a
	}
	return approvals["en"]
}
