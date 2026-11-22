package words

// What each language added after the first eight must say on a managed app's
// approval card: the words for "our sign-in partner" (a pattern, one for the
// language) and for "because" (a piece of text in the "asks for more access"
// line). The first eight are in approval_test.go.
var morePartner = map[string]string{
	"ja": `(?i)サインインパートナー`,
	"ko": `(?i)로그인 파트너`,
	"zh": `(?i)登入合作夥伴`,
	"th": `พันธมิตรด้านการเข้าสู่ระบบ`,
	"vi": `(?i)đối tác đăng nhập`,
	"id": `(?i)mitra masuk`,
	"tr": `oturum açma ortağımız`,
	"lt": `prisijungimo partner`,
	"lv": `pierakstīšanās partner`,
	"et": `sisselogimispartner`,
	"hu": `bejelentkezési partner`,
	"sl": `partner\w* za prijavo`,
	"sk": `partner\w* pre prihlásenie`,
	"hr": `partner\w* za prijavu`,
	"cs": `partner\w* pro přihlášení`,
	"ro": `(?i)partener\w* nostru de conectare`,
	"bg": `партньор за вход`,
}

var moreBecause = map[string]string{
	"ja": "ため",
	"ko": "때문",
	"zh": "因為",
	"th": "เพราะ",
	"vi": " vì ",
	"id": " karena ",
	"tr": " çünkü ",
	"lt": " nes ",
	"lv": " jo ",
	"et": " sest ",
	"hu": " mert ",
	"sl": " ker ",
	"sk": " pretože ",
	"hr": " jer ",
	"cs": " protože ",
	"ro": " deoarece ",
	"bg": " защото ",
}
