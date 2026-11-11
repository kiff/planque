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

	// An API key.
	APIAsReadF   string // the description's URL, its digest, the API's base
	KeyVaultF    string // the API's host, where the key goes
	DiscoverAPIF string // the API's host, the taf, APIAsRead, KeyVault
	WillGiveKeyF string // the API's host, APIAsRead
	KeyTail      string

	// Tools added to a connection.
	FirstToolsF string // the tools
	AddedToolsF string // the tools

	// A connection through a connection service.
	ConnectTitledF   string // the app, the taf
	ConnectReusedF   string // the app, the taf
	ReusedNext       string
	EnterKeyF        string // the app, the taf
	EnterTokenF      string // the app, the taf
	SignInThenChoose string
	OnlyYourAccountF string // the app
	MoreAccessF      string // the app, its tools
	OnlyWhatItUsesF  string // the app, its tools

	// A registered application.
	AccessAtF     string // the application
	ReviewedFromF string // where its definitions are from

	// A browser.
	AnySiteF     string // the taf, opening a sentence
	TheseSitesF  string // the taf, opening a sentence, then the sites
	BrowserFills string
	BrowserPaysF string // the currency
	BrowserTail  string

	// A Telegram bot.
	TelegramHookF     string // the taf
	TelegramVault     string
	TelegramConnectF  string // the taf, TelegramHook, TelegramVault
	TelegramWillGiveF string // TelegramHook
	TelegramTail      string

	// How a connection stands.
	PausedF          string // why
	Connected        string
	NewToolsWaiting  string
	ToolsWaiting     string
	WaitingForYou    string
	DoorEvents       string
	ToolsBeingChosen string
	SignedOutByF     string // the vendor
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
