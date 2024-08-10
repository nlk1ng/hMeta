package hentag

import hmeta "github.com/nlk1ng/hMeta"

type Metadata struct {
	hmeta.Metadata
	Character    []string
	ThumbnailUrl string
	Language     string
	Category     string
}
