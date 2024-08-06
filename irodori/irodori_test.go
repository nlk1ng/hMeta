package irodori

import "testing"

func TestIrodori(t *testing.T) {
	testCases := []struct {
		desc string
		url  string
	}{
		{
			desc: "",
			url:  "https://irodoricomics.com/Screwed-by-Step-Dad-All-About-Yui-2",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			d, err := Irodori(tC.url)
			if err != nil {
				t.Error(err.Error())
			}
			t.Logf("%+v", d)
		})
	}
}
