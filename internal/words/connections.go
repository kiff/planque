package words

import (
	"fmt"
	"time"
)

// Connections is one language's words for a connection as the owner's list
// shows it: its health in plain words, and the actions that can work for it,
// so every app draws the same line and the same buttons, and none keeps a
// rule of its own. The terms follow the approval words above.
type Connections struct {
	// Actions.
	SignInAgain      string
	Remove           string
	RemoveUnfinished string
	StopsF           string // the taf, the app: what stops when access is removed
	NothingStops     string

	// Health.
	WorkingLastF   string // when it was last used
	WorkingNotYet  string
	Working        string // when planque didn't look up its last use
	ExpiredOnF     string // the app, the day
	ExpiredF       string // the app
	UnfinishedOnF  string // the day it was signed in
	Unfinished     string
	Months         [12]string
	DayMonthF      string // the day, the month's short name
	TimeOnlyLayout string // a time on the day it is, as Go's time layout
}

var connections = map[string]Connections{
	"en": {
		SignInAgain: "Sign in again", Remove: "Remove access", RemoveUnfinished: "Remove unfinished connection",
		StopsF: "%s will stop using %s.", NothingStops: "No taf uses it, so nothing stops.",
		WorkingLastF: "Working · last used %s", WorkingNotYet: "Working · not used yet", Working: "Working",
		ExpiredOnF: "%s's sign-in expired on %s", ExpiredF: "%s's sign-in has expired",
		UnfinishedOnF: "Unfinished · signed in on %s, no tools chosen, no taf uses it", Unfinished: "Unfinished · no tools chosen, no taf uses it",
		Months:    [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		DayMonthF: "%d %s", TimeOnlyLayout: "15:04",
	},
	"de": {
		SignInAgain: "Erneut anmelden", Remove: "Zugriff entfernen", RemoveUnfinished: "Unvollständige Verbindung entfernen",
		StopsF: "%s nutzt %s dann nicht mehr.", NothingStops: "Kein taf nutzt sie, also hört nichts auf.",
		WorkingLastF: "Funktioniert · zuletzt genutzt %s", WorkingNotYet: "Funktioniert · noch nicht genutzt", Working: "Funktioniert",
		ExpiredOnF: "Die Anmeldung bei %s ist am %s abgelaufen", ExpiredF: "Die Anmeldung bei %s ist abgelaufen",
		UnfinishedOnF: "Unvollständig · angemeldet am %s, keine Tools gewählt, kein taf nutzt sie", Unfinished: "Unvollständig · keine Tools gewählt, kein taf nutzt sie",
		Months:    [12]string{"Jan.", "Feb.", "März", "Apr.", "Mai", "Juni", "Juli", "Aug.", "Sept.", "Okt.", "Nov.", "Dez."},
		DayMonthF: "%d. %s", TimeOnlyLayout: "15:04",
	},
	"fr": {
		SignInAgain: "Se reconnecter", Remove: "Retirer l'accès", RemoveUnfinished: "Retirer la connexion inachevée",
		StopsF: "%s n'utilisera plus %s.", NothingStops: "Aucun taf ne l'utilise, donc rien ne s'arrête.",
		WorkingLastF: "Fonctionne · dernière utilisation : %s", WorkingNotYet: "Fonctionne · pas encore utilisée", Working: "Fonctionne",
		ExpiredOnF: "La connexion à %s a expiré le %s", ExpiredF: "La connexion à %s a expiré",
		UnfinishedOnF: "Inachevée · connexion faite le %s, aucun outil choisi, aucun taf ne l'utilise", Unfinished: "Inachevée · aucun outil choisi, aucun taf ne l'utilise",
		Months:    [12]string{"janv.", "févr.", "mars", "avr.", "mai", "juin", "juil.", "août", "sept.", "oct.", "nov.", "déc."},
		DayMonthF: "%d %s", TimeOnlyLayout: "15:04",
	},
	"es": {
		SignInAgain: "Volver a iniciar sesión", Remove: "Quitar el acceso", RemoveUnfinished: "Quitar la conexión sin terminar",
		StopsF: "%s dejará de usar %s.", NothingStops: "Ningún taf la usa, así que nada se detiene.",
		WorkingLastF: "Funciona · último uso: %s", WorkingNotYet: "Funciona · aún sin usar", Working: "Funciona",
		ExpiredOnF: "La sesión de %s caducó el %s", ExpiredF: "La sesión de %s ha caducado",
		UnfinishedOnF: "Sin terminar · sesión iniciada el %s, sin herramientas elegidas, ningún taf la usa", Unfinished: "Sin terminar · sin herramientas elegidas, ningún taf la usa",
		Months:    [12]string{"ene.", "feb.", "mar.", "abr.", "may.", "jun.", "jul.", "ago.", "sept.", "oct.", "nov.", "dic."},
		DayMonthF: "%d %s", TimeOnlyLayout: "15:04",
	},
	"it": {
		SignInAgain: "Accedi di nuovo", Remove: "Rimuovi l'accesso", RemoveUnfinished: "Rimuovi la connessione incompleta",
		StopsF: "%s smetterà di usare %s.", NothingStops: "Nessun taf la usa, quindi non si ferma nulla.",
		WorkingLastF: "Funziona · ultimo uso: %s", WorkingNotYet: "Funziona · non ancora usata", Working: "Funziona",
		ExpiredOnF: "L'accesso a %s è scaduto il %s", ExpiredF: "L'accesso a %s è scaduto",
		UnfinishedOnF: "Incompleta · accesso effettuato il %s, nessuno strumento scelto, nessun taf la usa", Unfinished: "Incompleta · nessuno strumento scelto, nessun taf la usa",
		Months:    [12]string{"gen", "feb", "mar", "apr", "mag", "giu", "lug", "ago", "set", "ott", "nov", "dic"},
		DayMonthF: "%d %s", TimeOnlyLayout: "15:04",
	},
	"pt": {
		SignInAgain: "Iniciar sessão novamente", Remove: "Remover o acesso", RemoveUnfinished: "Remover a ligação por concluir",
		StopsF: "%s deixará de usar %s.", NothingStops: "Nenhum taf a usa, por isso nada para.",
		WorkingLastF: "A funcionar · última utilização: %s", WorkingNotYet: "A funcionar · ainda não utilizada", Working: "A funcionar",
		ExpiredOnF: "A sessão em %s expirou a %s", ExpiredF: "A sessão em %s expirou",
		UnfinishedOnF: "Por concluir · sessão iniciada a %s, sem ferramentas escolhidas, nenhum taf a usa", Unfinished: "Por concluir · sem ferramentas escolhidas, nenhum taf a usa",
		Months:    [12]string{"jan.", "fev.", "mar.", "abr.", "mai.", "jun.", "jul.", "ago.", "set.", "out.", "nov.", "dez."},
		DayMonthF: "%d de %s", TimeOnlyLayout: "15:04",
	},
	"nl": {
		SignInAgain: "Opnieuw inloggen", Remove: "Toegang intrekken", RemoveUnfinished: "Onvoltooide koppeling verwijderen",
		StopsF: "%s gebruikt %s dan niet meer.", NothingStops: "Geen taf gebruikt deze koppeling, dus er stopt niets.",
		WorkingLastF: "Werkt · laatst gebruikt %s", WorkingNotYet: "Werkt · nog niet gebruikt", Working: "Werkt",
		ExpiredOnF: "De inlog bij %s is verlopen op %s", ExpiredF: "De inlog bij %s is verlopen",
		UnfinishedOnF: "Onvoltooid · ingelogd op %s, geen tools gekozen, geen taf gebruikt deze koppeling", Unfinished: "Onvoltooid · geen tools gekozen, geen taf gebruikt deze koppeling",
		Months:    [12]string{"jan.", "feb.", "mrt.", "apr.", "mei", "jun.", "jul.", "aug.", "sep.", "okt.", "nov.", "dec."},
		DayMonthF: "%d %s", TimeOnlyLayout: "15:04",
	},
	"pl": {
		SignInAgain: "Zaloguj się ponownie", Remove: "Usuń dostęp", RemoveUnfinished: "Usuń niedokończone połączenie",
		StopsF: "%s przestanie korzystać z: %s.", NothingStops: "Żaden taf go nie używa, więc nic się nie zatrzyma.",
		WorkingLastF: "Działa · ostatnio użyte: %s", WorkingNotYet: "Działa · jeszcze nieużyte", Working: "Działa",
		ExpiredOnF: "Logowanie w %s wygasło %s", ExpiredF: "Logowanie w %s wygasło",
		UnfinishedOnF: "Niedokończone · zalogowano %s, nie wybrano narzędzi, żaden taf go nie używa", Unfinished: "Niedokończone · nie wybrano narzędzi, żaden taf go nie używa",
		Months:    [12]string{"sty", "lut", "mar", "kwi", "maj", "cze", "lip", "sie", "wrz", "paź", "lis", "gru"},
		DayMonthF: "%d %s", TimeOnlyLayout: "15:04",
	},
}

// Connection is lang's connection words, English's for a language planque
// doesn't speak.
func Connection(lang string) Connections {
	if c, ok := connections[lang]; ok {
		return c
	}
	return connections["en"]
}

// Day is t's day as the owner reads it ("14 Nov", "14. Nov.", "14 de nov.").
func (c Connections) Day(t time.Time) string {
	return fmt.Sprintf(c.DayMonthF, t.Day(), c.Months[t.Month()-1])
}

// When is t as a time when it falls on now's day, else as its day: what a
// last use says. Both are in the zone they come in.
func (c Connections) When(t, now time.Time) string {
	if y, m, d := t.Date(); func() bool { ny, nm, nd := now.Date(); return y == ny && m == nm && d == nd }() {
		return t.Format(c.TimeOnlyLayout)
	}
	return c.Day(t)
}
