package dictionary

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const apiKey = "dict.1.1.20140416T183822Z.3b90bf5bedccc85b.93d3bab6d7fb38c57e7fd1ebd1aa6442bb64876a"

func TestDictionaryAPI(t *testing.T) {
	dict := New(apiKey)

	Convey("#GetLangs", t, func() {
		Convey("On success it returns available languages", func() {
			langs, _ := dict.GetLangs()
			So(langs, ShouldNotBeEmpty)
			pairs := []string{"ru-ru", "ru-en", "ru-de", "ru-it", "ru-fr", "en-de", "en-it", "en-ru"}
			for _, pair := range pairs {
				So(langs, ShouldContain, pair)
			}
		})

		Convey("On failure it returns error code and message", func() {
			tr := New(apiKey + "a")
			response, err := tr.GetLangs()
			So(response, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "(401) API key is invalid")
		})
	})

	Convey("#Lookup", t, func() {
		Convey("On success it returns translation of the word", func() {
			Convey("with all possible fields filled in case of russian or english", func() {
				entry, err := dict.Lookup(&Params{Lang: "en-ru", Text: "dog"})

				So(err, ShouldBeNil)
				So(entry, ShouldNotBeNil)

				So(entry.Code, ShouldEqual, 0)
				So(entry.Message, ShouldBeBlank)

				So(entry.Def[0].Text, ShouldEqual, "dog")
				So(entry.Def[0].Pos, ShouldEqual, "noun")
				So(entry.Def[0].Ts, ShouldNotBeBlank)

				So(entry.Def[0].Tr[0].Text, ShouldEqual, "собака")
				So(entry.Def[0].Tr[0].Pos, ShouldEqual, "noun")

				So(entry.Def[0].Tr[0].Syn[0].Text, ShouldNotBeBlank)
				So(entry.Def[0].Tr[0].Mean[0].Text, ShouldNotBeBlank)
				So(entry.Def[0].Tr[0].Ex[0].Text, ShouldNotBeBlank)
				So(entry.Def[0].Tr[0].Ex[0].Tr[0].Text, ShouldNotBeBlank)
			})

			Convey("With some fields not fields in case of other languages", func() {
				entry, err := dict.Lookup(&Params{Lang: "en-de", Text: "dog"})

				So(err, ShouldBeNil)
				So(entry, ShouldNotBeNil)

				So(entry.Code, ShouldEqual, 0)
				So(entry.Message, ShouldBeBlank)

				So(entry.Def[0].Text, ShouldEqual, "dog")
				So(entry.Def[0].Pos, ShouldEqual, "noun")
				So(entry.Def[0].Ts, ShouldNotBeBlank)

				So(entry.Def[0].Tr[0].Text, ShouldEqual, "Hund")
				So(entry.Def[0].Tr[0].Pos, ShouldEqual, "noun")

				So(entry.Def[0].Tr[0].Syn, ShouldBeEmpty)
				So(entry.Def[0].Tr[0].Mean, ShouldNotBeEmpty)
				So(entry.Def[0].Tr[0].Ex, ShouldBeEmpty)
			})

			Convey("Using different language for the interface", func() {
				dict = NewUsingLang(apiKey, "ru")
				entry, _ := dict.Lookup(&Params{Lang: "en-ru", Text: "dog"})

				So(entry.Def[0].Text, ShouldEqual, "dog")
				So(entry.Def[0].Pos, ShouldEqual, "существительное")
				So(entry.Def[0].Ts, ShouldNotBeBlank)

				So(entry.Def[0].Tr[0].Text, ShouldEqual, "собака")
				So(entry.Def[0].Tr[0].Pos, ShouldEqual, "существительное")
			})
		})

		Convey("On failure it returns error", func() {
			entry, err := dict.Lookup(&Params{Lang: "en-mumbayumba", Text: "dog"})

			So(err, ShouldNotBeNil)
			So(entry, ShouldBeNil)
		})
	})

}
