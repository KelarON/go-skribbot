package utils

import (
	"errors"
	"image"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type Searcher struct {
	lastQuery string
	lastCount int
}

func NewSearcher() *Searcher {
	return &Searcher{}
}

const IMAGE_COUNT = 6

var QUERY_ADDONS = [...]string{
	" рисунок",
	" портрет",
	" pixel art",
	" drawing",
	" art",
	" sketch",
	" portrait",
	" мем",
	" чёрно-белый",
	" meme",
}

// SearchImages searches for images on Google and returns a slice of image.Image objects.
func (s *Searcher) SearchImages(query string) ([]image.Image, error) {

	if query == "" {
		return nil, errors.New("query cannot be empty")
	}

	if s.lastCount >= len(QUERY_ADDONS) {
		s.lastCount = 0
		s.lastQuery = ""
	}

	if query == s.lastQuery {
		query += QUERY_ADDONS[s.lastCount]
		s.lastCount++
	} else {
		s.lastCount = 0
		s.lastQuery = query
	}

	// prepare query params
	params := url.Values{}
	params.Add("q", query)
	params.Add("tbm", "isch")

	// make a GET request to Google's search engine
	resp, err := http.Get("https://www.google.com/search?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// find required image URLs from the response body
	by, err := io.ReadAll(resp.Body)
	urlList := regexp.MustCompile(`src="(.+?)"`).FindAllSubmatch(by, IMAGE_COUNT+1)[1:]

	imageList := make([]image.Image, 0, IMAGE_COUNT)

	// get images from URLs and add them to the list
	for _, url := range urlList {
		rsp, err := http.Get(string(url[1]))
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		img, _, err := image.Decode(rsp.Body)
		if err != nil {
			return nil, err
		}
		imageList = append(imageList, img)
	}
	return imageList, nil
}
