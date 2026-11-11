package words

// approvals is every language's approval words. Each language has every
// field, taking English's values (approval_test.go).
var approvals = map[string]Approval{
	"en": {
		ThisTaf: "this taf", ThisTafFirst: "This taf",

		DiscoverF: "Connect your account at %s to see what it offers for %s.\nYou will sign in at %s.\n\n" +
			"Only what it offers is read: its list of tools. None of them is called: nothing is read from your account, nothing is changed, and nothing is given to the taf. " +
			"The tools the taf may use come next, and you approve them again. If they aren't approved within a day, the sign-in is revoked.\n\n" +
			"Your sign-in is kept in a sealed vault; the taf never sees it. You can disconnect at any time.",
		ConnectionVault: "\nYour sign-in is kept in a sealed vault; the taf never sees it. Each call is sent only with KIFF's permit for that exact call. You can disconnect at any time.",

		KeyGiven:     "The key you gave is used: no second key.\n\n",
		SignedInToF:  "You're already signed in, to the account %s: no second sign-in.\n\n",
		SignedIn:     "You're already signed in: no second sign-in.\n\n",
		WillSignInF:  "You will sign in at %s.\n\n",
		MayF:         "%s may:\n",
		Changes:      "makes changes, only within the taf's Card",
		ReadsOnly:    "reads only",
		AnyChange:    "can make any change this tool allows; the taf's Card can allow or refuse each call, not limit what it changes",
		MovesMoneyF:  "moves money (%s), only within the taf's Card",
		SetupReadF:   "\nSetup reads, to set this up: %s. It only reads it, through KIFF with its own Card, at most %d reads a day; each read has KIFF's permit and is kept in KIFF's record.\n",
		SetupReadsF:  "\nSetup reads, to set this up: %s. It only reads them, through KIFF with its own Card, at most %d reads a day; each read has KIFF's permit and is kept in KIFF's record.\n",
		And:          "and",
		KeyInHeaderF: "in the %s header",
		KeyInQueryF:  "as the %s query parameter",

		APIAsReadF: "Its API is as described at %s when it was read (SHA-256 %s…), at %s.",
		KeyVaultF:  "Your key is kept in a sealed vault and is never shown again; the taf never sees it. It's sent only to %s, %s.",
		DiscoverAPIF: "Connect %s with an API key you give, to see what it offers for %s.\n%s\n\n" +
			"Only its API description is read: none of its operations is called, nothing is read from your account, nothing is changed, and nothing is given to the taf. " +
			"The tools the taf may use come next, and you approve them again. If they aren't approved within a day, the key is deleted.\n\n" +
			"%s You can disconnect at any time, which deletes it.",
		WillGiveKeyF: "You will give an API key for %s.\n%s\n\n",
		KeyTail:      " Each call is sent only with KIFF's permit for that exact call. You can disconnect at any time, which deletes it.",

		FirstToolsF: "These are the taf's first tools on this connection: %s. Setup's own reads keep working while you decide.\n\n",
		AddedToolsF: "This adds to a connection the taf already has: %s. What it does today keeps working while you decide.\n\n",

		ConnectTitledF:   "Connect %[1]s for %[2]s. Your sign-in to %[1]s stays with the connection service; the taf never sees it.",
		ConnectReusedF:   "Connect %[1]s for %[2]s with the sign-in you already have there: no second sign-in or key. It stays with the connection service; the taf never sees it.",
		ReusedNext:       " The actions to approve come next. No actions are granted yet.",
		EnterKeyF:        "Connect %[1]s for %[2]s. You'll enter your %[1]s key on the connection service's page; it keeps the key, and the taf never sees it.",
		EnterTokenF:      "Connect %[1]s for %[2]s. You'll enter your %[1]s token on the connection service's page; it keeps the token, and the taf never sees it.",
		SignInThenChoose: " Sign in and confirm on this device, then choose the actions to approve. No actions are granted yet.",
		OnlyYourAccountF: " Use only your connected %s account. Approve these actions:",
		MoreAccessF:      " %s will ask for more access than this taf uses. The taf can still only use: %s.",
		OnlyWhatItUsesF:  " %s is asked only for what this taf uses: %s.",

		AccessAtF:     "\n\nAccess at %s:",
		ReviewedFromF: "\nThe approved operations are the reviewed definitions from %s.",

		AnySiteF:     "%[1]s may use websites, in a browser planque runs for it alone. Any public site; nothing is kept between calls except a sign-in you make yourself.\n\n%[1]s may:\n",
		TheseSitesF:  "%[1]s may use these sites, in a browser planque runs for it alone: %[2]s. No other site; nothing is kept between calls except a sign-in you make yourself.\n\n%[1]s may:\n",
		BrowserFills: "fills in forms and clicks, only within the taf's Card; never a payment, a password, a card number or a code",
		BrowserPaysF: "pays a checkout in %s only, with the total read from the page, held for your approval each time",
		BrowserTail: "\nWhen a site needs a sign-in, you sign in yourself, on your phone; the taf never sees a password, a card number, a code or a cookie. " +
			"Each call is sent only with KIFF's permit for that exact call. You can disconnect at any time, which deletes the kept sign-ins.",

		TelegramHookF: "Giving it points the bot's webhook at planque, replacing any the bot had: messages people send the bot in a private chat, " +
			"and group messages that mention it, go to %s as they arrive. Other messages are dropped.",
		TelegramVault: "The token is kept in a sealed vault and is never shown again; the taf never sees it. It's sent only to Telegram.",
		TelegramConnectF: "Connect a Telegram bot with its token from BotFather, for %s. Telegram checks the token first.\n\n" +
			"%s The taf can't answer until you approve its tools next. If they aren't approved within a day, the token is deleted.\n\n" +
			"%s You can disconnect at any time, which deletes the webhook and the token.",
		TelegramWillGiveF: "You will give a Telegram bot's token from BotFather.\n%s\n\n",
		TelegramTail:      " Each message is sent only with KIFF's permit for that exact call. You can disconnect at any time, which deletes the webhook and the token.",

		PausedF:          "Paused · %s · connect again to approve it",
		Connected:        "Connected",
		NewToolsWaiting:  "Connected · new tools waiting for your approval",
		ToolsWaiting:     "Signed in · tools waiting for your approval",
		WaitingForYou:    "Waiting for your approval",
		DoorEvents:       "Connected · door events enabled",
		ToolsBeingChosen: "Signed in · the tools are being chosen",
		SignedOutByF:     "Signed out by %s · sign in again",
		ApprovedWaitF:    "Approved · waiting for the sign-in at %s",
		WaitingToConfirm: "Signed in · waiting for you to confirm",
		NotConnected:     "Not connected",
		Disconnected:     "Disconnected",
	},
	"de": {
		ThisTaf: "diesen taf", ThisTafFirst: "Dieser taf",

		DiscoverF: "Verbinde dein Konto bei %s, um zu sehen, was es für %s anbietet.\nDu meldest dich bei %s an.\n\n" +
			"Gelesen wird nur, was es anbietet: seine Liste von Tools. Keines davon wird aufgerufen: Aus deinem Konto wird nichts gelesen, nichts wird geändert, und der taf bekommt nichts. " +
			"Die Tools, die der taf nutzen darf, kommen als Nächstes, und du genehmigst sie erneut. Werden sie nicht innerhalb eines Tages genehmigt, wird die Anmeldung widerrufen.\n\n" +
			"Deine Anmeldung liegt in einem versiegelten Tresor; der taf sieht sie nie. Du kannst jederzeit trennen.",
		ConnectionVault: "\nDeine Anmeldung liegt in einem versiegelten Tresor; der taf sieht sie nie. Jeder Aufruf wird nur mit KIFFs Freigabe für genau diesen Aufruf gesendet. Du kannst jederzeit trennen.",

		KeyGiven:     "Der Schlüssel, den du gegeben hast, wird genutzt: kein zweiter Schlüssel.\n\n",
		SignedInToF:  "Du bist schon angemeldet, beim Konto %s: keine zweite Anmeldung.\n\n",
		SignedIn:     "Du bist schon angemeldet: keine zweite Anmeldung.\n\n",
		WillSignInF:  "Du meldest dich bei %s an.\n\n",
		MayF:         "%s darf:\n",
		Changes:      "nimmt Änderungen vor, nur im Rahmen der Card des tafs",
		ReadsOnly:    "liest nur",
		AnyChange:    "kann jede Änderung vornehmen, die dieses Tool erlaubt; die Card des tafs kann jeden Aufruf erlauben oder ablehnen, aber nicht begrenzen, was er ändert",
		MovesMoneyF:  "bewegt Geld (%s), nur im Rahmen der Card des tafs",
		SetupReadF:   "\nLesezugriffe von Setup, um das einzurichten: %s. Setup liest es nur, über KIFF mit seiner eigenen Card, höchstens %d Lesezugriffe am Tag; jeder hat KIFFs Freigabe und steht in KIFFs Protokoll.\n",
		SetupReadsF:  "\nLesezugriffe von Setup, um das einzurichten: %s. Setup liest sie nur, über KIFF mit seiner eigenen Card, höchstens %d Lesezugriffe am Tag; jeder hat KIFFs Freigabe und steht in KIFFs Protokoll.\n",
		And:          "und",
		KeyInHeaderF: "im Header %s",
		KeyInQueryF:  "als Abfrageparameter %s",

		APIAsReadF: "Seine API ist so, wie sie unter %s beim Lesen beschrieben war (SHA-256 %s…), unter %s.",
		KeyVaultF:  "Dein Schlüssel liegt in einem versiegelten Tresor und wird nie wieder angezeigt; der taf sieht ihn nie. Er wird nur an %s gesendet, %s.",
		DiscoverAPIF: "Verbinde %s mit einem API-Schlüssel, den du gibst, um zu sehen, was es für %s anbietet.\n%s\n\n" +
			"Gelesen wird nur seine API-Beschreibung: Keine seiner Operationen wird aufgerufen, aus deinem Konto wird nichts gelesen, nichts wird geändert, und der taf bekommt nichts. " +
			"Die Tools, die der taf nutzen darf, kommen als Nächstes, und du genehmigst sie erneut. Werden sie nicht innerhalb eines Tages genehmigt, wird der Schlüssel gelöscht.\n\n" +
			"%s Du kannst jederzeit trennen, wodurch er gelöscht wird.",
		WillGiveKeyF: "Du gibst einen API-Schlüssel für %s.\n%s\n\n",
		KeyTail:      " Jeder Aufruf wird nur mit KIFFs Freigabe für genau diesen Aufruf gesendet. Du kannst jederzeit trennen, wodurch er gelöscht wird.",

		FirstToolsF: "Das sind die ersten Tools des tafs in dieser Verbindung: %s. Die eigenen Lesezugriffe von Setup laufen weiter, während du entscheidest.\n\n",
		AddedToolsF: "Das kommt zu einer Verbindung hinzu, die der taf schon hat: %s. Was er heute tut, läuft weiter, während du entscheidest.\n\n",

		ConnectTitledF:   "Verbinde %[1]s für %[2]s. Deine Anmeldung bei %[1]s bleibt beim Verbindungsdienst; der taf sieht sie nie.",
		ConnectReusedF:   "Verbinde %[1]s für %[2]s mit der Anmeldung, die du dort schon hast: keine zweite Anmeldung und kein Schlüssel. Sie bleibt beim Verbindungsdienst; der taf sieht sie nie.",
		ReusedNext:       " Die Aktionen zum Genehmigen kommen als Nächstes. Noch ist keine Aktion gewährt.",
		EnterKeyF:        "Verbinde %[1]s für %[2]s. Du gibst deinen %[1]s-Schlüssel auf der Seite des Verbindungsdienstes ein; er behält den Schlüssel, und der taf sieht ihn nie.",
		EnterTokenF:      "Verbinde %[1]s für %[2]s. Du gibst dein %[1]s-Token auf der Seite des Verbindungsdienstes ein; er behält das Token, und der taf sieht es nie.",
		SignInThenChoose: " Melde dich an und bestätige auf diesem Gerät, dann wähle die Aktionen zum Genehmigen. Noch ist keine Aktion gewährt.",
		OnlyYourAccountF: " Nur dein verbundenes %s-Konto nutzen. Diese Aktionen genehmigen:",
		MoreAccessF:      " %s fragt nach mehr Zugriff, als dieser taf nutzt. Der taf kann trotzdem nur nutzen: %s.",
		OnlyWhatItUsesF:  " %s wird nur nach dem gefragt, was dieser taf nutzt: %s.",

		AccessAtF:     "\n\nZugriff bei %s:",
		ReviewedFromF: "\nDie genehmigten Operationen sind die geprüften Definitionen von %s.",

		AnySiteF:     "%[1]s darf Websites nutzen, in einem Browser, den planque nur für ihn betreibt. Jede öffentliche Seite; zwischen den Aufrufen wird nichts behalten außer einer Anmeldung, die du selbst machst.\n\n%[1]s darf:\n",
		TheseSitesF:  "%[1]s darf diese Seiten nutzen, in einem Browser, den planque nur für ihn betreibt: %[2]s. Keine andere Seite; zwischen den Aufrufen wird nichts behalten außer einer Anmeldung, die du selbst machst.\n\n%[1]s darf:\n",
		BrowserFills: "füllt Formulare aus und klickt, nur im Rahmen der Card des tafs; nie eine Zahlung, ein Passwort, eine Kartennummer oder ein Code",
		BrowserPaysF: "bezahlt einen Checkout nur in %s, mit dem Betrag von der Seite, jedes Mal zu deiner Genehmigung zurückgehalten",
		BrowserTail: "\nWenn eine Seite eine Anmeldung braucht, meldest du dich selbst an, auf deinem Telefon; der taf sieht nie ein Passwort, eine Kartennummer, einen Code oder ein Cookie. " +
			"Jeder Aufruf wird nur mit KIFFs Freigabe für genau diesen Aufruf gesendet. Du kannst jederzeit trennen, wodurch die behaltenen Anmeldungen gelöscht werden.",

		TelegramHookF: "Damit zeigt der Webhook des Bots auf planque und ersetzt einen, den der Bot hatte: Nachrichten, die Leute dem Bot privat schicken, " +
			"und Gruppennachrichten, die ihn erwähnen, gehen an %s, sobald sie ankommen. Andere Nachrichten werden verworfen.",
		TelegramVault: "Das Token liegt in einem versiegelten Tresor und wird nie wieder angezeigt; der taf sieht es nie. Es wird nur an Telegram gesendet.",
		TelegramConnectF: "Verbinde einen Telegram-Bot mit seinem Token von BotFather, für %s. Telegram prüft das Token zuerst.\n\n" +
			"%s Der taf kann nicht antworten, bis du als Nächstes seine Tools genehmigst. Werden sie nicht innerhalb eines Tages genehmigt, wird das Token gelöscht.\n\n" +
			"%s Du kannst jederzeit trennen, wodurch der Webhook und das Token gelöscht werden.",
		TelegramWillGiveF: "Du gibst das Token eines Telegram-Bots von BotFather.\n%s\n\n",
		TelegramTail:      " Jede Nachricht wird nur mit KIFFs Freigabe für genau diesen Aufruf gesendet. Du kannst jederzeit trennen, wodurch der Webhook und das Token gelöscht werden.",

		PausedF:          "Pausiert · %s · verbinde erneut, um es zu genehmigen",
		Connected:        "Verbunden",
		NewToolsWaiting:  "Verbunden · neue Tools warten auf deine Genehmigung",
		ToolsWaiting:     "Angemeldet · Tools warten auf deine Genehmigung",
		WaitingForYou:    "Wartet auf deine Genehmigung",
		DoorEvents:       "Verbunden · Tür-Ereignisse aktiviert",
		ToolsBeingChosen: "Angemeldet · die Tools werden gewählt",
		SignedOutByF:     "Von %s abgemeldet · erneut anmelden",
		ApprovedWaitF:    "Genehmigt · wartet auf die Anmeldung bei %s",
		WaitingToConfirm: "Angemeldet · wartet auf deine Bestätigung",
		NotConnected:     "Nicht verbunden",
		Disconnected:     "Getrennt",
	},
	"fr": {
		ThisTaf: "ce taf", ThisTafFirst: "Ce taf",

		DiscoverF: "Connectez votre compte chez %s pour voir ce qu'il propose pour %s.\nVous vous connecterez sur %s.\n\n" +
			"Seul ce qu'il propose est lu : sa liste d'outils. Aucun n'est appelé : rien n'est lu dans votre compte, rien n'est modifié et rien n'est donné au taf. " +
			"Les outils que le taf pourra utiliser viennent ensuite, et vous les approuvez à nouveau. S'ils ne sont pas approuvés dans la journée, la connexion est révoquée.\n\n" +
			"Votre connexion est gardée dans un coffre scellé ; le taf ne la voit jamais. Vous pouvez déconnecter à tout moment.",
		ConnectionVault: "\nVotre connexion est gardée dans un coffre scellé ; le taf ne la voit jamais. Chaque appel n'est envoyé qu'avec l'autorisation de KIFF pour cet appel précis. Vous pouvez déconnecter à tout moment.",

		KeyGiven:     "La clé que vous avez donnée est utilisée : pas de seconde clé.\n\n",
		SignedInToF:  "Vous êtes déjà connecté, au compte %s : pas de seconde connexion.\n\n",
		SignedIn:     "Vous êtes déjà connecté : pas de seconde connexion.\n\n",
		WillSignInF:  "Vous vous connecterez sur %s.\n\n",
		MayF:         "%s pourra :\n",
		Changes:      "fait des modifications, uniquement dans le cadre de la Card du taf",
		ReadsOnly:    "lecture seule",
		AnyChange:    "peut faire toute modification que cet outil permet ; la Card du taf peut autoriser ou refuser chaque appel, pas limiter ce qu'il modifie",
		MovesMoneyF:  "déplace de l'argent (%s), uniquement dans le cadre de la Card du taf",
		SetupReadF:   "\nLectures de Setup, pour la mise en place : %s. Setup ne fait que le lire, via KIFF avec sa propre Card, au plus %d lectures par jour ; chaque lecture a l'autorisation de KIFF et est gardée dans le registre de KIFF.\n",
		SetupReadsF:  "\nLectures de Setup, pour la mise en place : %s. Setup ne fait que les lire, via KIFF avec sa propre Card, au plus %d lectures par jour ; chaque lecture a l'autorisation de KIFF et est gardée dans le registre de KIFF.\n",
		And:          "et",
		KeyInHeaderF: "dans l'en-tête %s",
		KeyInQueryF:  "comme paramètre de requête %s",

		APIAsReadF: "Son API est telle que décrite sur %s au moment de sa lecture (SHA-256 %s…), sur %s.",
		KeyVaultF:  "Votre clé est gardée dans un coffre scellé et n'est plus jamais affichée ; le taf ne la voit jamais. Elle n'est envoyée qu'à %s, %s.",
		DiscoverAPIF: "Connectez %s avec une clé d'API que vous donnez, pour voir ce qu'il propose pour %s.\n%s\n\n" +
			"Seule sa description d'API est lue : aucune de ses opérations n'est appelée, rien n'est lu dans votre compte, rien n'est modifié et rien n'est donné au taf. " +
			"Les outils que le taf pourra utiliser viennent ensuite, et vous les approuvez à nouveau. S'ils ne sont pas approuvés dans la journée, la clé est supprimée.\n\n" +
			"%s Vous pouvez déconnecter à tout moment, ce qui la supprime.",
		WillGiveKeyF: "Vous donnerez une clé d'API pour %s.\n%s\n\n",
		KeyTail:      " Chaque appel n'est envoyé qu'avec l'autorisation de KIFF pour cet appel précis. Vous pouvez déconnecter à tout moment, ce qui la supprime.",

		FirstToolsF: "Ce sont les premiers outils du taf sur cette connexion : %s. Les lectures de Setup continuent pendant que vous décidez.\n\n",
		AddedToolsF: "Ceci s'ajoute à une connexion que le taf a déjà : %s. Ce qu'il fait aujourd'hui continue pendant que vous décidez.\n\n",

		ConnectTitledF:   "Connectez %[1]s pour %[2]s. Votre connexion à %[1]s reste chez le service de connexion ; le taf ne la voit jamais.",
		ConnectReusedF:   "Connectez %[1]s pour %[2]s avec la connexion que vous y avez déjà : pas de seconde connexion ni de clé. Elle reste chez le service de connexion ; le taf ne la voit jamais.",
		ReusedNext:       " Les actions à approuver viennent ensuite. Aucune action n'est encore accordée.",
		EnterKeyF:        "Connectez %[1]s pour %[2]s. Vous saisirez votre clé %[1]s sur la page du service de connexion ; il garde la clé, et le taf ne la voit jamais.",
		EnterTokenF:      "Connectez %[1]s pour %[2]s. Vous saisirez votre jeton %[1]s sur la page du service de connexion ; il garde le jeton, et le taf ne le voit jamais.",
		SignInThenChoose: " Connectez-vous et confirmez sur cet appareil, puis choisissez les actions à approuver. Aucune action n'est encore accordée.",
		OnlyYourAccountF: " Utiliser uniquement votre compte %s connecté. Approuver ces actions :",
		MoreAccessF:      " %s demandera plus d'accès que ce taf n'en utilise. Le taf ne peut toujours utiliser que : %s.",
		OnlyWhatItUsesF:  " %s n'est sollicité que pour ce que ce taf utilise : %s.",

		AccessAtF:     "\n\nAccès chez %s :",
		ReviewedFromF: "\nLes opérations approuvées sont les définitions revues de %s.",

		AnySiteF:     "%[1]s pourra utiliser des sites web, dans un navigateur que planque fait tourner pour lui seul. N'importe quel site public ; rien n'est gardé entre les appels, sauf une connexion que vous faites vous-même.\n\n%[1]s pourra :\n",
		TheseSitesF:  "%[1]s pourra utiliser ces sites, dans un navigateur que planque fait tourner pour lui seul : %[2]s. Aucun autre site ; rien n'est gardé entre les appels, sauf une connexion que vous faites vous-même.\n\n%[1]s pourra :\n",
		BrowserFills: "remplit des formulaires et clique, uniquement dans le cadre de la Card du taf ; jamais un paiement, un mot de passe, un numéro de carte ou un code",
		BrowserPaysF: "paie une commande en %s uniquement, avec le total lu sur la page, retenu pour votre approbation à chaque fois",
		BrowserTail: "\nQuand un site demande une connexion, vous vous connectez vous-même, sur votre téléphone ; le taf ne voit jamais un mot de passe, un numéro de carte, un code ou un cookie. " +
			"Chaque appel n'est envoyé qu'avec l'autorisation de KIFF pour cet appel précis. Vous pouvez déconnecter à tout moment, ce qui supprime les connexions gardées.",

		TelegramHookF: "Le donner fait pointer le webhook du bot vers planque, en remplaçant celui que le bot avait : les messages qu'on envoie au bot en privé, " +
			"et les messages de groupe qui le mentionnent, vont à %s dès leur arrivée. Les autres messages sont ignorés.",
		TelegramVault: "Le jeton est gardé dans un coffre scellé et n'est plus jamais affiché ; le taf ne le voit jamais. Il n'est envoyé qu'à Telegram.",
		TelegramConnectF: "Connectez un bot Telegram avec son jeton de BotFather, pour %s. Telegram vérifie d'abord le jeton.\n\n" +
			"%s Le taf ne peut pas répondre tant que vous n'approuvez pas ses outils, à l'étape suivante. S'ils ne sont pas approuvés dans la journée, le jeton est supprimé.\n\n" +
			"%s Vous pouvez déconnecter à tout moment, ce qui supprime le webhook et le jeton.",
		TelegramWillGiveF: "Vous donnerez le jeton d'un bot Telegram, de BotFather.\n%s\n\n",
		TelegramTail:      " Chaque message n'est envoyé qu'avec l'autorisation de KIFF pour cet appel précis. Vous pouvez déconnecter à tout moment, ce qui supprime le webhook et le jeton.",

		PausedF:          "En pause · %s · connectez à nouveau pour l'approuver",
		Connected:        "Connecté",
		NewToolsWaiting:  "Connecté · de nouveaux outils attendent votre approbation",
		ToolsWaiting:     "Connexion faite · des outils attendent votre approbation",
		WaitingForYou:    "En attente de votre approbation",
		DoorEvents:       "Connecté · événements de porte activés",
		ToolsBeingChosen: "Connexion faite · les outils sont en cours de choix",
		SignedOutByF:     "Déconnecté par %s · reconnectez-vous",
		ApprovedWaitF:    "Approuvé · en attente de la connexion chez %s",
		WaitingToConfirm: "Connexion faite · en attente de votre confirmation",
		NotConnected:     "Non connecté",
		Disconnected:     "Déconnecté",
	},
	"es": {
		ThisTaf: "este taf", ThisTafFirst: "Este taf",

		DiscoverF: "Conecta tu cuenta en %s para ver lo que ofrece para %s.\nIniciarás sesión en %s.\n\n" +
			"Solo se lee lo que ofrece: su lista de herramientas. No se llama a ninguna: no se lee nada de tu cuenta, no se cambia nada y no se le da nada al taf. " +
			"Las herramientas que el taf puede usar vienen después, y las vuelves a aprobar. Si no se aprueban en un día, se revoca el inicio de sesión.\n\n" +
			"Tu inicio de sesión se guarda en una caja fuerte sellada; el taf nunca lo ve. Puedes desconectar en cualquier momento.",
		ConnectionVault: "\nTu inicio de sesión se guarda en una caja fuerte sellada; el taf nunca lo ve. Cada llamada se envía solo con el permiso de KIFF para esa llamada exacta. Puedes desconectar en cualquier momento.",

		KeyGiven:     "Se usa la clave que diste: no hay una segunda clave.\n\n",
		SignedInToF:  "Ya has iniciado sesión, en la cuenta %s: no hay un segundo inicio de sesión.\n\n",
		SignedIn:     "Ya has iniciado sesión: no hay un segundo inicio de sesión.\n\n",
		WillSignInF:  "Iniciarás sesión en %s.\n\n",
		MayF:         "%s podrá:\n",
		Changes:      "hace cambios, solo dentro de la Card del taf",
		ReadsOnly:    "solo lee",
		AnyChange:    "puede hacer cualquier cambio que esta herramienta permita; la Card del taf puede permitir o rechazar cada llamada, no limitar lo que cambia",
		MovesMoneyF:  "mueve dinero (%s), solo dentro de la Card del taf",
		SetupReadF:   "\nLecturas de Setup, para configurarlo: %s. Setup solo lo lee, a través de KIFF con su propia Card, como mucho %d lecturas al día; cada lectura tiene el permiso de KIFF y queda en el registro de KIFF.\n",
		SetupReadsF:  "\nLecturas de Setup, para configurarlo: %s. Setup solo las lee, a través de KIFF con su propia Card, como mucho %d lecturas al día; cada lectura tiene el permiso de KIFF y queda en el registro de KIFF.\n",
		And:          "y",
		KeyInHeaderF: "en la cabecera %s",
		KeyInQueryF:  "como parámetro de consulta %s",

		APIAsReadF: "Su API es la descrita en %s cuando se leyó (SHA-256 %s…), en %s.",
		KeyVaultF:  "Tu clave se guarda en una caja fuerte sellada y no se vuelve a mostrar; el taf nunca la ve. Solo se envía a %s, %s.",
		DiscoverAPIF: "Conecta %s con una clave de API que das, para ver lo que ofrece para %s.\n%s\n\n" +
			"Solo se lee su descripción de API: no se llama a ninguna de sus operaciones, no se lee nada de tu cuenta, no se cambia nada y no se le da nada al taf. " +
			"Las herramientas que el taf puede usar vienen después, y las vuelves a aprobar. Si no se aprueban en un día, se borra la clave.\n\n" +
			"%s Puedes desconectar en cualquier momento, lo que la borra.",
		WillGiveKeyF: "Darás una clave de API para %s.\n%s\n\n",
		KeyTail:      " Cada llamada se envía solo con el permiso de KIFF para esa llamada exacta. Puedes desconectar en cualquier momento, lo que la borra.",

		FirstToolsF: "Estas son las primeras herramientas del taf en esta conexión: %s. Las lecturas propias de Setup siguen funcionando mientras decides.\n\n",
		AddedToolsF: "Esto se añade a una conexión que el taf ya tiene: %s. Lo que hace hoy sigue funcionando mientras decides.\n\n",

		ConnectTitledF:   "Conecta %[1]s para %[2]s. Tu inicio de sesión en %[1]s se queda en el servicio de conexión; el taf nunca lo ve.",
		ConnectReusedF:   "Conecta %[1]s para %[2]s con el inicio de sesión que ya tienes allí: sin un segundo inicio de sesión ni clave. Se queda en el servicio de conexión; el taf nunca lo ve.",
		ReusedNext:       " Las acciones que aprobar vienen después. Aún no se concede ninguna acción.",
		EnterKeyF:        "Conecta %[1]s para %[2]s. Introducirás tu clave de %[1]s en la página del servicio de conexión; este guarda la clave, y el taf nunca la ve.",
		EnterTokenF:      "Conecta %[1]s para %[2]s. Introducirás tu token de %[1]s en la página del servicio de conexión; este guarda el token, y el taf nunca lo ve.",
		SignInThenChoose: " Inicia sesión y confirma en este dispositivo; luego elige las acciones que aprobar. Aún no se concede ninguna acción.",
		OnlyYourAccountF: " Usar solo tu cuenta de %s conectada. Aprobar estas acciones:",
		MoreAccessF:      " %s pedirá más acceso del que usa este taf. El taf solo puede usar: %s.",
		OnlyWhatItUsesF:  " A %s solo se le pide lo que usa este taf: %s.",

		AccessAtF:     "\n\nAcceso en %s:",
		ReviewedFromF: "\nLas operaciones aprobadas son las definiciones revisadas de %s.",

		AnySiteF:     "%[1]s podrá usar sitios web, en un navegador que planque ejecuta solo para él. Cualquier sitio público; no se guarda nada entre llamadas, salvo un inicio de sesión que hagas tú.\n\n%[1]s podrá:\n",
		TheseSitesF:  "%[1]s podrá usar estos sitios, en un navegador que planque ejecuta solo para él: %[2]s. Ningún otro sitio; no se guarda nada entre llamadas, salvo un inicio de sesión que hagas tú.\n\n%[1]s podrá:\n",
		BrowserFills: "rellena formularios y hace clic, solo dentro de la Card del taf; nunca un pago, una contraseña, un número de tarjeta ni un código",
		BrowserPaysF: "paga un pedido solo en %s, con el total leído de la página, retenido para tu aprobación cada vez",
		BrowserTail: "\nCuando un sitio pide iniciar sesión, la inicias tú, en tu teléfono; el taf nunca ve una contraseña, un número de tarjeta, un código ni una cookie. " +
			"Cada llamada se envía solo con el permiso de KIFF para esa llamada exacta. Puedes desconectar en cualquier momento, lo que borra los inicios de sesión guardados.",

		TelegramHookF: "Darlo apunta el webhook del bot a planque, en lugar del que tuviera: los mensajes que la gente envía al bot en un chat privado, " +
			"y los mensajes de grupo que lo mencionan, van a %s en cuanto llegan. Los demás mensajes se descartan.",
		TelegramVault: "El token se guarda en una caja fuerte sellada y no se vuelve a mostrar; el taf nunca lo ve. Solo se envía a Telegram.",
		TelegramConnectF: "Conecta un bot de Telegram con su token de BotFather, para %s. Telegram comprueba primero el token.\n\n" +
			"%s El taf no puede responder hasta que apruebes sus herramientas, a continuación. Si no se aprueban en un día, se borra el token.\n\n" +
			"%s Puedes desconectar en cualquier momento, lo que borra el webhook y el token.",
		TelegramWillGiveF: "Darás el token de un bot de Telegram, de BotFather.\n%s\n\n",
		TelegramTail:      " Cada mensaje se envía solo con el permiso de KIFF para esa llamada exacta. Puedes desconectar en cualquier momento, lo que borra el webhook y el token.",

		PausedF:          "En pausa · %s · vuelve a conectar para aprobarlo",
		Connected:        "Conectado",
		NewToolsWaiting:  "Conectado · nuevas herramientas esperan tu aprobación",
		ToolsWaiting:     "Sesión iniciada · herramientas esperan tu aprobación",
		WaitingForYou:    "Esperando tu aprobación",
		DoorEvents:       "Conectado · eventos de puerta activados",
		ToolsBeingChosen: "Sesión iniciada · se están eligiendo las herramientas",
		SignedOutByF:     "Sesión cerrada por %s · vuelve a iniciar sesión",
		ApprovedWaitF:    "Aprobado · esperando el inicio de sesión en %s",
		WaitingToConfirm: "Sesión iniciada · esperando tu confirmación",
		NotConnected:     "No conectado",
		Disconnected:     "Desconectado",
	},
	"it": {
		ThisTaf: "questo taf", ThisTafFirst: "Questo taf",

		DiscoverF: "Collega il tuo account su %s per vedere cosa offre per %s.\nAccederai su %s.\n\n" +
			"Si legge solo ciò che offre: il suo elenco di strumenti. Nessuno viene chiamato: non si legge nulla dal tuo account, non si cambia nulla e non si dà nulla al taf. " +
			"Gli strumenti che il taf può usare vengono dopo, e li approvi di nuovo. Se non vengono approvati entro un giorno, l'accesso viene revocato.\n\n" +
			"Il tuo accesso è custodito in una cassaforte sigillata; il taf non lo vede mai. Puoi scollegarlo in qualsiasi momento.",
		ConnectionVault: "\nIl tuo accesso è custodito in una cassaforte sigillata; il taf non lo vede mai. Ogni chiamata è inviata solo con il permesso di KIFF per quella chiamata esatta. Puoi scollegarlo in qualsiasi momento.",

		KeyGiven:     "Si usa la chiave che hai dato: nessuna seconda chiave.\n\n",
		SignedInToF:  "Hai già effettuato l'accesso, all'account %s: nessun secondo accesso.\n\n",
		SignedIn:     "Hai già effettuato l'accesso: nessun secondo accesso.\n\n",
		WillSignInF:  "Accederai su %s.\n\n",
		MayF:         "%s potrà:\n",
		Changes:      "apporta modifiche, solo entro la Card del taf",
		ReadsOnly:    "solo lettura",
		AnyChange:    "può apportare qualsiasi modifica che questo strumento consente; la Card del taf può consentire o rifiutare ogni chiamata, non limitare ciò che cambia",
		MovesMoneyF:  "sposta denaro (%s), solo entro la Card del taf",
		SetupReadF:   "\nLetture di Setup, per configurarlo: %s. Setup lo legge soltanto, tramite KIFF con la sua Card, al massimo %d letture al giorno; ogni lettura ha il permesso di KIFF ed è conservata nel registro di KIFF.\n",
		SetupReadsF:  "\nLetture di Setup, per configurarlo: %s. Setup li legge soltanto, tramite KIFF con la sua Card, al massimo %d letture al giorno; ogni lettura ha il permesso di KIFF ed è conservata nel registro di KIFF.\n",
		And:          "e",
		KeyInHeaderF: "nell'intestazione %s",
		KeyInQueryF:  "come parametro di query %s",

		APIAsReadF: "La sua API è quella descritta su %s quando è stata letta (SHA-256 %s…), su %s.",
		KeyVaultF:  "La tua chiave è custodita in una cassaforte sigillata e non viene più mostrata; il taf non la vede mai. Viene inviata solo a %s, %s.",
		DiscoverAPIF: "Collega %s con una chiave API che dai, per vedere cosa offre per %s.\n%s\n\n" +
			"Si legge solo la sua descrizione dell'API: nessuna delle sue operazioni viene chiamata, non si legge nulla dal tuo account, non si cambia nulla e non si dà nulla al taf. " +
			"Gli strumenti che il taf può usare vengono dopo, e li approvi di nuovo. Se non vengono approvati entro un giorno, la chiave viene eliminata.\n\n" +
			"%s Puoi scollegarlo in qualsiasi momento, il che la elimina.",
		WillGiveKeyF: "Darai una chiave API per %s.\n%s\n\n",
		KeyTail:      " Ogni chiamata è inviata solo con il permesso di KIFF per quella chiamata esatta. Puoi scollegarlo in qualsiasi momento, il che la elimina.",

		FirstToolsF: "Questi sono i primi strumenti del taf su questo collegamento: %s. Le letture di Setup continuano mentre decidi.\n\n",
		AddedToolsF: "Questo si aggiunge a un collegamento che il taf ha già: %s. Ciò che fa oggi continua mentre decidi.\n\n",

		ConnectTitledF:   "Collega %[1]s per %[2]s. Il tuo accesso a %[1]s resta al servizio di collegamento; il taf non lo vede mai.",
		ConnectReusedF:   "Collega %[1]s per %[2]s con l'accesso che hai già lì: nessun secondo accesso né chiave. Resta al servizio di collegamento; il taf non lo vede mai.",
		ReusedNext:       " Le azioni da approvare vengono dopo. Nessuna azione è ancora concessa.",
		EnterKeyF:        "Collega %[1]s per %[2]s. Inserirai la tua chiave di %[1]s nella pagina del servizio di collegamento; il servizio conserva la chiave, e il taf non la vede mai.",
		EnterTokenF:      "Collega %[1]s per %[2]s. Inserirai il tuo token di %[1]s nella pagina del servizio di collegamento; il servizio conserva il token, e il taf non lo vede mai.",
		SignInThenChoose: " Accedi e conferma su questo dispositivo, poi scegli le azioni da approvare. Nessuna azione è ancora concessa.",
		OnlyYourAccountF: " Usa solo il tuo account %s collegato. Approva queste azioni:",
		MoreAccessF:      " %s chiederà più accesso di quanto usi questo taf. Il taf può comunque usare solo: %s.",
		OnlyWhatItUsesF:  " A %s si chiede solo ciò che usa questo taf: %s.",

		AccessAtF:     "\n\nAccesso su %s:",
		ReviewedFromF: "\nLe operazioni approvate sono le definizioni revisionate da %s.",

		AnySiteF:     "%[1]s potrà usare siti web, in un browser che planque esegue solo per lui. Qualsiasi sito pubblico; tra una chiamata e l'altra non si conserva nulla, tranne un accesso che fai tu.\n\n%[1]s potrà:\n",
		TheseSitesF:  "%[1]s potrà usare questi siti, in un browser che planque esegue solo per lui: %[2]s. Nessun altro sito; tra una chiamata e l'altra non si conserva nulla, tranne un accesso che fai tu.\n\n%[1]s potrà:\n",
		BrowserFills: "compila moduli e fa clic, solo entro la Card del taf; mai un pagamento, una password, un numero di carta o un codice",
		BrowserPaysF: "paga un checkout solo in %s, con il totale letto dalla pagina, trattenuto ogni volta per la tua approvazione",
		BrowserTail: "\nQuando un sito chiede un accesso, accedi tu, sul tuo telefono; il taf non vede mai una password, un numero di carta, un codice o un cookie. " +
			"Ogni chiamata è inviata solo con il permesso di KIFF per quella chiamata esatta. Puoi scollegarlo in qualsiasi momento, il che elimina gli accessi conservati.",

		TelegramHookF: "Darlo punta il webhook del bot su planque, sostituendo quello che il bot aveva: i messaggi che le persone inviano al bot in una chat privata, " +
			"e i messaggi di gruppo che lo menzionano, vanno a %s appena arrivano. Gli altri messaggi vengono scartati.",
		TelegramVault: "Il token è custodito in una cassaforte sigillata e non viene più mostrato; il taf non lo vede mai. Viene inviato solo a Telegram.",
		TelegramConnectF: "Collega un bot Telegram con il suo token di BotFather, per %s. Telegram verifica prima il token.\n\n" +
			"%s Il taf non può rispondere finché non approvi i suoi strumenti, al passo successivo. Se non vengono approvati entro un giorno, il token viene eliminato.\n\n" +
			"%s Puoi scollegarlo in qualsiasi momento, il che elimina il webhook e il token.",
		TelegramWillGiveF: "Darai il token di un bot Telegram, da BotFather.\n%s\n\n",
		TelegramTail:      " Ogni messaggio è inviato solo con il permesso di KIFF per quella chiamata esatta. Puoi scollegarlo in qualsiasi momento, il che elimina il webhook e il token.",

		PausedF:          "In pausa · %s · collega di nuovo per approvarlo",
		Connected:        "Collegato",
		NewToolsWaiting:  "Collegato · nuovi strumenti aspettano la tua approvazione",
		ToolsWaiting:     "Accesso effettuato · gli strumenti aspettano la tua approvazione",
		WaitingForYou:    "In attesa della tua approvazione",
		DoorEvents:       "Collegato · eventi della porta attivi",
		ToolsBeingChosen: "Accesso effettuato · si stanno scegliendo gli strumenti",
		SignedOutByF:     "Disconnesso da %s · accedi di nuovo",
		ApprovedWaitF:    "Approvato · in attesa dell'accesso su %s",
		WaitingToConfirm: "Accesso effettuato · in attesa della tua conferma",
		NotConnected:     "Non collegato",
		Disconnected:     "Scollegato",
	},
	"pt": {
		ThisTaf: "este taf", ThisTafFirst: "Este taf",

		DiscoverF: "Ligue a sua conta em %s para ver o que oferece para %s.\nVai iniciar sessão em %s.\n\n" +
			"Só se lê o que oferece: a sua lista de ferramentas. Nenhuma é chamada: nada é lido da sua conta, nada é alterado e nada é dado ao taf. " +
			"As ferramentas que o taf pode usar vêm a seguir, e aprova-as novamente. Se não forem aprovadas no prazo de um dia, o início de sessão é revogado.\n\n" +
			"O seu início de sessão fica num cofre selado; o taf nunca o vê. Pode desligar a qualquer momento.",
		ConnectionVault: "\nO seu início de sessão fica num cofre selado; o taf nunca o vê. Cada chamada só é enviada com a autorização da KIFF para essa chamada exata. Pode desligar a qualquer momento.",

		KeyGiven:     "É usada a chave que deu: sem segunda chave.\n\n",
		SignedInToF:  "Já tem sessão iniciada, na conta %s: sem segundo início de sessão.\n\n",
		SignedIn:     "Já tem sessão iniciada: sem segundo início de sessão.\n\n",
		WillSignInF:  "Vai iniciar sessão em %s.\n\n",
		MayF:         "%s poderá:\n",
		Changes:      "faz alterações, só dentro da Card do taf",
		ReadsOnly:    "só lê",
		AnyChange:    "pode fazer qualquer alteração que esta ferramenta permita; a Card do taf pode permitir ou recusar cada chamada, não limitar o que altera",
		MovesMoneyF:  "movimenta dinheiro (%s), só dentro da Card do taf",
		SetupReadF:   "\nLeituras do Setup, para configurar isto: %s. O Setup só o lê, através da KIFF com a sua própria Card, no máximo %d leituras por dia; cada leitura tem a autorização da KIFF e fica no registo da KIFF.\n",
		SetupReadsF:  "\nLeituras do Setup, para configurar isto: %s. O Setup só as lê, através da KIFF com a sua própria Card, no máximo %d leituras por dia; cada leitura tem a autorização da KIFF e fica no registo da KIFF.\n",
		And:          "e",
		KeyInHeaderF: "no cabeçalho %s",
		KeyInQueryF:  "como parâmetro de consulta %s",

		APIAsReadF: "A sua API é a descrita em %s quando foi lida (SHA-256 %s…), em %s.",
		KeyVaultF:  "A sua chave fica num cofre selado e nunca mais é mostrada; o taf nunca a vê. Só é enviada para %s, %s.",
		DiscoverAPIF: "Ligue %s com uma chave de API que dá, para ver o que oferece para %s.\n%s\n\n" +
			"Só se lê a sua descrição de API: nenhuma das suas operações é chamada, nada é lido da sua conta, nada é alterado e nada é dado ao taf. " +
			"As ferramentas que o taf pode usar vêm a seguir, e aprova-as novamente. Se não forem aprovadas no prazo de um dia, a chave é apagada.\n\n" +
			"%s Pode desligar a qualquer momento, o que a apaga.",
		WillGiveKeyF: "Vai dar uma chave de API para %s.\n%s\n\n",
		KeyTail:      " Cada chamada só é enviada com a autorização da KIFF para essa chamada exata. Pode desligar a qualquer momento, o que a apaga.",

		FirstToolsF: "Estas são as primeiras ferramentas do taf nesta ligação: %s. As leituras do Setup continuam enquanto decide.\n\n",
		AddedToolsF: "Isto junta-se a uma ligação que o taf já tem: %s. O que faz hoje continua enquanto decide.\n\n",

		ConnectTitledF:   "Ligue %[1]s para %[2]s. O seu início de sessão em %[1]s fica no serviço de ligação; o taf nunca o vê.",
		ConnectReusedF:   "Ligue %[1]s para %[2]s com o início de sessão que já tem lá: sem segundo início de sessão nem chave. Fica no serviço de ligação; o taf nunca o vê.",
		ReusedNext:       " As ações a aprovar vêm a seguir. Ainda não é concedida nenhuma ação.",
		EnterKeyF:        "Ligue %[1]s para %[2]s. Vai introduzir a sua chave de %[1]s na página do serviço de ligação; este guarda a chave, e o taf nunca a vê.",
		EnterTokenF:      "Ligue %[1]s para %[2]s. Vai introduzir o seu token de %[1]s na página do serviço de ligação; este guarda o token, e o taf nunca o vê.",
		SignInThenChoose: " Inicie sessão e confirme neste dispositivo; depois escolha as ações a aprovar. Ainda não é concedida nenhuma ação.",
		OnlyYourAccountF: " Usar só a sua conta de %s ligada. Aprovar estas ações:",
		MoreAccessF:      " %s vai pedir mais acesso do que este taf usa. O taf continua a poder usar só: %s.",
		OnlyWhatItUsesF:  " A %s só é pedido o que este taf usa: %s.",

		AccessAtF:     "\n\nAcesso em %s:",
		ReviewedFromF: "\nAs operações aprovadas são as definições revistas de %s.",

		AnySiteF:     "%[1]s poderá usar sites, num navegador que planque executa só para ele. Qualquer site público; nada é guardado entre chamadas, exceto um início de sessão que faça você mesmo.\n\n%[1]s poderá:\n",
		TheseSitesF:  "%[1]s poderá usar estes sites, num navegador que planque executa só para ele: %[2]s. Nenhum outro site; nada é guardado entre chamadas, exceto um início de sessão que faça você mesmo.\n\n%[1]s poderá:\n",
		BrowserFills: "preenche formulários e clica, só dentro da Card do taf; nunca um pagamento, uma palavra-passe, um número de cartão ou um código",
		BrowserPaysF: "paga um checkout só em %s, com o total lido da página, retido para a sua aprovação de cada vez",
		BrowserTail: "\nQuando um site pede um início de sessão, é você que inicia sessão, no seu telemóvel; o taf nunca vê uma palavra-passe, um número de cartão, um código ou um cookie. " +
			"Cada chamada só é enviada com a autorização da KIFF para essa chamada exata. Pode desligar a qualquer momento, o que apaga os inícios de sessão guardados.",

		TelegramHookF: "Dá-lo aponta o webhook do bot para planque, substituindo o que o bot tinha: as mensagens que as pessoas enviam ao bot num chat privado, " +
			"e as mensagens de grupo que o mencionam, vão para %s assim que chegam. As outras mensagens são descartadas.",
		TelegramVault: "O token fica num cofre selado e nunca mais é mostrado; o taf nunca o vê. Só é enviado para o Telegram.",
		TelegramConnectF: "Ligue um bot do Telegram com o seu token do BotFather, para %s. O Telegram verifica primeiro o token.\n\n" +
			"%s O taf não pode responder até aprovar as suas ferramentas, a seguir. Se não forem aprovadas no prazo de um dia, o token é apagado.\n\n" +
			"%s Pode desligar a qualquer momento, o que apaga o webhook e o token.",
		TelegramWillGiveF: "Vai dar o token de um bot do Telegram, do BotFather.\n%s\n\n",
		TelegramTail:      " Cada mensagem só é enviada com a autorização da KIFF para essa chamada exata. Pode desligar a qualquer momento, o que apaga o webhook e o token.",

		PausedF:          "Em pausa · %s · ligue novamente para aprovar",
		Connected:        "Ligado",
		NewToolsWaiting:  "Ligado · novas ferramentas aguardam a sua aprovação",
		ToolsWaiting:     "Sessão iniciada · ferramentas aguardam a sua aprovação",
		WaitingForYou:    "A aguardar a sua aprovação",
		DoorEvents:       "Ligado · eventos da porta ativados",
		ToolsBeingChosen: "Sessão iniciada · as ferramentas estão a ser escolhidas",
		SignedOutByF:     "Sessão terminada por %s · inicie sessão novamente",
		ApprovedWaitF:    "Aprovado · a aguardar o início de sessão em %s",
		WaitingToConfirm: "Sessão iniciada · a aguardar a sua confirmação",
		NotConnected:     "Não ligado",
		Disconnected:     "Desligado",
	},
	"nl": {
		ThisTaf: "dit taf", ThisTafFirst: "Dit taf",

		DiscoverF: "Koppel je account bij %s om te zien wat het biedt voor %s.\nJe logt in bij %s.\n\n" +
			"Alleen wat het biedt wordt gelezen: de lijst met tools. Geen ervan wordt aangeroepen: er wordt niets uit je account gelezen, niets veranderd en niets aan het taf gegeven. " +
			"De tools die het taf mag gebruiken komen hierna, en die keur je opnieuw goed. Worden ze niet binnen een dag goedgekeurd, dan wordt de inlog ingetrokken.\n\n" +
			"Je inlog wordt in een verzegelde kluis bewaard; het taf ziet hem nooit. Je kunt op elk moment ontkoppelen.",
		ConnectionVault: "\nJe inlog wordt in een verzegelde kluis bewaard; het taf ziet hem nooit. Elke aanroep wordt alleen verzonden met de toestemming van KIFF voor precies die aanroep. Je kunt op elk moment ontkoppelen.",

		KeyGiven:     "De sleutel die je gaf wordt gebruikt: geen tweede sleutel.\n\n",
		SignedInToF:  "Je bent al ingelogd, op het account %s: geen tweede inlog.\n\n",
		SignedIn:     "Je bent al ingelogd: geen tweede inlog.\n\n",
		WillSignInF:  "Je logt in bij %s.\n\n",
		MayF:         "%s mag:\n",
		Changes:      "brengt wijzigingen aan, alleen binnen de Card van het taf",
		ReadsOnly:    "leest alleen",
		AnyChange:    "kan elke wijziging aanbrengen die deze tool toestaat; de Card van het taf kan elke aanroep toestaan of weigeren, niet beperken wat hij wijzigt",
		MovesMoneyF:  "verplaatst geld (%s), alleen binnen de Card van het taf",
		SetupReadF:   "\nLeesacties van Setup, om dit in te stellen: %s. Setup leest het alleen, via KIFF met een eigen Card, hoogstens %d leesacties per dag; elke leesactie heeft de toestemming van KIFF en staat in het logboek van KIFF.\n",
		SetupReadsF:  "\nLeesacties van Setup, om dit in te stellen: %s. Setup leest ze alleen, via KIFF met een eigen Card, hoogstens %d leesacties per dag; elke leesactie heeft de toestemming van KIFF en staat in het logboek van KIFF.\n",
		And:          "en",
		KeyInHeaderF: "in de header %s",
		KeyInQueryF:  "als queryparameter %s",

		APIAsReadF: "De API is zoals beschreven op %s toen die werd gelezen (SHA-256 %s…), op %s.",
		KeyVaultF:  "Je sleutel wordt in een verzegelde kluis bewaard en nooit meer getoond; het taf ziet hem nooit. Hij wordt alleen naar %s gestuurd, %s.",
		DiscoverAPIF: "Koppel %s met een API-sleutel die je geeft, om te zien wat het biedt voor %s.\n%s\n\n" +
			"Alleen de API-beschrijving wordt gelezen: geen van de bewerkingen wordt aangeroepen, er wordt niets uit je account gelezen, niets veranderd en niets aan het taf gegeven. " +
			"De tools die het taf mag gebruiken komen hierna, en die keur je opnieuw goed. Worden ze niet binnen een dag goedgekeurd, dan wordt de sleutel verwijderd.\n\n" +
			"%s Je kunt op elk moment ontkoppelen, waarmee hij wordt verwijderd.",
		WillGiveKeyF: "Je geeft een API-sleutel voor %s.\n%s\n\n",
		KeyTail:      " Elke aanroep wordt alleen verzonden met de toestemming van KIFF voor precies die aanroep. Je kunt op elk moment ontkoppelen, waarmee hij wordt verwijderd.",

		FirstToolsF: "Dit zijn de eerste tools van het taf op deze koppeling: %s. De eigen leesacties van Setup blijven werken terwijl je beslist.\n\n",
		AddedToolsF: "Dit komt bij een koppeling die het taf al heeft: %s. Wat het vandaag doet blijft werken terwijl je beslist.\n\n",

		ConnectTitledF:   "Koppel %[1]s voor %[2]s. Je inlog bij %[1]s blijft bij de koppelingsdienst; het taf ziet hem nooit.",
		ConnectReusedF:   "Koppel %[1]s voor %[2]s met de inlog die je daar al hebt: geen tweede inlog of sleutel. Hij blijft bij de koppelingsdienst; het taf ziet hem nooit.",
		ReusedNext:       " De acties om goed te keuren komen hierna. Er is nog geen actie toegekend.",
		EnterKeyF:        "Koppel %[1]s voor %[2]s. Je voert je %[1]s-sleutel in op de pagina van de koppelingsdienst; die bewaart de sleutel, en het taf ziet hem nooit.",
		EnterTokenF:      "Koppel %[1]s voor %[2]s. Je voert je %[1]s-token in op de pagina van de koppelingsdienst; die bewaart het token, en het taf ziet het nooit.",
		SignInThenChoose: " Log in en bevestig op dit apparaat, en kies dan de acties om goed te keuren. Er is nog geen actie toegekend.",
		OnlyYourAccountF: " Alleen je gekoppelde %s-account gebruiken. Deze acties goedkeuren:",
		MoreAccessF:      " %s vraagt om meer toegang dan dit taf gebruikt. Het taf kan nog steeds alleen gebruiken: %s.",
		OnlyWhatItUsesF:  " Aan %s wordt alleen gevraagd wat dit taf gebruikt: %s.",

		AccessAtF:     "\n\nToegang bij %s:",
		ReviewedFromF: "\nDe goedgekeurde bewerkingen zijn de beoordeelde definities van %s.",

		AnySiteF:     "%[1]s mag websites gebruiken, in een browser die planque alleen ervoor draait. Elke openbare site; tussen aanroepen wordt niets bewaard behalve een inlog die je zelf doet.\n\n%[1]s mag:\n",
		TheseSitesF:  "%[1]s mag deze sites gebruiken, in een browser die planque alleen ervoor draait: %[2]s. Geen andere site; tussen aanroepen wordt niets bewaard behalve een inlog die je zelf doet.\n\n%[1]s mag:\n",
		BrowserFills: "vult formulieren in en klikt, alleen binnen de Card van het taf; nooit een betaling, een wachtwoord, een kaartnummer of een code",
		BrowserPaysF: "betaalt een checkout alleen in %s, met het totaal van de pagina, elke keer vastgehouden voor je goedkeuring",
		BrowserTail: "\nAls een site een inlog nodig heeft, log je zelf in, op je telefoon; het taf ziet nooit een wachtwoord, een kaartnummer, een code of een cookie. " +
			"Elke aanroep wordt alleen verzonden met de toestemming van KIFF voor precies die aanroep. Je kunt op elk moment ontkoppelen, waarmee de bewaarde inlogs worden verwijderd.",

		TelegramHookF: "Daarmee wijst de webhook van de bot naar planque, in plaats van een die de bot had: berichten die mensen de bot in een privéchat sturen, " +
			"en groepsberichten die hem noemen, gaan naar %s zodra ze binnenkomen. Andere berichten worden genegeerd.",
		TelegramVault: "Het token wordt in een verzegelde kluis bewaard en nooit meer getoond; het taf ziet het nooit. Het wordt alleen naar Telegram gestuurd.",
		TelegramConnectF: "Koppel een Telegram-bot met het token van BotFather, voor %s. Telegram controleert het token eerst.\n\n" +
			"%s Het taf kan niet antwoorden tot je hierna de tools goedkeurt. Worden ze niet binnen een dag goedgekeurd, dan wordt het token verwijderd.\n\n" +
			"%s Je kunt op elk moment ontkoppelen, waarmee de webhook en het token worden verwijderd.",
		TelegramWillGiveF: "Je geeft het token van een Telegram-bot van BotFather.\n%s\n\n",
		TelegramTail:      " Elk bericht wordt alleen verzonden met de toestemming van KIFF voor precies die aanroep. Je kunt op elk moment ontkoppelen, waarmee de webhook en het token worden verwijderd.",

		PausedF:          "Gepauzeerd · %s · koppel opnieuw om het goed te keuren",
		Connected:        "Gekoppeld",
		NewToolsWaiting:  "Gekoppeld · nieuwe tools wachten op je goedkeuring",
		ToolsWaiting:     "Ingelogd · tools wachten op je goedkeuring",
		WaitingForYou:    "Wacht op je goedkeuring",
		DoorEvents:       "Gekoppeld · deurgebeurtenissen aan",
		ToolsBeingChosen: "Ingelogd · de tools worden gekozen",
		SignedOutByF:     "Uitgelogd door %s · opnieuw inloggen",
		ApprovedWaitF:    "Goedgekeurd · wacht op de inlog bij %s",
		WaitingToConfirm: "Ingelogd · wacht op je bevestiging",
		NotConnected:     "Niet gekoppeld",
		Disconnected:     "Ontkoppeld",
	},
	"pl": {
		ThisTaf: "tego tafa", ThisTafFirst: "Ten taf",

		DiscoverF: "Połącz swoje konto w %s, aby zobaczyć, co oferuje dla %s.\nZalogujesz się w %s.\n\n" +
			"Odczytywane jest tylko to, co oferuje: jego lista narzędzi. Żadne nie jest wywoływane: nic nie jest odczytywane z Twojego konta, nic nie jest zmieniane i nic nie trafia do tafa. " +
			"Narzędzia, których taf może używać, są następne i zatwierdzasz je ponownie. Jeśli nie zostaną zatwierdzone w ciągu doby, logowanie zostanie cofnięte.\n\n" +
			"Twoje logowanie jest przechowywane w zapieczętowanym sejfie; taf nigdy go nie widzi. Możesz rozłączyć w każdej chwili.",
		ConnectionVault: "\nTwoje logowanie jest przechowywane w zapieczętowanym sejfie; taf nigdy go nie widzi. Każde wywołanie jest wysyłane tylko ze zgodą KIFF na dokładnie to wywołanie. Możesz rozłączyć w każdej chwili.",

		KeyGiven:     "Używany jest klucz, który podałeś: bez drugiego klucza.\n\n",
		SignedInToF:  "Jesteś już zalogowany, na koncie %s: bez drugiego logowania.\n\n",
		SignedIn:     "Jesteś już zalogowany: bez drugiego logowania.\n\n",
		WillSignInF:  "Zalogujesz się w %s.\n\n",
		MayF:         "%s może:\n",
		Changes:      "wprowadza zmiany, tylko w ramach Card tafa",
		ReadsOnly:    "tylko odczyt",
		AnyChange:    "może wprowadzić każdą zmianę, na którą pozwala to narzędzie; Card tafa może zezwolić na każde wywołanie lub je odrzucić, ale nie ograniczyć tego, co zmienia",
		MovesMoneyF:  "przesyła pieniądze (%s), tylko w ramach Card tafa",
		SetupReadF:   "\nOdczyty Setup, aby to skonfigurować: %s. Setup tylko to odczytuje, przez KIFF z własną Card, najwyżej %d odczytów dziennie; każdy odczyt ma zgodę KIFF i jest zapisany w rejestrze KIFF.\n",
		SetupReadsF:  "\nOdczyty Setup, aby to skonfigurować: %s. Setup tylko je odczytuje, przez KIFF z własną Card, najwyżej %d odczytów dziennie; każdy odczyt ma zgodę KIFF i jest zapisany w rejestrze KIFF.\n",
		And:          "i",
		KeyInHeaderF: "w nagłówku %s",
		KeyInQueryF:  "jako parametr zapytania %s",

		APIAsReadF: "Jego API jest takie, jak opisano pod %s w chwili odczytu (SHA-256 %s…), pod %s.",
		KeyVaultF:  "Twój klucz jest przechowywany w zapieczętowanym sejfie i nigdy więcej nie jest pokazywany; taf nigdy go nie widzi. Jest wysyłany tylko do %s, %s.",
		DiscoverAPIF: "Połącz %s kluczem API, który podasz, aby zobaczyć, co oferuje dla %s.\n%s\n\n" +
			"Odczytywany jest tylko opis jego API: żadna z jego operacji nie jest wywoływana, nic nie jest odczytywane z Twojego konta, nic nie jest zmieniane i nic nie trafia do tafa. " +
			"Narzędzia, których taf może używać, są następne i zatwierdzasz je ponownie. Jeśli nie zostaną zatwierdzone w ciągu doby, klucz zostanie usunięty.\n\n" +
			"%s Możesz rozłączyć w każdej chwili, co go usuwa.",
		WillGiveKeyF: "Podasz klucz API dla %s.\n%s\n\n",
		KeyTail:      " Każde wywołanie jest wysyłane tylko ze zgodą KIFF na dokładnie to wywołanie. Możesz rozłączyć w każdej chwili, co go usuwa.",

		FirstToolsF: "To pierwsze narzędzia tafa w tym połączeniu: %s. Własne odczyty Setup działają dalej, gdy decydujesz.\n\n",
		AddedToolsF: "To dochodzi do połączenia, które taf już ma: %s. To, co robi dziś, działa dalej, gdy decydujesz.\n\n",

		ConnectTitledF:   "Połącz %[1]s dla %[2]s. Twoje logowanie do %[1]s zostaje w usłudze połączeń; taf nigdy go nie widzi.",
		ConnectReusedF:   "Połącz %[1]s dla %[2]s logowaniem, które już tam masz: bez drugiego logowania ani klucza. Zostaje w usłudze połączeń; taf nigdy go nie widzi.",
		ReusedNext:       " Działania do zatwierdzenia są następne. Żadne działanie nie jest jeszcze przyznane.",
		EnterKeyF:        "Połącz %[1]s dla %[2]s. Wpiszesz swój klucz %[1]s na stronie usługi połączeń; ona przechowuje klucz, a taf nigdy go nie widzi.",
		EnterTokenF:      "Połącz %[1]s dla %[2]s. Wpiszesz swój token %[1]s na stronie usługi połączeń; ona przechowuje token, a taf nigdy go nie widzi.",
		SignInThenChoose: " Zaloguj się i potwierdź na tym urządzeniu, a potem wybierz działania do zatwierdzenia. Żadne działanie nie jest jeszcze przyznane.",
		OnlyYourAccountF: " Używaj tylko swojego połączonego konta %s. Zatwierdź te działania:",
		MoreAccessF:      " %s poprosi o większy dostęp, niż używa ten taf. Taf nadal może używać tylko: %s.",
		OnlyWhatItUsesF:  " %s jest proszony tylko o to, czego używa ten taf: %s.",

		AccessAtF:     "\n\nDostęp w %s:",
		ReviewedFromF: "\nZatwierdzone operacje to sprawdzone definicje z %s.",

		AnySiteF:     "%[1]s może korzystać ze stron internetowych, w przeglądarce, którą planque uruchamia tylko dla niego. Dowolna publiczna strona; między wywołaniami nic nie jest zachowywane poza logowaniem, które zrobisz sam.\n\n%[1]s może:\n",
		TheseSitesF:  "%[1]s może korzystać z tych stron, w przeglądarce, którą planque uruchamia tylko dla niego: %[2]s. Żadnej innej strony; między wywołaniami nic nie jest zachowywane poza logowaniem, które zrobisz sam.\n\n%[1]s może:\n",
		BrowserFills: "wypełnia formularze i klika, tylko w ramach Card tafa; nigdy płatność, hasło, numer karty ani kod",
		BrowserPaysF: "płaci za zamówienie tylko w %s, kwotą odczytaną ze strony, za każdym razem wstrzymaną do Twojego zatwierdzenia",
		BrowserTail: "\nGdy strona wymaga logowania, logujesz się sam, na swoim telefonie; taf nigdy nie widzi hasła, numeru karty, kodu ani pliku cookie. " +
			"Każde wywołanie jest wysyłane tylko ze zgodą KIFF na dokładnie to wywołanie. Możesz rozłączyć w każdej chwili, co usuwa zachowane logowania.",

		TelegramHookF: "Podanie go kieruje webhook bota do planque, zastępując ten, który bot miał: wiadomości, które ludzie wysyłają botowi na czacie prywatnym, " +
			"i wiadomości grupowe, które go wspominają, trafiają do %s, gdy tylko dotrą. Inne wiadomości są odrzucane.",
		TelegramVault: "Token jest przechowywany w zapieczętowanym sejfie i nigdy więcej nie jest pokazywany; taf nigdy go nie widzi. Jest wysyłany tylko do Telegrama.",
		TelegramConnectF: "Połącz bota Telegrama jego tokenem od BotFather, dla %s. Telegram najpierw sprawdza token.\n\n" +
			"%s Taf nie może odpowiadać, dopóki w następnym kroku nie zatwierdzisz jego narzędzi. Jeśli nie zostaną zatwierdzone w ciągu doby, token zostanie usunięty.\n\n" +
			"%s Możesz rozłączyć w każdej chwili, co usuwa webhook i token.",
		TelegramWillGiveF: "Podasz token bota Telegrama od BotFather.\n%s\n\n",
		TelegramTail:      " Każda wiadomość jest wysyłana tylko ze zgodą KIFF na dokładnie to wywołanie. Możesz rozłączyć w każdej chwili, co usuwa webhook i token.",

		PausedF:          "Wstrzymane · %s · połącz ponownie, aby zatwierdzić",
		Connected:        "Połączono",
		NewToolsWaiting:  "Połączono · nowe narzędzia czekają na Twoje zatwierdzenie",
		ToolsWaiting:     "Zalogowano · narzędzia czekają na Twoje zatwierdzenie",
		WaitingForYou:    "Czeka na Twoje zatwierdzenie",
		DoorEvents:       "Połączono · zdarzenia drzwi włączone",
		ToolsBeingChosen: "Zalogowano · narzędzia są wybierane",
		SignedOutByF:     "Wylogowano przez %s · zaloguj się ponownie",
		ApprovedWaitF:    "Zatwierdzono · czeka na logowanie w %s",
		WaitingToConfirm: "Zalogowano · czeka na Twoje potwierdzenie",
		NotConnected:     "Nie połączono",
		Disconnected:     "Rozłączono",
	},
}
