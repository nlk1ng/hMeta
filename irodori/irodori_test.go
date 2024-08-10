package irodori

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	hmeta "github.com/nlk1ng/hMeta"
)

func TestIrodori(t *testing.T) {
	var ua string
	if v, hv := os.LookupEnv("USER_AGENT"); hv {
		ua = v
	} else {
		panic("missing user agent")
	}
	testCases := []struct {
		desc string
		url  string
	}{
		{
			desc: "1",
			url:  "https://irodoricomics.com/Screwed-by-Step-Dad-All-About-Yui-2",
		}, {
			"2",
			"https://irodoricomics.com/Top-Class-MILF",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			d, err := ByURL(tC.url, hmeta.SetUserAgent(ua))
			if err != nil {
				t.Error(err.Error())
			}
			t.Logf("%+v\n", d)
		})
	}
}
