package hentag

import "errors"

var (
	ErrNotFound = errors.New("no entry found")
)

// ByURL will search Hentag using the url of original source. This will only return the first one if multiple metadatas are found.
// language is omitted if it's empty string
func ByURL(url string, language string) (Metadata, error) {
	resp, err := SearchByURL([]string{url}, language)
	if err != nil {
		return Metadata{}, err
	}
	if len(resp) <= 1 {
		return Metadata{}, ErrNotFound
	}
	return resp[0].ToHentagMetadata(), nil
}
