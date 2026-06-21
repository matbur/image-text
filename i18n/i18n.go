package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//go:embed locales/*.json
var localeFS embed.FS

const (
	EN         = "en"
	PL         = "pl"
	ES         = "es"
	Default    = EN
	CookieName = "lang"
)

var (
	Supported = []string{EN, PL, ES}
	flags     = map[string]string{
		EN: "🇺🇸",
		PL: "🇵🇱",
		ES: "🇪🇸",
	}
	catalogs map[string]map[string]string
)

func init() {
	catalogs = make(map[string]map[string]string, len(Supported))

	for _, tag := range Supported {
		bb, err := localeFS.ReadFile(fmt.Sprintf("locales/%s.json", tag))
		if err != nil {
			panic(fmt.Sprintf("i18n: load locale %q: %v", tag, err))
		}

		var messages map[string]string
		if err := json.Unmarshal(bb, &messages); err != nil {
			panic(fmt.Sprintf("i18n: parse locale %q: %v", tag, err))
		}

		catalogs[tag] = messages
	}
}

type Locale struct {
	Tag string
}

func (l Locale) T(key string) string {
	if messages, ok := catalogs[l.Tag]; ok {
		if value, ok := messages[key]; ok {
			return value
		}
	}

	if value, ok := catalogs[Default][key]; ok {
		return value
	}

	return key
}

func NativeLabel(tag string) string {
	if messages, ok := catalogs[tag]; ok {
		if value, ok := messages["lang."+tag]; ok {
			return value
		}
	}

	if value, ok := catalogs[Default]["lang."+tag]; ok {
		return value
	}

	return tag
}

func Flag(tag string) string {
	if flag, ok := flags[tag]; ok {
		return flag
	}
	return ""
}

func (l Locale) Flag() string {
	return Flag(l.Tag)
}

func FromRequest(r *http.Request) Locale {
	if cookie, err := r.Cookie(CookieName); err == nil {
		if tag := normalizeTag(cookie.Value); tag != "" {
			return Locale{Tag: tag}
		}
	}

	return Locale{Tag: matchAcceptLanguage(r.Header.Get("Accept-Language"))}
}

func normalizeTag(raw string) string {
	tag := strings.ToLower(strings.TrimSpace(raw))
	base, _, _ := strings.Cut(tag, "-")

	switch base {
	case EN, PL, ES:
		return base
	default:
		return ""
	}
}

func matchAcceptLanguage(header string) string {
	if header == "" {
		return Default
	}

	for _, part := range strings.Split(header, ",") {
		tag := strings.TrimSpace(strings.Split(part, ";")[0])
		if normalized := normalizeTag(tag); normalized != "" {
			return normalized
		}
	}

	return Default
}
