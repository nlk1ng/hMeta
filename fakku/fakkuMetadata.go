package fakku

import hmeta "github.com/nlk1ng/hMeta"

type Metadata struct {
	hmeta.Metadata
	Description  string
	Magazine     []string
	Publisher    []string
	Book         []string
	Favorites    int
	Pages        int
	ThumbnailUrl string
	Collection   []Collection
}

type Collection struct {
	Name string
	Link string
}

func (c Collection) String() string {
	return c.Name
}
