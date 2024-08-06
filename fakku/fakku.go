package fakku

import (
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type FakkuScraperConfig struct {
	CfClearance string
	UserAgent   string
	FakkuSid    string
	Delay       time.Duration
}

const fakkuMainContentSelector = `body > div.block.max-w-screen-xl.relative.mx-auto.flex.flex-col.h-screen > div.grid.grid-flow-row-dense.grid-cols-2.sm\:grid-cols-4.md\:grid-cols-6.lg\:grid-cols-10.gap-6.md\:gap-8.relative.text-center.w-full.px-6.pt-20.mx-auto.xl\:px-0 > div:nth-child(2) > div > div.block.md\:table-cell.relative.w-full.align-top.p-8.pt-0.md\:pt-8.md\:pl-2.space-y-4`
const fakkuThumbSelector = `body > div.block.max-w-screen-xl.relative.mx-auto.flex.flex-col.h-screen > div.grid.grid-flow-row-dense.grid-cols-2.sm\:grid-cols-4.md\:grid-cols-6.lg\:grid-cols-10.gap-6.md\:gap-8.relative.text-center.w-full.px-6.pt-20.mx-auto.xl\:px-0 > div:nth-child(2) > div > div.block.sm\:inline-block.relative.w-full.align-top.p-2.md\:p-8.text-center.space-y-4 > div > a > img`
const fakkuDescriptionSelector = `head > meta:nth-child(17)`

const otherContentContainer = `body > div.block.max-w-screen-xl.relative.mx-auto.flex.flex-col.h-screen > div.grid.grid-flow-row-dense.grid-cols-2.sm\:grid-cols-4.md\:grid-cols-6.lg\:grid-cols-10.gap-6.md\:gap-8.relative.text-center.w-full.px-6.pt-20.mx-auto.xl\:px-0 > div.col-span-full.grid.grid-flow-row-dense.grid-cols-10.gap-x-4.gap-y-4.sm\:mx-24`

func Fakku(config FakkuScraperConfig, urls ...string) (Metadata, error) {
	var gal Metadata
	collector := colly.NewCollector(
		colly.UserAgent(config.UserAgent),
	)

	collector.Limit(&colly.LimitRule{
		DomainRegexp: "https://www.fakku.net",
		RandomDelay:  time.Millisecond * 15,
		Delay:        config.Delay,
	})

	// Set Cookies
	collector.OnRequest(func(r *colly.Request) {
		var rawCookie strings.Builder
		rawCookie.WriteString("cf_clearance=" + config.CfClearance + ";")
		if config.FakkuSid != "" {
			rawCookie.WriteString("fakku_sid=" + config.FakkuSid + ";")
		}
		r.Headers.Set("Cookie", rawCookie.String())
	})

	// Description
	collector.OnHTML(fakkuDescriptionSelector,
		func(h *colly.HTMLElement) {
			gal.Description = html.UnescapeString(h.Attr("content"))
		},
	)

	// Thumbnail
	collector.OnHTML(fakkuThumbSelector, func(h *colly.HTMLElement) {
		att, exst := h.DOM.Attr("src")
		if exst {
			gal.ThumbnailUrl = att
		}
	})

	// Main Content
	collector.OnHTML(fakkuMainContentSelector, func(h *colly.HTMLElement) {
		// Title
		titleElem := h.DOM.Find(`h1`)
		gal.Title = titleElem.Text()

		// Tags
		h.DOM.Children().Filter(`div.table`).Each(func(i int, s *goquery.Selection) {
			// gal.Tag = append(gal.Tag, s.Text())

			switch s.Find("div:nth-child(1)").Text() {
			case "Artist":
				gal.Artist = append(gal.Artist, parseTags(s)...)
			case "Circle":
				gal.Circle = append(gal.Circle, parseTags(s)...)
			case "Parody":
				gal.Parody = append(gal.Parody, parseTags(s)...)
			case "Magazine":
				gal.Magazine = append(gal.Magazine, parseTags(s)...)
			case "Publisher":
				gal.Publisher = append(gal.Publisher, parseTags(s)...)
			case "Book":
				gal.Book = append(gal.Book, parseTags(s)...)
			case "Pages":
				p, err := strconv.Atoi(strings.TrimSuffix(s.Find("div:nth-child(2)").Text(), " pages"))
				if err == nil {
					gal.Pages = p
				}
			case "Favorites":
				p, err := strconv.Atoi(strings.TrimSuffix(s.Find("div:nth-child(2)").Text(), " favorites"))
				if err == nil {
					gal.Favorites = p
				}
			default:
				if s.Find("div:nth-child(1) > a").Text() != "" {
					s.Find("div:nth-child(1) > a").Each(func(i int, s *goquery.Selection) {
						gal.Tag = append(gal.Tag, strings.TrimSpace(s.Text()))
					})
				}
			}
		})
	})

	// Collection
	collector.OnHTML(otherContentContainer,
		func(h *colly.HTMLElement) {
			// Check if the second tab is named "Collection"
			if h.DOM.Find(`div.col-span-full.border-b.border-gray-300.overflow-hidden.js-tab-container > ul > li:nth-child(2) > div > a`).Text() == "Collections" {
				aElm := h.DOM.Find(`div.col-span-full.block.js-tab-targets>div:nth-child(2)>div:nth-child(1)>h2>em>a`)
				if hrf, hrfExst := aElm.Attr("href"); hrfExst {
					gal.Collection = append(gal.Collection, Collection{aElm.Text(), "https://www.fakku.net" + hrf})
				}
			}
		},
	)

	for _, url := range urls {
		err := collector.Visit(url)
		if err != nil {
			return gal, err
		}
	}
	collector.Wait()
	return gal, nil
}

func parseTags(s *goquery.Selection) (tags []string) {
	s.Find("div:nth-child(2)").Children().Each(func(i int, s *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(s.Text()))
	})

	return
}
