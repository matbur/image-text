package i18n_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matbur/image-text/i18n"
)

func TestFlag(t *testing.T) {
	assert.Equal(t, "🇵🇱", i18n.Flag(i18n.PL))
	assert.Equal(t, "🇪🇸", i18n.Flag(i18n.ES))
	assert.Equal(t, "", i18n.Flag("fr"))
}

func TestNativeLabel(t *testing.T) {
	assert.Equal(t, "English", i18n.NativeLabel(i18n.EN))
	assert.Equal(t, "Polski", i18n.NativeLabel(i18n.PL))
	assert.Equal(t, "Español", i18n.NativeLabel(i18n.ES))

	locale := i18n.Locale{Tag: i18n.PL}
	assert.Equal(t, "English", i18n.NativeLabel(i18n.EN))
	assert.Equal(t, "Polski", i18n.NativeLabel(i18n.PL))
	assert.Equal(t, locale.T("lang.en"), "Angielski")
}

func TestLocaleT(t *testing.T) {
	locale := i18n.Locale{Tag: i18n.PL}

	assert.Equal(t, "Ustawienia", locale.T("common.settings"))
	assert.Equal(t, "common.missing", locale.T("common.missing"))
}

func TestFromRequestCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: i18n.CookieName, Value: "es"})

	locale := i18n.FromRequest(req)
	assert.Equal(t, i18n.ES, locale.Tag)
	assert.Equal(t, "Ajustes", locale.T("common.settings"))
}

func TestFromRequestAcceptLanguage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en;q=0.8")

	locale := i18n.FromRequest(req)
	assert.Equal(t, i18n.PL, locale.Tag)
}

func TestFromRequestDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	locale := i18n.FromRequest(req)
	assert.Equal(t, i18n.Default, locale.Tag)
}

func TestAllLocalesHaveSameKeys(t *testing.T) {
	base := i18n.Locale{Tag: i18n.EN}
	pl := i18n.Locale{Tag: i18n.PL}
	es := i18n.Locale{Tag: i18n.ES}

	keys := []string{
		"index.lead",
		"common.settings",
		"docs.title",
		"lang.pl",
	}

	for _, key := range keys {
		require.NotEqual(t, key, base.T(key))
		require.NotEqual(t, key, pl.T(key))
		require.NotEqual(t, key, es.T(key))
		assert.NotEqual(t, base.T(key), pl.T(key))
		assert.NotEqual(t, base.T(key), es.T(key))
	}
}
