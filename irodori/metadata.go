package irodori

import hmeta "github.com/nlk1ng/hMeta"

type Metadata struct {
	hmeta.Metadata
	Description  string
	Pages        int
	ThumbnailUrl string
}
