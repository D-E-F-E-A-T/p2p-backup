package middleware

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

var (
	CryptoCurrencies = map[uint8]string{
		1: "btc",
		2: "eth",
		3: "wvs",
		4: "xrp",
		5: "ltc",
		6: "bch",
	}
	FiatCurrencies = map[uint8]string{
		7:  "usd",
		8:  "eur",
		9:  "rub",
		10: "gbp",
	}
	CurrencySign = map[uint8]string{
		7:  "$",
		8:  "€",
		9:  "₽",
		10: "£",
	}
	Providers = map[uint8]string{
		1:  "Cash",
		2:  "Sberbank",
		3:  "Qiwi",
		4:  "Yandex.Money",
		5:  "Alpha Bank",
		6:  "Tinkoff",
		7:  "VTB",
		8:  "Rosbank",
		9:  "MTS",
		10: "Beeline",
		11: "Tele2",
		12: "Megafon",
		13: "Rostelekon",
		14: "WebMoney",
		15: "Payeer",
	}
	ComplaintCategories = map[uint8]string{
		1: "user_gone",
		2: "incorrect_sum",
		3: "script_error",
	}
	Groups = map[int]string{
		1: "admin",
		2: "support",
	}
	Languages = []string{
		"cn",
		"de",
		"ee",
		"en",
		"es",
		"fr",
		"it",
		"lv",
		"ru",
	}
)

type HTMLRenderer struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

type HTMLInstance struct {
	Template *template.Template
}

