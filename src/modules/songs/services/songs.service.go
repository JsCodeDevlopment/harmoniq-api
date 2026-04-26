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

	var versions []dto.SongVersion
	var simplifiedUrl, principalUrl, keyboardUrl string

	// Extract versions from the page
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || strings.HasPrefix(href, "javascript") || strings.Contains(href, "academy") {
			return
		}

		// Normalize URL
		fullUrl := href
		if !strings.HasPrefix(href, "http") {
			fullUrl = "https://www.cifraclub.com.br" + href
		}

		// Check if it's a version link (contains the song path and ends with .html or /)
		// We can use the current URL to determine the song path
		songPath := strings.TrimSuffix(url, "/")
		songPath = strings.TrimSuffix(songPath, ".html")
		// Get the artist/song part
		parts := strings.Split(strings.TrimPrefix(songPath, "https://www.cifraclub.com.br/"), "/")
		if len(parts) < 2 {
			return
		}
		basePath := "/" + parts[0] + "/" + parts[1] + "/"

		if strings.Contains(fullUrl, basePath) {
			name := ""
			// Try to find the name in a strong tag or first span, or just the first text node
			if strong := s.Find("strong"); strong.Length() > 0 {
				name = strings.TrimSpace(strong.Text())
			} else if span := s.Find("span").First(); span.Length() > 0 {
				name = strings.TrimSpace(span.Text())
			} else {
				name = strings.TrimSpace(s.Text())
			}

			// Clean up and normalize
			if strings.Contains(name, "Cifra:") {
				name = strings.TrimSpace(strings.Replace(name, "Cifra:", "", 1))
			}

			// Further cleanup if it still contains metadata (heuristic)
			// If name contains common metadata markers, truncate it
			for _, marker := range []string{"Básico", "Intermediário", "Avançado", "exibições"} {
				if idx := strings.Index(name, marker); idx != -1 {
					name = strings.TrimSpace(name[:idx])
					break
				}
			}

			// Categorize and add if it looks like a version name
			isVersion := false
			if strings.HasSuffix(fullUrl, basePath) || strings.HasSuffix(fullUrl, basePath+"index.html") {
				isVersion = true
				name = "Principal"
				principalUrl = fullUrl
			} else if strings.HasSuffix(fullUrl, "/simplificada.html") || strings.Contains(name, "Simplificada") {
				isVersion = true
				name = "Simplificada"
				simplifiedUrl = fullUrl
			} else if strings.HasSuffix(fullUrl, "/teclado.html") || strings.Contains(name, "Teclado") {
				isVersion = true
				name = "Teclado"
				keyboardUrl = fullUrl
			} else if strings.Contains(fullUrl, "versao-") || strings.Contains(strings.ToLower(name), "versao") {
				isVersion = true
				// Normalize name to "Versão X"
				if strings.Contains(fullUrl, "versao-") {
					vParts := strings.Split(fullUrl, "versao-")
					if len(vParts) > 1 {
						vNum := strings.TrimSuffix(vParts[1], ".html")
						name = "Versão " + vNum
					}
				}
			}

			if isVersion && name != "" {
				// Avoid duplicates
				exists := false
				for _, v := range versions {
					if v.Url == fullUrl || v.Name == name {
						exists = true
						break
					}
				}
				if !exists {
					versions = append(versions, dto.SongVersion{
						Name: name,
						Url:  fullUrl,
					})
				}
			}
		}
	})

	// Heuristic fallbacks if not found in links
	if principalUrl == "" {
		// Try to construct it from current URL
		parts := strings.Split(strings.TrimPrefix(url, "https://www.cifraclub.com.br/"), "/")
		if len(parts) >= 2 {
			principalUrl = "https://www.cifraclub.com.br/" + parts[0] + "/" + parts[1] + "/"
		}
	}

	if simplifiedUrl == "" && principalUrl != "" {
		simplifiedUrl = principalUrl + "simplificada.html"
	}
	if keyboardUrl == "" && principalUrl != "" {
		keyboardUrl = principalUrl + "teclado.html"
	}

	// Ensure Principal is in the list
	foundPrincipal := false
	for _, v := range versions {
		if v.Name == "Principal" {
			foundPrincipal = true
			break
		}
	}
	if !foundPrincipal && principalUrl != "" {
		versions = append([]dto.SongVersion{{Name: "Principal", Url: principalUrl}}, versions...)
	}

	// Scrape Artist Image
	artistImage, _ := doc.Find(".header-nav nav a img").Attr("src")
	if artistImage == "" {
		artistImage, _ = doc.Find(".t3 a img").Attr("src")
	}
	if artistImage == "" {
		artistImage, _ = doc.Find(".header-artist-img img").Attr("src")
	}
	if artistImage == "" {
		// Try to find any image that looks like an artist profile in the top section
		artistImage, _ = doc.Find("img[alt='" + artist + "']").Attr("src")
	}

	var recommendations []dto.SongSearchResponse
	// Strategy 1: "Toque também" (Related songs)
	// Using classes from Cifra Club's modern layout as shown in DevTools
	doc.Find(".playToo-listItem, .js-side-related a, .related-songs a, .cifra-footer-related a").Each(func(i int, s *goquery.Selection) {
		if len(recommendations) >= 12 {
			return
		}

		// Handle both the specific classes and generic fallbacks
		title := s.Find(".playToo--primaryText").Text()
		if title == "" {
			title = s.Find("strong").Text()
		}
		if title == "" {
			title = s.Find("b").Text()
		}

		artistName := s.Find(".playToo--secondaryText").Text()
		if artistName == "" {
			artistName = s.Find("span").First().Text()
		}

		href := ""
		if s.Is("a") {
			href, _ = s.Attr("href")
		} else {
			href, _ = s.Find("a").Attr("href")
		}

		// Image extraction
		img := ""
		imgSel := s.Find(".playToo--artistImage img")
		if imgSel.Length() > 0 {
			img, _ = imgSel.Attr("src")
		} else {
			// Check if it's in a data-src or a regular img inside the item
			imgSel = s.Find("img")
			img, _ = imgSel.Attr("src")
			if img == "" || strings.Contains(img, "placeholder") {
				img, _ = imgSel.Attr("data-src")
			}
		}

		// Special case for the thumb layout in screenshot: image might be in a background-image or child
		if img == "" {
			img, _ = s.Find(".thumb").Attr("data-src")
		}

		if title != "" && href != "" {
			// Avoid duplicates
			isDuplicate := false
			for _, r := range recommendations {
				if r.Title == title && r.Artist == artistName {
					isDuplicate = true
					break
				}
			}

			if !isDuplicate {
				if !strings.HasPrefix(href, "http") {
					href = "https://www.cifraclub.com.br" + href
				}

				// Ensure we have a high quality image if possible
				if strings.Contains(img, "-tb2.jpg") {
					img = strings.Replace(img, "-tb2.jpg", "-tb5.jpg", 1)
				}

				recommendations = append(recommendations, dto.SongSearchResponse{
					Title:  strings.TrimSpace(title),
					Artist: strings.TrimSpace(artistName),
					Url:    href,
					Image:  img,
				})
			}
		}
	})

	// Strategy 2: Popular songs from the same artist (fallback)
	if len(recommendations) < 4 {
		doc.Find(".art_musics li").Each(func(i int, s *goquery.Selection) {
			if len(recommendations) >= 12 {
				return
			}
			link := s.Find("a")
			songTitle := link.Find("div div div").Text()
			if songTitle == "" {
				songTitle = link.Find("span").First().Text()
			}
			if songTitle == "" {
				songTitle = link.Text()
			}

			href, _ := link.Attr("href")
			if songTitle != "" && href != "" {
				if !strings.HasPrefix(href, "http") {
					href = "https://www.cifraclub.com.br" + href
				}

				// Avoid adding current song as recommendation
				if strings.Contains(href, url) || strings.Contains(url, href) {
					return
				}

				recommendations = append(recommendations, dto.SongSearchResponse{
					Title:  strings.TrimSpace(songTitle),
					Artist: strings.TrimSpace(artist),
					Url:    href,
					Image:  artistImage,
				})
			}
		})
	}

	return &dto.SongDetailResponse{
		Title:           strings.TrimSpace(title),
		Artist:          strings.TrimSpace(artist),
		Key:             strings.TrimSpace(key),
		Chords:          chords,
		Content:         content,
		SimplifiedUrl:   simplifiedUrl,
		PrincipalUrl:    principalUrl,
		KeyboardUrl:     keyboardUrl,
		Versions:        versions,
		Recommendations: recommendations,
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
