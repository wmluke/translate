package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	a "github.com/stretchr/testify/assert"
	goproperties "github.com/dmotylev/goproperties"
	"fmt"
	"encoding/json"
)

func (r TranslationResponse) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func TestEscapeNonAscii(t *testing.T) {
	assert := a.New(t)
	assert.Equal("El gato salt\\u00f3 sobre el sombrero", escapeNonAscii("El gato saltó sobre el sombrero"))
	assert.Equal("\\u732b\\u306f\\u5e3d\\u5b50\\u3092\\u98db\\u3073\\u8d8a\\u3048\\u305f", escapeNonAscii("猫は帽子を飛び越えた"))
}

func TestKeys(t *testing.T) {
	p := make(goproperties.Properties)

	p["aa.bb.cc"] = "foo"
	p["a.b.c"] = "foo"
	p["aaa.bbb.ccc"] = "foo"

	k := keys(p);

	assert := a.New(t)

	assert.Equal("a.b.c", k[0])
	assert.Equal("aa.bb.cc", k[1])
	assert.Equal("aaa.bbb.ccc", k[2])
}

func TestTranslate(t *testing.T) {
	params := struct {
			key    string
			q      string
			source string
			target string
		}{}

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		params.key = r.URL.Query()["key"][0]
		params.q = r.URL.Query()["q"][0]
		params.source = r.URL.Query()["source"][0]
		params.target = r.URL.Query()["target"][0]

		translation := Translation{TranslatedText: "El gato saltó sobre el sombrero"}
		translations := []Translation{translation}
		data := TranslationData{ Translations: translations}
		resp := TranslationResponse{ Data: data}

		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprint(rw, resp)
	}))

	defer ts.Close()

	apiUrl = ts.URL

	translation, err := translate("abc", "The cat jumped over the hat", "en", "es")

	assert := a.New(t)

	assert.Empty(err)
	assert.Equal("abc", params.key)
	assert.Equal("The cat jumped over the hat", params.q)
	assert.Equal("en", params.source)
	assert.Equal("es", params.target)
	assert.Equal("El gato saltó sobre el sombrero", translation)
}



