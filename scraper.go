package hmeta

import (
	"time"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	*colly.Collector
}
type ScraperOption func(*Scraper)

func SetRequestTimeout(timeout time.Duration) ScraperOption {
	return func(s *Scraper) {
		s.SetRequestTimeout(timeout)
	}
}
func SetUserAgent(ua string) ScraperOption {
	return func(s *Scraper) {
		s.UserAgent = ua
	}
}

func IgnoreRobotsTxt() ScraperOption {
	return func(s *Scraper) {
		s.IgnoreRobotsTxt = true
	}
}
