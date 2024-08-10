package hentag

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	urlEndpoint   = "/url"
	titleEndpoint = "/title"
	idEndpoint    = "/title"
)

type HentagResponse struct {
	Title         string   `json:"title"`
	CoverImageUrl *string  `json:"coverImageUrl,omitempty"`
	Parodies      []string `json:"parodies,omitempty"`
	Circles       []string `json:"circles,omitempty"`
	Artists       []string `json:"artists,omitempty"`
	Characters    []string `json:"characters,omitempty"`
	MaleTags      []string `json:"maleTags,omitempty"`
	FemaleTags    []string `json:"femaleTags,omitempty"`
	OtherTags     []string `json:"otherTags,omitempty"`
	Language      string   `json:"language"`
	Category      string   `json:"category"`
	CreatedAt     uint64   `json:"createdAt"`
	LastModified  uint64   `json:"lastModified"`
	PublishedOn   *uint64  `json:"publishedOn,omitempty"`
	Locations     []string `json:"locations"`
}

func (h HentagResponse) ToHentagMetadata() Metadata {
	var data Metadata
	data.Title = h.Title
	data.Artist = h.Artists
	data.Circle = h.Circles
	data.Parody = h.Parodies
	data.Character = h.Characters
	data.ThumbnailUrl = *h.CoverImageUrl
	data.Language = h.Language
	data.Category = h.Category
	for _, t := range h.FemaleTags {
		data.Tag = append(data.Tag, "female:"+t)
	}
	for _, t := range h.MaleTags {
		data.Tag = append(data.Tag, "male:"+t)
	}
	for _, t := range h.OtherTags {
		data.Tag = append(data.Tag, "other:"+t)
	}
	return data
}

type urlApiPayload struct {
	Urls     []string `json:"urls"`
	Language string   `json:"language,omitempty"`
}

type idApiPayload struct {
	Ids      []string `json:"ids"`
	Language string   `json:"language,omitempty"`
}

type titleApiPayload struct {
	Title    string `json:"title"`
	Language string `json:"language,omitempty"`
}

// SearchByURLWithContext will call the urls endpoint of Hentag. If language is an empty string, it will be omitted.
// Visit https://hentag.com/api-documentation for more information
func SearchByURLWithContext(ctx context.Context, urls []string, language string) ([]HentagResponse, error) {
	payload := urlApiPayload{Urls: urls, Language: language}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return requestSearchApi(ctx, urlEndpoint, bytes.NewReader(jsonPayload))
}

// SearchByURL will call SearchByURLWithContext using a context.Background
func SearchByURL(urls []string, language string) ([]HentagResponse, error) {
	return SearchByURLWithContext(context.Background(), urls, language)
}

// SearchByIDWithContext will call the ids (vault id) endpoint of Hentag. If language is an empty string, it will be omitted.
// Visit https://hentag.com/api-documentation for more information
func SearchByIDWithContext(ctx context.Context, ids []string, language string) ([]HentagResponse, error) {
	payload := idApiPayload{Ids: ids, Language: language}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return requestSearchApi(ctx, idEndpoint, bytes.NewReader(jsonPayload))
}

// SearchByURL will call SearchByIDWithContext using a context.Background
func SearchByID(ids []string, language string) ([]HentagResponse, error) {
	return SearchByIDWithContext(context.Background(), ids, language)
}

// SearchByTitleWithContext will call the title endpoint of Hentag. If language is an empty string, it will be omitted.
// Visit https://hentag.com/api-documentation for more information
func SearchByTitleWithContext(ctx context.Context, title string, language string) ([]HentagResponse, error) {
	payload := titleApiPayload{Title: title, Language: language}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return requestSearchApi(ctx, titleEndpoint, bytes.NewReader(jsonPayload))
}

// SearchByURL will call SearchByTitleWithContext using a context.Background
func SearchByTitle(title string, language string) ([]HentagResponse, error) {
	return SearchByTitleWithContext(context.Background(), title, language)
}

func requestSearchApi(ctx context.Context, endpoint string, body io.Reader) ([]HentagResponse, error) {
	var resp []HentagResponse
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://hentag.com/api/v1/search/vault"+endpoint,
		body,
	)
	if err != nil {
		return nil, err
	}

	rq.Header.Set("Content-Type", "application/json")
	clnt := http.Client{}
	res, err := clnt.Do(rq)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
