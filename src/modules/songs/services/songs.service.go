package services

import (
	"api/src/modules/songs/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SongsService struct{}

func NewSongsService() *SongsService {
	return &SongsService{}
}

func (s *SongsService) Search(query string) ([]dto.SongSearchResponse, error) {
	// Using CifraClub's public Solr API which is what the website uses now
	// Format: https://solr.sscdn.co/cc/c7/?q=QUERY&limit=10
	url := fmt.Sprintf("https://solr.sscdn.co/cc/c7/?q=%s&limit=30", strings.ReplaceAll(query, " ", "+"))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cifraclub returned status %d", res.StatusCode)
	}

	var apiRes struct {
		Response struct {
			Docs []struct {
				Txt  string `json:"txt"`  // Song Title
				Art  string `json:"art"`  // Artist Name
				Dns  string `json:"dns"`  // Artist Slug
				Url  string `json:"url"`  // Song Slug
				Tipo string `json:"tipo"` // Type (2 for song, 1 for artist)
			} `json:"docs"`
		} `json:"response"`
	}

	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		return nil, err
	}

	var songs []dto.SongSearchResponse
	for _, r := range apiRes.Response.Docs {
		// tipo 2 is a song. tipo 1 is an artist which we skip for now.
		if r.Tipo == "2" {
			songs = append(songs, dto.SongSearchResponse{
				Title:  r.Txt,
				Artist: r.Art,
				Url:    fmt.Sprintf("https://www.cifraclub.com.br/%s/%s/", r.Dns, r.Url),
			})
		}
	}

	return songs, nil
}

func (s *SongsService) GetSong(url string) (*dto.SongDetailResponse, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cifraclub returned status %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find(".t1").Text()
	artist := doc.Find(".t3").Text()
	key := doc.Find("#cifra_tom a").Text()

	content, _ := doc.Find("pre").Html()

	var chords []string
	doc.Find("#cifra_capo").NextAll().Find("b").Each(func(i int, s *goquery.Selection) {
		chords = append(chords, s.Text())
	})

	return &dto.SongDetailResponse{
		Title:   strings.TrimSpace(title),
		Artist:  strings.TrimSpace(artist),
		Key:     strings.TrimSpace(key),
		Chords:  chords,
		Content: content,
	}, nil
}

func (s *SongsService) GetTrending() ([]dto.SongSearchResponse, error) {
	// Try multiple URLs if needed, but start with the most specific gospel one
	url := "https://www.cifraclub.com.br/mais-tocadas/gospel-religioso/"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	// Modern real browser user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var songs []dto.SongSearchResponse

	// Cifra Club relies heavily on generated css modules, .primaryLabel and .secondaryLabel are stable hooks
	count := 0
	doc.Find("li").Each(func(i int, sel *goquery.Selection) {
		if count >= 12 { // Limit to 12 for the hero UI
			return
		}

		// Try to find the title using the new stable CSS classes
		title := sel.Find(".primaryLabel span").First().Text()
		if title == "" {
			title = sel.Find(".primaryLabel").Text()
		}

		// Fallbacks for older layout versions
		if title == "" {
			title = sel.Find("b").Text()
		}
		if title == "" {
			title = sel.Find(".mais-tocadas-song").Text()
		}

		artist := sel.Find(".secondaryLabel").Text()
		if artist == "" {
			artist = sel.Find("span").Text()
		}
		if artist == "" {
			artist = sel.Find(".mais-tocadas-artist").Text()
		}

		songUrl, _ := sel.Find("a").Attr("href")
		imgUrl, _ := sel.Find("img").Attr("src")

		if imgUrl == "" || strings.Contains(imgUrl, "placeholder") {
			imgUrl, _ = sel.Find("img").Attr("data-src")
		}

		// Cleaning up title if it contains rank numbers (e.g. "01")
		if len(title) > 2 && title[:2] == fmt.Sprintf("%02d", count+1) {
			title = title[2:]
		}

		if !strings.HasPrefix(songUrl, "http") && songUrl != "" && !strings.HasPrefix(songUrl, "javascript") {
			songUrl = "https://www.cifraclub.com.br" + songUrl
		}

		if title != "" && artist != "" && songUrl != "" {
			// Increase image resolution for better UI experience
			if strings.Contains(imgUrl, "-tb2.jpg") {
				imgUrl = strings.Replace(imgUrl, "-tb2.jpg", "-tb5.jpg", 1)
			}

			songs = append(songs, dto.SongSearchResponse{
				Title:  strings.TrimSpace(title),
				Artist: strings.TrimSpace(artist),
				Url:    songUrl,
				Image:  imgUrl,
			})
			count++
		}
	})

	return songs, nil
}
