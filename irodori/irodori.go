package irodori

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

const irodoriTitleSelector = `#product > div.title.page-title`
const irodoriDescriptionSelector = "#product > div.product-blocks.blocks-top > div.product_extra.product_blocks.product_blocks-top > div > div > div > div"
const irodoriTumbSelector = `#content > div.product-info.has-extra-button > div.product-left > div.lightgallery.lightgallery-product-images`
const irodoriArtistSelector = `#product > div.product-stats > ul > li.product-manufacturer > a`
const irodoriPagesSelector = `#product > div.product-stats > ul > li.product-upc > span`

type irodoriThumbJson struct {
	Src     string `json:"src"`
	Thumb   string `json:"thumb"`
	SubHtml string `json:"subHtml"`
}

func Irodori(url string) (Metadata, error) {
	var gal Metadata
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	// Title
	collector.OnHTML(irodoriTitleSelector, func(h *colly.HTMLElement) {
		gal.Title = h.DOM.Text()
	})

	// Description
	collector.OnHTML(irodoriDescriptionSelector,
		func(h *colly.HTMLElement) {
			txts := make([]string, 0)
			h.ForEach("p", func(i int, h *colly.HTMLElement) {
				txts = append(txts, h.Text)
			})

			gal.Description = strings.Join(txts, "\n")
		},
	)

	// Thumbnail
	collector.OnHTML(irodoriTumbSelector, func(h *colly.HTMLElement) {

		imgJsn := h.Attr("data-images")
		if imgJsn == "" {
			return
		}
		imgDta := []irodoriThumbJson{}

		err := json.Unmarshal([]byte(imgJsn), &imgDta)
		if err != nil {
			println(err.Error())
			return
		}

		for _, v := range imgDta {
			gal.ThumbnailUrl = v.Src
			break
		}
	})

	// Artist
	collector.OnHTML(irodoriArtistSelector, func(h *colly.HTMLElement) {
		gal.Artist = append(gal.Artist, h.DOM.Text())
	})

	// Pages
	collector.OnHTML(irodoriPagesSelector, func(h *colly.HTMLElement) {
		p, err := strconv.Atoi(h.DOM.Text())
		if err == nil {
			gal.Pages = p
		}
	})

	// Tags
	collector.OnHTML(`body > ul.ctagList`, func(h *colly.HTMLElement) {
		h.ForEach("li > span > a", func(i int, h *colly.HTMLElement) {
			gal.Tag = append(gal.Tag, h.Text)
		})
	})

	// Product ID
	collector.OnHTML("html", func(h *colly.HTMLElement) {
		class, exst := h.DOM.Attr("class")
		if !exst {
			return
		}
		trgetClass := regexp.MustCompile(`product-\d+`).FindString(class)
		prodId, err := parseIrodoriProductIdClass(trgetClass)
		if err != nil {
			return
		}
		h.Request.Visit(fmt.Sprintf("https://irodoricomics.com/index.php?route=product/product/cattags&product_id=%v", prodId))
	})

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", "_rps=1;")
	})

	err := collector.Visit(url)
	if err != nil {
		return gal, err
	}
	return gal, nil
}

func parseIrodoriProductIdClass(class string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(class, "product-"))
}
