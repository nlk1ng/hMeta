package fakku

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestFakkuScraper(t *testing.T) {
	var fakkucf, fakkusid, ua string
	if v, hv := os.LookupEnv("FAKKU_CF"); hv {
		fakkucf = v
	}

	if v, hv := os.LookupEnv("FAKKU_SID"); hv {
		fakkusid = v
	}

	if v, hv := os.LookupEnv("USER_AGENT"); hv {
		ua = v
	}
	if fakkucf == "" || ua == "" {
		panic("missing fakku cloudflare tokken or user agent")
	}
	testCases := []struct {
		desc string
		url  string
	}{
		{
			desc: "normal",
			url:  "https://www.fakku.net/hentai/x-eros-girls-collection-113-henkuma-english",
		},
		{
			desc: "Multiple artist",
			url:  "https://www.fakku.net/hentai/black-tights-english",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			gal, err := Fakku(FakkuScraperConfig{
				UserAgent:   ua,
				CfClearance: fakkucf,
				FakkuSid:    fakkusid,
			}, tC.url)

			if err != nil {
				t.Error(err.Error())
			}

			t.Logf("%+v\n", gal)
		})
	}
}