func (r HTMLInstance) Instance(name string, data interface{}) render.Render {
	return HTMLRenderer{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

// Render (HTML) executes template and writes its result with custom ContentType for response.
func (r HTMLRenderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if h, ok := r.Data.(gin.H); ok {
		h["FiatCurrencies"] = FiatCurrencies
		h["CryptoCurrencies"] = CryptoCurrencies
		h["Providers"] = Providers
		h["ComplaintCategories"] = ComplaintCategories
		h["Groups"] = Groups
	}

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

func LoadRenderer(pattern string, funcMap template.FuncMap) HTMLInstance {
	templ := template.Must(template.New("").Funcs(funcMap).ParseGlob(pattern))

	return HTMLInstance{Template: templ}
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HTMLRenderer) WriteContentType(w http.ResponseWriter) {
	//writeContentType(w, htmlContentType)
}

func Client(c *gin.Context) {
	sess := sessions.Default(c)
	_, ok := sess.Get("Token").(string)
	c.Set("IsLogged", ok)

	var count int
	var notifications interface{}
	if user, ok := sess.Get("User").(*models.User); ok {
		numOffers, _ := postgres.Postgres.GetUserOffersAmount(user.ID)
		c.Set("NumOffers", numOffers)
		numDeals, _ := postgres.Postgres.GetUserDealsAmount(user.ID)
		c.Set("NumDeals", numDeals)
		count, notifications, _ = postgres.Postgres.SearchUserNotifications(user.ID, 10, 0)
		c.Set("NotificationsAmount", count)
	}
	c.Set("Notifications", notifications)
	c.Next()
}

func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

func FormatDate(t time.Time) string {
	return t.Format("01/02/2006")
}

func FormatFloat(i interface{}) string {
	switch i.(type) {
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 64)
	}
	return ""
}

// TODO rewrite validation

func isSliceContain(value interface{}, slice ...interface{}) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func isMapContain(value interface{}, map_ interface{}) bool {
	for _, v := range map_.(map[interface{}]interface{}) {
		if v == value {
			return true
		}
	}
	return false
}

func IsFiatUint(currency uint8) (y bool) {
	for id := range FiatCurrencies {
		if id == currency {
			return true
		}
	}
	return false
}

func IsFiatString(currency string) (y bool) {
	for _, id := range FiatCurrencies {
		if id == currency {
			return true
		}
	}
	return false
}

func IsCryptoUint(currency uint8) (y bool) {
	for id := range CryptoCurrencies {
		if id == currency {
			return true
		}
	}
	return false
}

func IsCryptoString(currency string) (y bool) {
	for _, id := range CryptoCurrencies {
		if id == currency {
			return true
		}
	}
	return false
}

func Currency(currency uint8) string {
	for key, value := range FiatCurrencies {
		if key == currency {
			return value
		}
	}
	for key, value := range CryptoCurrencies {
		if key == currency {
			return value
		}
	}
	return ""
}

func CurrencyID(currency string) uint8 {
	for key, value := range FiatCurrencies {
		if value == currency {
			return key
		}
	}
	for key, value := range CryptoCurrencies {
		if value == currency {
			return key
		}
	}
	return 0
}

func Provider(provider uint8) string {
	for key, value := range Providers {
		if key == provider {
			return value
		}
	}
	return ""
}

func IsActiveDay(days []time.Weekday, day int) bool {
	for _, d := range days {
		if d == time.Weekday(day) {
			return true
		}
	}
	return false
}

func Country(countries []*models.Country, code string) string {
	for _, country := range countries {
		if country.ID == code {
			return country.Value
		}
	}
	return ""
}

func IsOnline(t time.Time) bool {
	return t.After(time.Now().Truncate(time.Hour * 1))
}

func IsGroup(groups []int, group int) bool {
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}

func Sign(i interface{}) string {
	switch i.(type) {
	case string:
		i = CurrencyID(i.(string))
	}
	return CurrencySign[i.(uint8)]
}

func IsAdmin(values []int) bool {
	for _, value := range values {
		if value == 1 {
			return true
		}
	}
	return false
}

var Countries = []string{
	"au",
	"at",
	"az",
	"ax",
	"al",
	"dz",
	"as",
	"ai",
	"ao",
	"ad",
	"aq",
	"ag",
	"ar",
	"am",
	"aw",
	"af",
	"bs",
	"bd",
	"bb",
	"bh",
	"by",
	"bz",
	"be",
	"bj",
	"bm",
	"bg",
	"bo",
	"bq",
	"ba",
	"bw",
	"br",
	"io",
	"bn",
	"bf",
	"bi",
	"bt",
	"vu",
	"va",
	"gb",
	"hu",
	"ve",
	"vg",
	"vi",
	"um",
	"tl",
	"vn",
	"ga",
	"ht",
	"gy",
	"gm",
	"gh",
	"gp",
	"gt",
	"gn",
	"gw",
	"de",
	"gg",
	"gi",
	"hn",
	"hk",
	"gd",
	"gl",
	"gr",
	"ge",
	"gu",
	"dk",
	"je",
	"dj",
	"dg",
	"dm",
	"do",
	"eg",
	"zm",
	"eh",
	"zw",
	"il",
	"in",
	"id",
	"jo",
	"iq",
	"ir",
	"ie",
	"is",
	"es",
	"it",
	"ye",
	"cv",
	"kz",
	"kh",
	"cm",
	"ca",
	"ic",
	"qa",
	"ke",
	"cy",
	"kg",
	"ki",
	"cn",
	"kp",
	"cc",
	"co",
	"km",
	"cg",
	"cd",
	"xk",
	"cr",
	"ci",
	"cu",
	"kw",
	"cw",
	"la",
	"lv",
	"ls",
	"lr",
	"lb",
	"ly",
	"lt",
	"li",
	"lu",
	"mu",
	"mr",
	"mg",
	"yt",
	"mo",
	"mw",
	"my",
	"ml",
	"mv",
	"mt",
	"ma",
	"mq",
	"mh",
	"mx",
	"mz",
	"md",
	"mc",
	"mn",
	"ms",
	"mm",
	"na",
	"nr",
	"np",
	"ne",
	"ng",
	"nl",
	"ni",
	"nu",
	"nz",
	"nc",
	"no",
	"ac",
	"im",
	"nf",
	"cx",
	"sh",
	"pn",
	"tc",
	"ae",
	"om",
	"ky",
	"ck",
	"pk",
	"pw",
	"ps",
	"pa",
	"pg",
	"py",
	"pe",
	"pl",
	"pt",
	"xb",
	"xa",
	"pr",
	"kr",
	"re",
	"ru",
	"rw",
	"ro",
	"sv",
	"ws",
	"sm",
	"st",
	"sa",
	"mk",
	"mp",
	"sc",
	"bl",
	"mf",
	"pm",
	"sn",
	"vc",
	"kn",
	"lc",
	"rs",
	"ea",
	"sg",
	"sx",
	"sy",
	"sk",
	"si",
	"us",
	"sb",
	"so",
	"sd",
	"sr",
	"sl",
	"tj",
	"th",
	"tw",
	"tz",
	"tg",
	"tk",
	"to",
	"tt",
	"ta",
	"tv",
	"tn",
	"tm",
	"tr",
	"ug",
	"uz",
	"ua",
	"wf",
	"uy",
	"fo",
	"fm",
	"fj",
	"ph",
	"fi",
	"fk",
	"fr",
	"gf",
	"pf",
	"tf",
	"hr",
	"cf",
	"td",
	"me",
	"cz",
	"cl",
	"ch",
	"se",
	"sj",
	"lk",
	"ec",
	"gq",
	"er",
	"sz",
	"ee",
	"et",
	"gs",
	"za",
	"ss",
	"jm",
	"jp",
}
